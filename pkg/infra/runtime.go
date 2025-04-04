package infra

// Runtime represents a union of all possible runtime types
type Runtime struct {
	DockerfileRuntime  *DockerFileRuntime  `json:"dockerfile,omitempty" jsonschema:"oneof_required=dockerfile"`
	DockerImageRuntime *DockerImageRuntime `json:"docker,omitempty" jsonschema:"oneof_required=docker"`
}

// DockerFileRuntime represents a runtime that uses a Dockerfile
type DockerFileRuntime struct {
	Dockerfile string `json:"dockerfile"`
}

type DockerImageRuntime struct {
	Image string `json:"image"`
}
