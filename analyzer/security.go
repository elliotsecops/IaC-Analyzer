package analyzer

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// SecurityCheck represents a security check function.
type SecurityCheck func(*Resource) *Issue

// CheckOpenSSH checks for open SSH access in security groups.
func CheckOpenSSH(resource *Resource) *Issue {
	if resource.Type == "aws_security_group" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{Type: "ingress"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing security group: %v", diags), resource.Name, 0)
		}

		for _, block := range bodyContent.Blocks {
			if isSSHAccessOpen(block) {
				return NewIssue("HIGH", "Open SSH access in security group", resource.Name, block.DefRange.Start.Line)
			}
		}
	}
	return nil
}

// isSSHAccessOpen checks if the given ingress block allows open SSH access.
func isSSHAccessOpen(block *hcl.Block) bool {
	content, _, diags := block.Body.PartialContent(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "from_port"},
			{Name: "to_port"},
			{Name: "cidr_blocks"},
		},
	})
	if diags.HasErrors() {
		return false
	}

	fromPort, _ := content.Attributes["from_port"].Expr.Value(nil)
	toPort, _ := content.Attributes["to_port"].Expr.Value(nil)

	if fromPort.Equals(cty.NumberIntVal(22)) && toPort.Equals(cty.NumberIntVal(22)) {
		cidrBlocks, _ := content.Attributes["cidr_blocks"].Expr.Value(nil)
		if cidrBlocks.Type().IsTupleType() {
			for _, cidr := range cidrBlocks.AsValueSlice() {
				if cidr.Equals(cty.StringVal("0.0.0.0/0")) {
					return true
				}
			}
		}
	}
	return false
}

// CheckPublicS3 checks for publicly accessible S3 buckets.
func CheckPublicS3(resource *Resource) *Issue {
	if resource.Type == "aws_s3_bucket" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "acl"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing S3 bucket: %v", diags), resource.Name, 0)
		}

		acl, _ := bodyContent.Attributes["acl"].Expr.Value(nil)
		if isPublicACL(acl) {
			return NewIssue("HIGH", "Public access enabled on S3 bucket", resource.Name, bodyContent.Attributes["acl"].Expr.Range().Start.Line)
		}
	}
	return nil
}

// isPublicACL checks if the given ACL value indicates public access.
func isPublicACL(acl cty.Value) bool {
	publicACLs := []cty.Value{
		cty.StringVal("public-read"),
		cty.StringVal("public-read-write"),
		cty.StringVal("authenticated-read"),
	}
	for _, publicACL := range publicACLs {
		if acl.Equals(publicACL).True() {
			return true
		}
	}
	return false
}

// CheckHardcodedSecrets checks for hardcoded secrets using regex patterns.
func CheckHardcodedSecrets(resource *Resource) *Issue {
	secretsRegex := regexp.MustCompile(`(?i)(password|secret|token|key)\s*=\s*".+"`)
	bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "user_data"},
			{Name: "environment"},
		},
	})
	if diags.HasErrors() {
		return NewIssue("ERROR", fmt.Sprintf("Error parsing resource: %v", diags), resource.Name, 0)
	}

	for _, attr := range bodyContent.Attributes {
		if secretsRegex.MatchString(attr.Expr.Range().String()) {
			return NewIssue("HIGH", "Hardcoded secret detected", resource.Name, attr.Expr.Range().Start.Line)
		}
	}
	return nil
}

// RunSecurityChecks runs all security checks on the given resources.
func RunSecurityChecks(resources []*Resource) []Issue {
	var issues []Issue
	checks := []SecurityCheck{
		CheckOpenSSH,
		CheckPublicS3,
		CheckHardcodedSecrets,
	}

	for _, resource := range resources {
		for _, check := range checks {
			if issue := check(resource); issue != nil {
				issues = append(issues, *issue)
			}
		}
	}

	return issues
}
