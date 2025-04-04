package infra

type ServiceType string

const (
	ServiceType_Serverless ServiceType = "serverless"
	ServiceType_Container  ServiceType = "container"
)

type Service struct {
	Name    string      `json:"name"`
	Type    ServiceType `json:"type" jsonschema:"enum=serverless,enum=container"`
	Port    int         `json:"port"`
	Lang    string      `json:"lang"`
	Runtime Runtime     `json:"runtime"`
}
