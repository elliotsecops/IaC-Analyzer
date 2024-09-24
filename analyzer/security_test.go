package analyzer

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

// Helper function to create a Resource struct for testing
func createResource(resourceType, resourceName string, body hcl.BodyContent) *Resource {
	return &Resource{
		Type: resourceType,
		Name: resourceName,
		Body: body,
	}
}

func TestCheckOpenSSH(t *testing.T) {
	testCases := []struct {
		name     string
		resource *Resource
		expected *Issue
	}{
		{
			name: "Open SSH Access from Anywhere",
			resource: createResource("aws_security_group", "test-sg", hcl.BodyContent{
				Blocks: hcl.Blocks{
					{
						Type: "ingress",
						Body: hcl.BodyContent{
							Attributes: hcl.Attributes{
								"from_port":   {Expr: &hcl.LiteralValueExpr{Val: cty.NumberIntVal(22)}},
								"to_port":     {Expr: &hcl.LiteralValueExpr{Val: cty.NumberIntVal(22)}},
								"cidr_blocks": {Expr: &hcl.LiteralValueExpr{Val: cty.TupleVal([]cty.Value{cty.StringVal("0.0.0.0/0")})}},
							},
						},
					},
				},
			}),
			expected: NewIssue("HIGH", "Open SSH access in security group", "test-sg", 0),
		},
		{
			name: "Open SSH Access from Specific IP",
			resource: createResource("aws_security_group", "test-sg", hcl.BodyContent{
				Blocks: hcl.Blocks{
					{
						Type: "ingress",
						Body: hcl.BodyContent{
							Attributes: hcl.Attributes{
								"from_port":   {Expr: &hcl.LiteralValueExpr{Val: cty.NumberIntVal(22)}},
								"to_port":     {Expr: &hcl.LiteralValueExpr{Val: cty.NumberIntVal(22)}},
								"cidr_blocks": {Expr: &hcl.LiteralValueExpr{Val: cty.TupleVal([]cty.Value{cty.StringVal("192.168.1.0/24")})}},
							},
						},
					},
				},
			}),
			expected: nil,
		},
		{
			name: "No Open SSH Access",
			resource: createResource("aws_security_group", "test-sg", hcl.BodyContent{
				Blocks: hcl.Blocks{
					{
						Type: "ingress",
						Body: hcl.BodyContent{
							Attributes: hcl.Attributes{
								"from_port":   {Expr: &hcl.LiteralValueExpr{Val: cty.NumberIntVal(80)}},
								"to_port":     {Expr: &hcl.LiteralValueExpr{Val: cty.NumberIntVal(80)}},
								"cidr_blocks": {Expr: &hcl.LiteralValueExpr{Val: cty.TupleVal([]cty.Value{cty.StringVal("192.168.1.0/24")})}},
							},
						},
					},
				},
			}),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issue := CheckOpenSSH(tc.resource)
			assert.Equal(t, tc.expected, issue)
		})
	}
}

func TestCheckPublicS3(t *testing.T) {
	testCases := []struct {
		name     string
		resource *Resource
		expected *Issue
	}{
		{
			name: "Public Read Access",
			resource: createResource("aws_s3_bucket", "test-bucket", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"acl": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("public-read")}},
				},
			}),
			expected: NewIssue("HIGH", "Public access enabled on S3 bucket", "test-bucket", 0),
		},
		{
			name: "Public Read-Write Access",
			resource: createResource("aws_s3_bucket", "test-bucket", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"acl": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("public-read-write")}},
				},
			}),
			expected: NewIssue("HIGH", "Public access enabled on S3 bucket", "test-bucket", 0),
		},
		{
			name: "Authenticated Read Access",
			resource: createResource("aws_s3_bucket", "test-bucket", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"acl": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("authenticated-read")}},
				},
			}),
			expected: NewIssue("HIGH", "Public access enabled on S3 bucket", "test-bucket", 0),
		},
		{
			name: "Private Access",
			resource: createResource("aws_s3_bucket", "test-bucket", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"acl": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("private")}},
				},
			}),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issue := CheckPublicS3(tc.resource)
			assert.Equal(t, tc.expected, issue)
		})
	}
}

func TestCheckHardcodedSecrets(t *testing.T) {
	testCases := []struct {
		name     string
		resource *Resource
		expected *Issue
	}{
		{
			name: "Hardcoded Secret in User Data",
			resource: createResource("aws_instance", "test-instance", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"user_data": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal(`password="secret"`)}},
				},
			}),
			expected: NewIssue("HIGH", "Hardcoded secret detected", "test-instance", 0),
		},
		{
			name: "Hardcoded Secret in Environment Variables",
			resource: createResource("aws_instance", "test-instance", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"environment": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal(`SECRET_KEY="value"`)}},
				},
			}),
			expected: NewIssue("HIGH", "Hardcoded secret detected", "test-instance", 0),
		},
		{
			name: "No Hardcoded Secret",
			resource: createResource("aws_instance", "test-instance", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"user_data": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal(`echo "hello"`)}},
				},
			}),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issue := CheckHardcodedSecrets(tc.resource)
			assert.Equal(t, tc.expected, issue)
		})
	}
}
