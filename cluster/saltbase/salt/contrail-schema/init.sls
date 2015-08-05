/etc/kubernetes/manifests/contrail-schema.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/contrail-schema.manifest
    - source_hash: md5=494f796f276d0cff719f414734a7b476
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
