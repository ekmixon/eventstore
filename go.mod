module github.com/triggermesh/eventstore

go 1.15

replace k8s.io/client-go => k8s.io/client-go v0.19.7

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.33.2
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/code-generator v0.19.7
)
