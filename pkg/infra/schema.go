package infra

type Schema struct {
	Services    []Service    `json:"services"`
	Entrypoints []Entrypoint `json:"entrypoints"`
	Websites    []Website    `json:"websites"`
	Topics      []Topic      `json:"topics"`
	Buckets     []Bucket     `json:"buckets"`
}
