/etc/kubernetes/manifests/contrail-control.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/contrail-control.manifest
    - source_hash: md5=c709bc650cdcfc274f26310194261489
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
