# consumers-status

##### Add kafka servers to /etc/hosts
```
192.168.1.2 kafka-01.local1.server
192.168.1.3 kafka-02.local2.server
192.168.1.4 kafka-03.local3.server
```
* You should be able to connect to the kafka brokers via port 9092

##### Download kafka and extract it in /usr/local/
```
sudo -i
wget -c wget -c https://downloads.apache.org/kafka/3.6.0/kafka_2.13-3.6.0.tgz
tar -xzf kafka_2.13-3.6.0.tgz -C /usr/local
mv /usr/local/kafka_2.13-3.6.0 /usr/local/kafka
```

##### Why node_exporter?
If you put a prom file in `/var/lib/node_exporter/textfile_collector/` path your metrics expose along other
node_exporter metrics.