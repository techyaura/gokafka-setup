# Prerequisite

### OS
```
UBUNTU 18.04
```

### Programming Language

```
Golang 1.13 +
```

### Other Software
```
Docker
```

### Install librdkafka

```
git clone https://github.com/edenhill/librdkafka.git
cd librdkafka
./configure --prefix /usr
make
sudo make install

```

### Update hosts

```
OPEN - sudo nano /etc/hosts
ADD - 0.0.0.0 kafka
SAVE & EXIT
```

# HOW TO RUN

```
1. Start Kafka & zookeeper
2. RUN `docker-compose up`
3. git clone `https://github.com/techyaura/gokafka-setup.git`
4. cd gokafka-setup/producer
6. RUN ./producer <HOST: 0.0.0.0> '<message>'
7. Open new terminal
8. cd gokafka-setup/consumer
8. RUN ./consumer <HOST: 0.0.0.0>
OBSERVE
```

# TEST
```
RUN>
cd gokafka-setup/consumer
./producer 0.0.0.0 '{"message": "Hello world"}

OBSERVE> On other terminal by running
./consumer.go 0.0.0.0
```


NOTE: On running `docker-compose up`, KafkaManager (https://github.com/yahoo/CMAK) will also start, which can be access on `localhost:9000`. It provides a decent UI for access or managing KAFKA.
Default credentials - `U/P`: `admin/admin`
