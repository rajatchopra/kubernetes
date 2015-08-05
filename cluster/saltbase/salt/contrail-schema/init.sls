/etc/kubernetes/manifests/contrail-schema.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/contrail-schema.manifest
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
