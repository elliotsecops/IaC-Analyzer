package analyzer

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// CostCheck represents a cost optimization check function.
type CostCheck func(*Resource) *Issue

// CheckOversizedInstances checks for oversized EC2 instances.
func CheckOversizedInstances(resource *Resource, config *Config) *Issue {
	if resource.Type == "aws_instance" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "instance_type"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing instance, invalid 'instance_type' attribute: %v", diags), resource.Name, 0)
		}

		instanceType, _ := bodyContent.Attributes["instance_type"].Expr.Value(nil)
		if recommendedType, exists := config.Cost.OversizedInstances[instanceType.AsString()]; exists {
			return NewIssue("INFO", fmt.Sprintf("Consider downsizing instance '%s' from %s to %s", resource.Name, instanceType.AsString(), recommendedType), resource.Name, bodyContent.Attributes["instance_type"].Expr.Range().Start.Line)
		}
	}
	return nil
}

// CheckUnattachedEBS checks for unattached EBS volumes.
func CheckUnattachedEBS(resource *Resource) *Issue {
	if resource.Type == "aws_ebs_volume" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "id"},
				{Name: "attachment"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing EBS volume, invalid 'id' or 'attachment' attribute: %v", diags), resource.Name, 0)
		}

		if _, exists := bodyContent.Attributes["attachment"]; !exists {
			id, _ := bodyContent.Attributes["id"].Expr.Value(nil)
			return NewIssue("LOW", fmt.Sprintf("Unattached EBS volume: %s", id.AsString()), resource.Name, bodyContent.Attributes["id"].Expr.Range().Start.Line)
		}
	}
	return nil
}

// CheckUnusedElasticIPs checks for unused Elastic IPs.
func CheckUnusedElasticIPs(resource *Resource) *Issue {
	if resource.Type == "aws_eip" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "instance"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing Elastic IP: %v", diags), resource.Name, 0)
		}

		if _, exists := bodyContent.Attributes["instance"]; !exists {
			return NewIssue("LOW", fmt.Sprintf("Unused Elastic IP: %s", resource.Name), resource.Name, 0)
		}
	}
	return nil
}

// CheckIdleLoadBalancers checks for idle load balancers.
func CheckIdleLoadBalancers(resource *Resource) *Issue {
	if resource.Type == "aws_lb" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "load_balancer_type"},
				{Name: "subnets"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing Load Balancer: %v", diags), resource.Name, 0)
		}

		loadBalancerType, _ := bodyContent.Attributes["load_balancer_type"].Expr.Value(nil)
		if loadBalancerType.Equals(cty.StringVal("application")) {
			subnets, _ := bodyContent.Attributes["subnets"].Expr.Value(nil)
			if subnets.LengthInt() == 0 {
				return NewIssue("LOW", fmt.Sprintf("Idle Load Balancer: %s", resource.Name), resource.Name, 0)
			}
		}
	}
	return nil
}

// CheckUnderutilizedRDSInstances checks for underutilized RDS instances.
func CheckUnderutilizedRDSInstances(resource *Resource) *Issue {
	if resource.Type == "aws_db_instance" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "instance_class"},
				{Name: "allocated_storage"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing RDS instance: %v", diags), resource.Name, 0)
		}

		instanceClass, _ := bodyContent.Attributes["instance_class"].Expr.Value(nil)
		allocatedStorage, _ := bodyContent.Attributes["allocated_storage"].Expr.Value(nil)

		if instanceClass.Equals(cty.StringVal("db.t3.medium")) && allocatedStorage.LessThan(cty.NumberIntVal(20)) {
			return NewIssue("LOW", fmt.Sprintf("Underutilized RDS instance: %s", resource.Name), resource.Name, 0)
		}
	}
	return nil
}

// CheckExpiredSSLCertificates checks for expired SSL certificates.
func CheckExpiredSSLCertificates(resource *Resource) *Issue {
	if resource.Type == "aws_acm_certificate" {
		bodyContent, _, diags := resource.Body.PartialContent(&hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "expiration_days"},
			},
		})
		if diags.HasErrors() {
			return NewIssue("ERROR", fmt.Sprintf("Error parsing SSL certificate: %v", diags), resource.Name, 0)
		}

		expirationDays, _ := bodyContent.Attributes["expiration_days"].Expr.Value(nil)
		if expirationDays.LessThan(cty.NumberIntVal(30)) {
			return NewIssue("LOW", fmt.Sprintf("Expired SSL certificate: %s", resource.Name), resource.Name, 0)
		}
	}
	return nil
}

// RunCostChecks runs all cost optimization checks on the given resources.
func RunCostChecks(resources []*Resource, config *Config) []Issue {
	var issues []Issue
	checks := []CostCheck{
		func(r *Resource) *Issue { return CheckOversizedInstances(r, config) },
		CheckUnattachedEBS,
		CheckUnusedElasticIPs,
		CheckIdleLoadBalancers,
		CheckUnderutilizedRDSInstances,
		CheckExpiredSSLCertificates,
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
