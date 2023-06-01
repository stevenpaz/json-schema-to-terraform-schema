package main

import (
	"fmt"
	"os"

	"github.com/stevenpaz/tf-schema-gen/openapi"
)

func main() {
	RunOpenAPIGen()
}

func RunOpenAPIGen() {
	// check args
	if len(os.Args) != 3 {
		fmt.Println("usage: tf-schema-gen <openapi.yaml> <output-folder>")
		os.Exit(1)
	}

	// get filePath and outputFolderPath from args
	filePath := os.Args[1]
	outputFolderPath := os.Args[2]

	err := openapi.CreateTFSchemaFromOpenAPI(filePath, outputFolderPath)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
