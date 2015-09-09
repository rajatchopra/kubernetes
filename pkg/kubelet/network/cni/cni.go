/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cni

import (
	"fmt"
	"github.com/appc/cni/libcni"
	cniTypes "github.com/appc/cni/pkg/types"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/kubelet/network"
	kubeletTypes "k8s.io/kubernetes/pkg/kubelet/types"
	"net"
	"sort"
)

const (
	CNIPluginName        = "cni"
	DefaultPluginName    = "kubernetes-bridge"
	DefaultNetDir        = "/etc/cni/net.d"
	DefaultCNIDir        = "/opt/cni/bin"
	VendorCNIDirTemplate = "/opt/%s/bin"
	DefaultPodNetworkKey = "net.experimental.kubernetes.io/networks"
)

type cniNetworkPlugin struct {
	cniNetworkMap map[string]*cniNetwork
	defaultNetwork *cniNetwork
	host          network.Host
}

type cniNetwork struct {
	name          string
	NetworkConfig *libcni.NetworkConfig
	CNIConfig     *libcni.CNIConfig
}

func ProbeNetworkPlugins() []network.NetworkPlugin {
	configList := make([]network.NetworkPlugin, 0)
	allNetworks, defaultNetwork := getAllCNINetworks()
	return append(configList, &cniNetworkPlugin{cniNetworkMap: getAllCNINetworks(), defaultNetwork: defaultNetwork})
}

func getAllCNINetworks() (map[string]*cniNetwork, *cniNetwork) {
	defaultNetwork := nil
	networkMap := make(map[string]*cniNetwork)
	files, err := libcni.ConfFiles(DefaultNetDir)
	switch {
	case err != nil:
		return networkMap, nil
	case len(files) == 0:
		return networkMap, nil
	}

	sort.Strings(files)
	for _, confFile := range files {
		conf, err := libcni.ConfFromFile(confFile)
		if err != nil {
			glog.Warningf("Error loading CNI config file %s: %v", confFile, err)
			continue
		}
		// Search for vendor-specific plugins as well as default plugins in the CNI codebase.
		vendorCNIDir := fmt.Sprintf(VendorCNIDirTemplate, conf.Network.Name)
		cninet := &libcni.CNIConfig{
			Path: []string{DefaultCNIDir, vendorCNIDir},
		}
		network := &cniNetwork{name: conf.Network.Name, NetworkConfig: conf, CNIConfig: cninet}
		if defaultNetwork == nil {
			defaultNetwork = net
		}
		networkMap[conf.Network.Name] = network
	}
	return networkMap, defaultNetwork
}

func (plugin *cniNetworkPlugin) selectNetworks(pod *api.Pod) ([]*cniNetwork, error) {
	selectedNetworks := make([]*cniNetwork, 0)

	// loading networks all over again maybe inefficient, but required if networks can be created on the fly
	plugin.cniNetworkMap,_ = getAllCNINetworks()

	if len(plugin.cniNetworkMap) == 0 {
		return nil, fmt.Errorf("No available CNI network available in %s", DefaultNetDir)
	}

	// check if the namespace of the pod is a network name itself
	network, ok := plugin.cniNetworkMap[pod.Namespace]
	if ok {
		selectedNetworks = append(selectedNetworks, network)
	}

	// label called "network"?
	var netNames []string
	err := json.Marshal(pod.Labels[DefaultPodNetworkKey], &netNames)
	if err != nil {
		for _, netName := range(netNames) {
			network, ok = plugin.cniNetworkMap[netName]
			if ok {
				selectedNetworks = append(selectedNetworks, network)
			}
		}
	}

	// annotation called "network"?
	err := json.Marshal(pod.Annotations[DefaultPodNetworkKey], &netNames)
	if err != nil {
		for _, netName := range(netNames) {
			network, ok = plugin.cniNetworkMap[netName]
			if ok {
				selectedNetworks = append(selectedNetworks, network)
			}
		}
	}

	if len(selectedNetworks)==0 {
		// return the default one
		selectedNetworks = append(selectedNetworks, plugin.defaultNetwork)
	}

	return selectedNetworks, nil
}

func (plugin *cniNetworkPlugin) Init(host network.Host) error {
	plugin.host = host
	return nil
}

func (plugin *cniNetworkPlugin) Name() string {
	return CNIPluginName
}

