package analyzer

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// Resource represents a Terraform resource.
type Resource struct {
	Type string
	Name string
	Body hcl.Body
}

// ParseTerraformFiles parses Terraform HCL files in the given directory
// and returns a list of parsed Resource structs.
func ParseTerraformFiles(dir string) ([]*Resource, error) {
	parser := hclparse.NewParser()
	var resources []*Resource

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".tf" {
			filePath := filepath.Join(dir, file.Name())
			f, diags := parser.ParseHCLFile(filePath)
			if diags.HasErrors() {
				for _, diag := range diags {
					log.Printf("Error in file %s: %s\n", filePath, diag.Error())
				}
				return nil, fmt.Errorf("failed to parse file %s: %w", filePath, diags)
			}

			content, diags := f.Body.Content(&hcl.BodySchema{
				Blocks: []hcl.BlockHeaderSchema{
					{Type: "resource"},
					{Type: "data"},
					// Add schemas for other HCL elements here if needed
				},
			})
			if diags.HasErrors() {
				return nil, fmt.Errorf("failed to extract content from file %s: %w", filePath, diags)
			}

			for _, block := range content.Blocks {
				resource := &Resource{
					Type: block.Labels[0],
					Name: block.Labels[1],
					Body: block.Body,
				}
				resources = append(resources, resource)
			}
		}
	}

	return resources, nil
}
