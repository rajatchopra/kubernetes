/etc/kubernetes/manifests/contrail-cassandra.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/cassandra.manifest
    - source_hash: md5=2193a46160758f71c933c46c880125b5
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
