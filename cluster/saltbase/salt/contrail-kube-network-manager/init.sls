/etc/kubernetes/manifests/contrail-kube-network-manager.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/kube-network-manager.manifest
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
