import logging

def check_security_issues(resources):
    issues = []
    logging.debug(f"Checking security issues for resources: {resources}")

    issue_id = 1  # Unique ID for each issue
    for resource in resources:
        for resource_type, resource_config in resource.items():
            for resource_name, config in resource_config.items():
                if resource_type == 'aws_security_group':
                    ingress_rules = config.get('ingress', [])
                    for rule in ingress_rules:
                        if rule.get('cidr_blocks') == ['0.0.0.0/0'] and rule.get('to_port') == 22:
                            issue = ('HIGH', f"Security issue: Open SSH access in security group {resource_name}", issue_id)
                            issues.append(issue)
                            logging.debug(f"Identified issue: {issue}")
                            issue_id += 1

                elif resource_type == 'aws_s3_bucket':
                    if config.get('acl') == 'public-read':
                        issue = ('MEDIUM', f"Public read access enabled on S3 bucket '{config.get('bucket', 'Unknown')}'", issue_id)
                        issues.append(issue)
                        logging.debug(f"Identified issue: {issue}")
                        issue_id += 1

                elif resource_type == 'aws_db_instance':
                    if 'storage_encrypted' not in config or not config['storage_encrypted']:
                        issue = ('HIGH', f"Security issue: Unencrypted RDS instance {resource_name}", issue_id)
                        issues.append(issue)
                        logging.debug(f"Identified issue: {issue}")
                        issue_id += 1

    logging.debug(f"Final security issues list: {issues}")
    return issues