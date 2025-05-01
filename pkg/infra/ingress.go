package infra

type TargetType string

const (
	TargetType_Service TargetType = "service"
	TargetType_Website TargetType = "website"
)

type Ingress struct {
	Name    string   `json:"name"`
	Targets []Target `json:"targets"`
}

type Target struct {
	Type       TargetType `json:"type" jsonschema:"enum=service,enum=website"`
	TargetName string     `json:"targetName"`
	Path       string     `json:"path"`
}
