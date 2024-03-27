module message-service

go 1.20

require (
	common v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.47
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
)

replace common => ../common
