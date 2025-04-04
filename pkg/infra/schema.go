package infra

type Schema struct {
	Services    []Service    `json:"services,omitempty"`
	Entrypoints []Entrypoint `json:"entrypoints,omitempty"`
	Websites    []Website    `json:"websites,omitempty"`
	Topics      []Topic      `json:"topics,omitempty"`
	Buckets     []Bucket     `json:"buckets,omitempty"`
}
