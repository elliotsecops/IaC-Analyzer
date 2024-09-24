package analyzer

import (
	"fmt"
	"log"
	"path/filepath"
)

// Analyzer is the core struct for performing analysis.
type Analyzer struct {
	Config   *Config
	Verbose  bool
	analyzers []AnalyzerFunc // Functions to perform analysis
}

// AnalyzerFunc represents a function that performs analysis on resources.
type AnalyzerFunc func([]*Resource) ([]Issue, error)

// Issue represents a generic issue (security or cost).
type Issue struct {
	Severity    string
	Description string
	FilePath    string // Added field for file path
	LineNumber  int    // Added field for line number
}

// NewIssue creates a new Issue with the given details.
func NewIssue(severity, description, filePath string, lineNumber int) *Issue {
	return &Issue{
		Severity:    severity,
		Description: description,
		FilePath:    filePath,
		LineNumber:  lineNumber,
	}
}

// NewAnalyzer creates a new Analyzer instance.
func NewAnalyzer(config *Config, verbose bool) *Analyzer {
	return &Analyzer{
		Config:  config,
		Verbose: verbose,
		analyzers: []AnalyzerFunc{
			// Add your analyzer functions here
			CheckSecurity,
			CheckCost,
		},
	}
}

// Analyze analyzes the Terraform files in the given directory.
func (a *Analyzer) Analyze(terraformDir string) (*AnalysisResults, error) {
	// 1. Find all Terraform files in the directory
	files, err := findTerraformFiles(terraformDir)
	if err != nil {
		return nil, fmt.Errorf("error finding Terraform files: %w", err)
	}

	if a.Verbose {
		log.Printf("Found %d Terraform files to analyze.\n", len(files))
	}

	// 2. Parse each Terraform file
	var allResources []*Resource
	for _, file := range files {
		if a.Verbose {
			log.Printf("Parsing: %s\n", file)
		}

		resources, err := ParseTerraformFiles(file) // Use your function name here
		if err != nil {
			return nil, fmt.Errorf("error parsing HCL file '%s': %w", file, err)
		}
		allResources = append(allResources, resources...)
	}

	// 3. Perform analysis using the registered analyzers
	results := &AnalysisResults{}
	for _, analyzeFunc := range a.analyzers {
		issues, err := analyzeFunc(allResources)
		if err != nil {
			return nil, fmt.Errorf("error during analysis: %w", err)
		}
		results.categorizeIssues(issues)
	}

	return results, nil
}

// AnalysisResults holds the results of the analysis.
type AnalysisResults struct {
	SecurityIssues  []Issue
	CostSuggestions []Issue
}

// categorizeIssues classifies issues into security and cost categories.
func (r *AnalysisResults) categorizeIssues(issues []Issue) {
	for _, issue := range issues {
		if issue.Severity == "HIGH" || issue.Severity == "MEDIUM" {
			r.SecurityIssues = append(r.SecurityIssues, issue)
		} else if issue.Severity == "INFO" || issue.Severity == "LOW" {
			r.CostSuggestions = append(r.CostSuggestions, issue)
		}
	}
}

// findTerraformFiles finds all .tf files in the given directory.
func findTerraformFiles(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*.tf")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("error finding files: %w", err)
	}
	return files, nil
}

// CheckSecurity performs security checks on the given resources.
func CheckSecurity(resources []*Resource) ([]Issue, error) {
	var issues []Issue
	for _, resource := range resources {
		// Implement security checks here
		// Example:
		// if someSecurityProblem {
		//     issues = append(issues, NewIssue("HIGH", "Description of the problem", filePath, lineNumber))
		// }
	}
	return issues, nil
}

// CheckCost performs cost optimization checks on the given resources.
func CheckCost(resources []*Resource) ([]Issue, error) {
	var issues []Issue
	for _, resource := range resources {
		// Implement cost optimization checks here
		// Example:
		// if someCostProblem {
		//     issues = append(issues, NewIssue("INFO", "Description of the problem", filePath, lineNumber))
		// }
	}
	return issues, nil
}
