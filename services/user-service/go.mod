module github.com/lucas/gokafka/user-service

go 1.22.2

require (
	github.com/google/uuid v1.6.0
	github.com/segmentio/kafka-go v0.4.48
	golang.org/x/crypto v0.14.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lib/pq v1.10.9
	github.com/lucas/gokafka/shared v0.0.0
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace github.com/lucas/gokafka/shared => ../../shared
