/etc/kubernetes/manifests/contrail-rabbitmq.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/rabbitmq.manifest
    - source_hash: md5=fe26b2f66c2adfd94c85c1971be50267
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
