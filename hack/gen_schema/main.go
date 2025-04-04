package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/nitrictech/cli/pkg/infra"
)

func main() {
	schema := jsonschema.Reflect(&infra.Schema{})

	jsonOutput, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(string(jsonOutput))
}
