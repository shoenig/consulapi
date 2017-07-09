/tmp/consul agent -bootstrap -config-dir hack/consul-config -data-dir /tmp/consuldata &
sleep 3
ps -ef | grep consul