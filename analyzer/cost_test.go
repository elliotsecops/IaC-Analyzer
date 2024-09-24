package analyzer

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestCheckOversizedInstances(t *testing.T) {
	testCases := []struct {
		name     string
		resource *Resource
		expected *Issue
	}{
		{
			name: "Oversized Instance (t3.large)",
			resource: createResource("aws_instance", "test-instance", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"instance_type": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("t3.large")}},
				},
			}),
			expected: NewIssue("INFO", "Consider downsizing instance 'test-instance' from t3.large to t3.medium", "test-instance", 0),
		},
		{
			name: "Oversized Instance (m5.large)",
			resource: createResource("aws_instance", "test-instance-2", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"instance_type": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("m5.large")}},
				},
			}),
			expected: NewIssue("INFO", "Consider downsizing instance 'test-instance-2' from m5.large to m5.medium", "test-instance-2", 0),
		},
		{
			name: "Not Oversized Instance",
			resource: createResource("aws_instance", "test-instance-3", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"instance_type": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("t2.micro")}},
				},
			}),
			expected: nil, // No issue expected
		},
		{
			name: "Invalid Resource Type",
			resource: createResource("aws_vpc", "test-vpc", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"instance_type": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("t3.large")}}, 
				},
			}),
			expected: nil, // No issue expected as it's not an aws_instance
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issue := CheckOversizedInstances(tc.resource)
			assert.Equal(t, tc.expected, issue)
		})
	}
}

func TestCheckUnattachedEBS(t *testing.T) {
	testCases := []struct {
		name     string
		resource *Resource
		expected *Issue
	}{
		{
			name: "Unattached EBS Volume",
			resource: createResource("aws_ebs_volume", "test-volume", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"id": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("vol-1234567890")}},
				},
			}),
			expected: NewIssue("LOW", "Unattached EBS volume: vol-1234567890", "test-volume", 0),
		},
		{
			name: "Attached EBS Volume",
			resource: createResource("aws_ebs_volume", "test-volume-2", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"id":         {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("vol-abcdefg1234")}},
					"attachment": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("attached")}}, // Simulate attachment
				},
			}),
			expected: nil, // No issue expected
		},
		{
			name: "Invalid Resource Type",
			resource: createResource("aws_s3_bucket", "test-bucket", hcl.BodyContent{
				Attributes: hcl.Attributes{
					"id": {Expr: &hcl.LiteralValueExpr{Val: cty.StringVal("vol-1234567890")}}, 
				},
			}),
			expected: nil, // No issue expected as it's not an aws_ebs_volume
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			issue := CheckUnattachedEBS(tc.resource)
			assert.Equal(t, tc.expected, issue)
		})
	}
}
