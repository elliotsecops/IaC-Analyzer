import os
import sys
import logging
import argparse
from analyzer import analyze_iac_file
from analyzer.security import check_security_issues
from analyzer.cost import check_cost_optimization

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def main():
    parser = argparse.ArgumentParser(description='IaC Configuration Analyzer')
    parser.add_argument('path', nargs='?', default='test_main.tf', help='Path to Terraform file or directory')
    parser.add_argument('-v', '--verbose', action='store_true', help='Increase output verbosity')
    args = parser.parse_args()

    # Set up logging
    log_level = logging.DEBUG if args.verbose else logging.INFO
    logging.getLogger().setLevel(log_level)

    path = args.path
    checks = [check_security_issues, check_cost_optimization]  # Default checks
    all_issues = []

    if os.path.isdir(path):
        for root, _, files in os.walk(path):
            for file in files:
                if file.endswith('.tf'):
                    file_path = os.path.join(root, file)
                    logging.info(f"Analyzing file: {file_path}")
                    file_issues = analyze_iac_file(file_path, checks)
                    logging.debug(f"Appending issues from file {file_path}: {file_issues}")
                    all_issues.extend(file_issues)
    elif os.path.isfile(path) and path.endswith('.tf'):
        logging.info(f"Analyzing file: {path}")
        file_issues = analyze_iac_file(path, checks)
        logging.debug(f"Appending issues from file {path}: {file_issues}")
        all_issues.extend(file_issues)
    else:
        logging.error("Error: Invalid path or non-Terraform file")
        sys.exit(1)

    if all_issues:
        print_issues(all_issues)
    else:
        print("No issues found.")

def print_issues(issues):
    severity_count = {'HIGH': 0, 'MEDIUM': 0, 'LOW': 0, 'INFO': 0}
    logging.debug(f"Total issues: {len(issues)}")
    
    # Log full list of issues before counting
    logging.debug(f"Issues list: {issues}")

    for i, (severity, message, issue_id) in enumerate(issues):
        # Normalize severity labels to ensure consistency
        severity_normalized = severity.strip().upper()  # Ensure all labels are uppercase and without leading/trailing spaces
        logging.debug(f"Issue #{i + 1}: ID: {issue_id}, Normalized severity: '{severity_normalized}', Message: '{message}'")
        
        # Verify if the severity is valid and count it
        if severity_normalized in severity_count:
            severity_count[severity_normalized] += 1
            logging.debug(f"Counting severity: {severity_normalized}, Current count: {severity_count[severity_normalized]}")
        else:
            logging.warning(f"Unknown severity level: {severity_normalized}")
        
    total = sum(severity_count.values())
    logging.debug(f"Final severity count: {severity_count}")
    print(f"\n{total} issues found ({severity_count['HIGH']} HIGH, {severity_count['MEDIUM']} MEDIUM, {severity_count['LOW']} LOW, {severity_count['INFO']} INFO)")

if __name__ == "__main__":
    main()