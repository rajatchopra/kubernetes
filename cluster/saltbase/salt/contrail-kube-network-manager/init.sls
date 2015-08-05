/etc/kubernetes/manifests/contrail-kube-network-manager.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/kube-network-manager.manifest
    - source_hash: md5=0ad938ed7db88b455d83834b2cc29a2a
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
