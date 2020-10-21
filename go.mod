module github.com/triggermesh/eventstore

go 1.15

// Top-level module control over the exact version used for important direct dependencies.
// https://github.com/golang/go/wiki/Modules#when-should-i-use-the-replace-directive
replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.8
	k8s.io/client-go => k8s.io/client-go v0.16.8
	k8s.io/code-generator => k8s.io/code-generator v0.16.8
)

require google.golang.org/grpc v1.33.0

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.1 // indirect
	github.com/aws/aws-sdk-go v1.30.16 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/go-cmp v0.4.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.2 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	go.opencensus.io v0.22.4
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.0 // indirect
	google.golang.org/protobuf v1.23.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/api v0.17.6 // indirect
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29 // indirect
	knative.dev/pkg v0.0.0-20200625173728-dfb81cf04a7c
	sigs.k8s.io/yaml v1.2.0 // indirect
)
