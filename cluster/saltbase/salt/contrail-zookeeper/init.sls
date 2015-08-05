/etc/kubernetes/manifests/contrail-zookeeper.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/zookeeper.manifest
    - source_hash: md5=c665836b3d5fe7b535d2c91a5efbb824
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
