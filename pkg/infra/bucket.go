package infra

type Bucket struct {
	Name     string          `json:"name"`
	Triggers []BucketTrigger `json:"triggers,omitempty"`
}

type BucketTriggerType string

const (
	BucketTriggerType_Write  BucketTriggerType = "write"
	BucketTriggerType_Delete BucketTriggerType = "delete"
)

type BucketTrigger struct {
	TriggerType BucketTriggerType   `json:"triggerType" jsonschema:"enum=write,enum=delete"`
	KeyPrefix   string              `json:"keyPrefix"`
	Target      BucketTriggerTarget `json:"target"`
}

type BucketTriggerTargetType string

const (
	BucketTriggerTargetType_Service BucketTriggerTargetType = "service"
)

type BucketTriggerTarget struct {
	TargetType BucketTriggerTargetType `json:"type" jsonschema:"enum=service"`
	TargetName string                  `json:"name"`
	Path       string                  `json:"path"`
}
