import logging

def check_cost_optimization(resources):
    issues = []
    logging.debug(f"Checking cost optimization for resources: {resources}")

    issue_id = 1  # Unique ID for each issue
    for resource in resources:
        for resource_type, resource_config in resource.items():
            for resource_name, config in resource_config.items():
                if resource_type == 'aws_instance':
                    instance_type = config.get('instance_type')
                    if instance_type in ['t3.large', 'm5.large']:
                        issue = ('INFO', f"Consider downsizing instance '{resource_name}' from {instance_type} to t3.medium", issue_id)
                        issues.append(issue)
                        logging.debug(f"Identified issue: {issue}")
                        issue_id += 1

                elif resource_type == 'aws_ebs_volume':
                    size = config.get('size', 0)
                    if size > 1000:
                        issue = ('LOW', f"Consider resizing EBS volume '{resource_name}' to reduce cost", issue_id)
                        issues.append(issue)
                        logging.debug(f"Identified issue: {issue}")
                        issue_id += 1

                elif resource_type == 'aws_eip':
                    if 'instance' not in config:
                        issue = ('LOW', f"Unattached Elastic IP: {resource_name}", issue_id)
                        issues.append(issue)
                        logging.debug(f"Identified issue: {issue}")
                        issue_id += 1

    logging.debug(f"Final cost optimization issues list: {issues}")
    return issues