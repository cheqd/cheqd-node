module github.com/cheqd/cheqd-node

go 1.15

require (
	github.com/armon/go-metrics v0.3.8 // indirect
	github.com/confio/ics23/go v0.6.6 // indirect
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.3-0.20201103224600-674baa8c7fc3 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/otiai10/copy v1.6.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/prometheus/common v0.23.0 // indirect
	github.com/regen-network/cosmos-proto v0.3.1 // indirect
	github.com/rs/zerolog v1.21.0 // indirect
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.1.3 // indirect
	google.golang.org/genproto v0.0.0-20210207032614-bba0dbe2a9ea
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	gopkg.in/ini.v1 v1.61.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