func (plugin *cniNetworkPlugin) SetUpPod(namespace string, name string, id kubeletTypes.DockerID) error {
	netns, err := plugin.host.GetRuntime().GetNetNs(string(id))
	if err != nil {
		return err
	}

	pod, ok := plugin.host.GetPodByName(namespace, name)
	if !ok {
		return fmt.Errorf("pod %q namespace %q: unable to find pod", name, namespace)
	}

	// TODO: pick one network? all networks? which one writes to Status.PodIP?
	// Perhaps pick the network by some annotation/label on pod, or the namespace.
	networks, err := plugin.selectNetworks(pod)
	if err != nil {
		return err
	}
	for network,_ := range(networks) {
		res, err := network.addToNetwork(name, namespace, string(id), netns)
		if err != nil {
			return err
		}

		var ip string
		if res.IP4 != nil {
			ip = res.IP4.IP.String()
		} else {
			ip = res.IP6.IP.String()
		}
		// TODO-PAT: check that PodIP can be an IPv6.
		// TODO-rajat: this keeps on updating the same field, pick the main one
		pod.Status.PodIP = ip
	}
	return err
}

func (plugin *cniNetworkPlugin) TearDownPod(namespace string, name string, id kubeletTypes.DockerID) error {
	netns, err := plugin.host.GetRuntime().GetNetNs(string(id))
	if err != nil {
		return err
	}
	pod, ok := plugin.host.GetPodByName(namespace, name)
	if !ok {
		return fmt.Errorf("pod %q namespace %q: unable to find pod", name, namespace)
	}

	networks, err := plugin.selectNetworks(pod)
	if err != nil {
		return err
	}
	for network,_ := range(networks) {
		err := network.deleteFromNetwork(name, namespace, string(id), netns)
		if er != nil {
			return err
		}
	}
	return nil
}

func (plugin *cniNetworkPlugin) Status(namespace string, name string, id kubeletTypes.DockerID) (*network.PodNetworkStatus, error) {
	pod, ok := plugin.host.GetPodByName(namespace, name)
	if !ok {
		return nil, fmt.Errorf("pod %q namespace %q: unable to find pod", name, namespace)
	}

	return &network.PodNetworkStatus{IP: net.ParseIP(pod.Status.PodIP)}, nil
}

func (network *cniNetwork) addToNetwork(podName string, podNamespace string, podInfraContainerID string, podNetnsPath string) (*cniTypes.Result, error) {
	rt, err := buildCNIRuntimeConf(podName, podNamespace, podInfraContainerID, podNetnsPath)
	if err != nil {
		glog.Errorf("Error adding network: %v", err)
		return nil, err
	}

	netconf, cninet := network.NetworkConfig, network.CNIConfig
	glog.V(2).Infof("About to run with conf.Network.Type=%v, c.Path=%v", netconf.Network.Type, cninet.Path)
	res, err := cninet.AddNetwork(netconf, rt)
	if err != nil {
		glog.Errorf("Error adding network: %v", err)
		return nil, err
	}

	return res, nil
}

func (network *cniNetwork) deleteFromNetwork(podName string, podNamespace string, podInfraContainerID string, podNetnsPath string) error {
	rt, err := buildCNIRuntimeConf(podName, podNamespace, podInfraContainerID, podNetnsPath)
	if err != nil {
		glog.Errorf("Error deleting network: %v", err)
		return err
	}

	netconf, cninet := network.NetworkConfig, network.CNIConfig
	glog.V(2).Infof("About to run with conf.Network.Type=%v, c.Path=%v", netconf.Network.Type, cninet.Path)
	err = cninet.DelNetwork(netconf, rt)
	if err != nil {
		glog.Errorf("Error deleting network: %v", err)
		return err
	}
	return nil
}

func buildCNIRuntimeConf(podName string, podNs string, podInfraContainerID string, podNetnsPath string) (*libcni.RuntimeConf, error) {
	glog.V(2).Infof("Got netns path %v", podNetnsPath)
	glog.V(2).Infof("Using netns path %v", podNs)

	rt := &libcni.RuntimeConf{
		ContainerID: podInfraContainerID,
		NetNS:       podNetnsPath,
		IfName:      "eth0",
		Args: [][2]string{
			{"K8S_POD_NAMESPACE", podNs},
			{"K8S_POD_NAME", podName},
			{"K8S_POD_INFRA_CONTAINER_ID", podInfraContainerID},
		},
	}

	return rt, nil
}
