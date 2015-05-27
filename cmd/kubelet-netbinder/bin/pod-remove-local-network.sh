#!/bin/sh


# ARGS : docker_id pod_name ip mac netid ovs_port

docker_id=$1
pod_name=$2
ip=$3
mac=$4
netid=$5
ovs_port=$6

iphex=`printf '%02x' ${ip//./ }; echo`
machex=`printf '%s' ${mac//:/ }; echo`


## del flows

# del rule that adds vnid to outbound traffic from pod
ovs-ofctl del-flows -O OpenFlow13 obr0 "table=0,cookie=${netid}/0xffffffff,in_port=${ovs_port}"
# del rule to identify inbound traffic for pod and redirect it to correct ovs-port/veth-pair
ovs-ofctl del-flows -O OpenFlow13 obr0 "table=1,cookie=${netid}/0xffffffff,tun_id=${netid},dl_dst=${mac}"
# rule to respond to arp requests that come locally from this host
ovs-ofctl del-flows -O OpenFlow13 obr0 "table=1,cookie=${netid}/0xffffffff,tun_id=${netid},dl_type=0x0806,nw_dst=${ip}"

# We don't get the docker ID here so we have to find it ourselves
CONTAINERS=`grep -m1 -riIl "HOSTNAME=${pod_name}" /var/lib/docker/execdriver/native/ | sed 's/\/container.json//g'`
for cdir in $CONTAINERS; do
	veth_host=`jq .network_state.veth_host ${cdir}/state.json | tr -d '"'`
	if [ -n $veth_host -a "$veth_host" != "null" ]; then
		ovs-vsctl del-port obr0 $veth_host
	fi
done

