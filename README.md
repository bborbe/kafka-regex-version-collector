# Kafka Regex Version Collector

Fetch a URL and parse available versions via regex.

## Run version collector

```bash
go run main.go \
-kafka-brokers=kafka:9092 \
-kafka-topic=application-version-available \
-kafka-schema-registry-url=http://localhost:8081 \
-application=Golang \
-url=https://golang.org/dl/ \
-regex='https://dl.google.com/go/go([\d\.]+)\.src\.tar\.gz'
-v=2
```
