module github.com/opencord/opendevice-manager

go 1.13

replace google.golang.org/grpc v1.29.1 => google.golang.org/grpc v1.26.0

require (
	github.com/Shopify/sarama v1.23.1
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.2.0
	github.com/jinzhu/copier v0.2.9
	github.com/opencord/device-management-interface v0.12.0
	github.com/opencord/voltha-lib-go/v4 v4.2.4
	google.golang.org/grpc v1.29.1
)
