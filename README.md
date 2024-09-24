# IaC Configuration Analyzer (iac-analyzer)

## Description

`iac-analyzer` is a command-line tool written in Go that helps you identify potential security vulnerabilities and cost optimization opportunities in your Terraform infrastructure code. It performs static analysis on your `.tf` files, looking for common misconfigurations and providing actionable recommendations.

## Features

- **Security Checks:**
    - Detects open SSH access in security groups.
    - Identifies publicly accessible S3 buckets.
    - Flags unencrypted sensitive resources (e.g., databases, S3 buckets).
    - Checks for hardcoded secrets (API keys, passwords, etc.).
- **Cost Optimization:**
    - Suggests downsizing oversized EC2 instances.
    - Identifies unattached EBS volumes. 
- **Customizable Configuration:**  
    - Allows you to enable/disable specific checks.
    - Configure custom mappings for oversized instance types.
    - Set the verbosity level of the output.
- **Simple Reporting:**
    - Generates clear and concise reports in text or JSON format.
    - Provides actionable recommendations for resolving issues.
- **CI/CD Integration:** 
    - Uses exit codes to signal analysis results, making it easy to integrate into CI/CD pipelines.

## Installation

1. **Prerequisites:**
   - Go (version 1.13 or later) installed on your system: [https://golang.org/](https://golang.org/)

2. **Install using `go get`:**
   ```bash
   go get github.com/elliotsecops/iac-analyzer
   ```
Create a Configuration File (config.yaml):

See the Configuration section below for options and an example.

Run the Analyzer:

iac-analyzer -config config.yaml -dir path/to/terraform/files


Command-Line Flags:

-config: Path to the configuration file (default: config.yaml).

-dir: Path to the directory containing Terraform files (default: current directory).

-output: Output format (text or json, default: text).

-verbose: Enable verbose mode for more detailed output.

Configuration

The config.yaml file allows you to customize the behavior of the analyzer:

Example config.yaml:

files: 
  - "path/to/terraform/files" # You can specify multiple directories or files.

checks:
  security:
    open_ssh: true
    public_s3: true
    hardcoded_secrets: true
    unencrypted_data: true
    compliance:
      pci: true
      hipaa: false
      gdpr: true
  cost:
    oversized_instances:
      t3.large: t3.medium
      m5.large: m5.medium
      # Add more instance type mappings as needed... 
    underutilized_ebs: true
    unattached_eip: true

output:
  verbosity: medium  # or low, high

Configuration Options:

files: A list of directories or individual .tf files to analyze.

checks.security: Enables/disables specific security checks.

checks.cost: Enables/disables specific cost optimization checks.

checks.cost.oversized_instances: A map of oversized instance types to their recommended alternatives.

output.verbosity: Sets the verbosity level of the output report (low, medium, high).

Exit Codes

0: No issues found.

1: Non-critical issues found (INFO, LOW).

2: Critical issues found (MEDIUM, HIGH).

3: Invalid configuration file format.

4: Error reading the configuration file.

Contributing

Contributions are welcome! If you'd like to contribute to this project, please:

Open an issue to discuss your proposed changes.

Fork the repository and create a pull request with your changes.
