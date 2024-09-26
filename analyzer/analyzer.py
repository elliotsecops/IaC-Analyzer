from .parser import parse_hcl_file
from .security import check_security_issues
from .cost import check_cost_optimization

def analyze_iac_file(file_path, checks):
    resources = parse_hcl_file(file_path)
    print(f"Parsed resources: {resources}")  # Debug print
    issues = []
    for check in checks:
        check_issues = check(resources)
        print(f"Issues from {check.__name__}: {check_issues}")  # Debug print
        issues.extend(check_issues)
    return issues