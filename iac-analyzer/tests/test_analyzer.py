import unittest
from analyzer.analyzer import analyze_iac_file
from analyzer.security import check_security_issues
from analyzer.cost import check_cost_optimization

class TestAnalyzer(unittest.TestCase):

    def test_security_check(self):
        issues = analyze_iac_file('tests/test_security.tf', [check_security_issues])
        self.assertTrue(any(issue[0] == 'HIGH' for issue in issues))

    def test_cost_check(self):
        issues = analyze_iac_file('tests/test_cost.tf', [check_cost_optimization])
        self.assertTrue(any(issue[0] == 'INFO' for issue in issues))

if __name__ == '__main__':
    unittest.main()