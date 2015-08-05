/etc/kubernetes/manifests/contrail-ifmap-server.manifest:
  file.managed:
    - source: https://raw.githubusercontent.com/pedro-r-marques/contrail-kubernetes/manifests/cluster/ifmap-server.manifest
    - source_hash: md5=7346ad0f28610b69760a2fdf052eebc7
    - user: root
    - group: root
    - mode: 644
    - makedirs: true
    - dir_mode: 755
