package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func yml2Json(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("failed to read data: %s", err)
	}
	content := make(map[string]any)
	err = yaml.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(content, "", "  ")
}
func json2Yml(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("failed to read data: %s", err)
	}
	content := make(map[string]any)
	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(content)
}

func main() {
	fileName := flag.String("file", "", "File to open to convert to json or yml")
	t := flag.String("type", "yaml", "type to convert to")
	flag.Parse()
	r := os.Stdin
	var err error
	if *fileName != "" {
		r, err = os.Open(*fileName)
		if err != nil {
			log.Fatalf("failed to read data from stdin: %s", err)
		}
	}
	switch *t {
	case "yaml", "yml":
		data, err := json2Yml(r)
		if err != nil {
			log.Fatalf("failed to convert json to yaml format: %s", err)
		}
		fmt.Printf("%s", data)
	case "json":
		data, err := yml2Json(r)
		if err != nil {
			log.Fatalf("failed to convert yaml to json format: %s", err)
		}
		fmt.Printf("%s\n", data)
	default:
		log.Fatalf("unknown type provided: %s", *t)
	}

}
