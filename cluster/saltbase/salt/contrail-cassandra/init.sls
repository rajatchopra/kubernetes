/etc/kubernetes/manifests/contrail-cassandra.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/cassandra.manifest
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
