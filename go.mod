module github.com/nitrictech/newcli

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.3.2
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/docker v20.10.11+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/fatih/color v1.13.0
	github.com/getkin/kin-openapi v0.88.0
	github.com/golangci/golangci-lint v1.43.0
	github.com/google/go-github/v41 v41.0.0
	github.com/hashicorp/consul/sdk v0.8.0
	github.com/jhoonb/archivex v0.0.0-20201016144719-6a343cdae81d
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.4.2
	github.com/nitrictech/apis v0.13.0-rc.4
	github.com/nitrictech/boxygen v0.0.1-rc.7.0.20211212231606-62c668408f91
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	google.golang.org/genproto v0.0.0-20211005153810-c76a74d43a8e // indirect
	google.golang.org/grpc v1.41.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/opencontainers/image-spec => github.com/opencontainers/image-spec v1.0.2-0.20211123152302-43a7dee1ec31
	github.com/rootless-containers/rootlesskit => github.com/rootless-containers/rootlesskit v0.14.6
)
