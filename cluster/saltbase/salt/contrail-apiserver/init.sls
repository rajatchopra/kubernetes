/etc/kubernetes/manifests/contrail-apiserver.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/contrail-api.manifest
    - source_hash: md5=cf18f795f064bb6cb2f4eece99366c57
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
