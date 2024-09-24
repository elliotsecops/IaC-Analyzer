package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"iac-analyzer/analyzer"

	"github.com/fatih/color"
)

const (
	exitCodeOK                  = iota
	exitCodeIssuesFound
	exitCodeCriticalIssuesFound
	exitCodeInvalidConfigFormat
	exitCodeConfigLoadError
)

func main() {
	// Define command-line flags
	configFile := flag.String("config", "config.yaml", "Path to the configuration file")
	terraformDir := flag.String("dir", ".", "Path to the directory containing Terraform files")
	outputFormat := flag.String("output", "text", "Output format (text or json)")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	flag.Parse()

	if *verbose {
		log.Println("Verbose mode enabled")
		log.Printf("Loading configuration from %s", *configFile)
	}

	// Load configuration
	config, err := analyzer.LoadConfig(*configFile)
	if err != nil {
		log.Printf("Error loading configuration from %s: %v", *configFile, err)
		os.Exit(exitCodeConfigLoadError)
	}

	// Initialize the analyzer
	analyzerInstance := analyzer.NewAnalyzer(config, *verbose) // Pass verbose flag to Analyzer

	// Analyze the Terraform files
	results, err := analyzerInstance.Analyze(*terraformDir)
	if err != nil {
		log.Fatalf("Failed to analyze Terraform files in %s: %v", *terraformDir, err)
	}

	// Generate and print the report
	var report string
	switch *outputFormat {
	case "text":
		report = generateTextReport(results) // Renamed for clarity
	case "json":
		report = generateJSONReport(results)
	default:
		log.Fatalf("Unsupported output format: %s", *outputFormat)
	}
	fmt.Println(report)

	// Set the exit code based on the analysis results
	exitCode := getExitCode(results)
	os.Exit(exitCode)
}

// generateTextReport generates a colorized text report
func generateTextReport(results *analyzer.AnalysisResults) string {
	report := "Security Issues:\n"
	for _, issue := range results.SecurityIssues {
		switch issue.Severity {
		case "HIGH":
			report += color.RedString("- [%s] %s\n", issue.Severity, issue.Description)
		case "MEDIUM":
			report += color.YellowString("- [%s] %s\n", issue.Severity, issue.Description)
		default:
			report += fmt.Sprintf("- [%s] %s\n", issue.Severity, issue.Description)
		}
	}

	report += "\nCost Optimization:\n"
	for _, suggestion := range results.CostSuggestions {
		switch suggestion.Severity {
		case "INFO":
			report += color.CyanString("- [%s] %s\n", suggestion.Severity, suggestion.Description)
		case "LOW":
			report += color.GreenString("- [%s] %s\n", suggestion.Severity, suggestion.Description)
		default:
			report += fmt.Sprintf("- [%s] %s\n", suggestion.Severity, suggestion.Description)
		}
	}

	report += fmt.Sprintf("\n%d security issues found (%d HIGH, %d MEDIUM)\n",
		len(results.SecurityIssues), countSeverity(results.SecurityIssues, "HIGH"), countSeverity(results.SecurityIssues, "MEDIUM"))
	report += fmt.Sprintf("%d cost optimization suggestions found (%d INFO, %d LOW)\n",
		len(results.CostSuggestions), countSeverity(results.CostSuggestions, "INFO"), countSeverity(results.CostSuggestions, "LOW"))

	return report
}

// generateJSONReport generates a JSON formatted report
func generateJSONReport(results *analyzer.AnalysisResults) string {
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error serializing results to JSON: %v", err)
	}
	return string(jsonData)
}

// countSeverity counts the number of issues with a specific severity
func countSeverity(issues []analyzer.Issue, severity string) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == severity {
			count++
		}
	}
	return count
}

// getExitCode determines the appropriate exit code based on the analysis results
func getExitCode(results *analyzer.AnalysisResults) int {
	if len(results.SecurityIssues) > 0 {
		for _, issue := range results.SecurityIssues {
			if issue.Severity == "HIGH" {
				return exitCodeCriticalIssuesFound
			}
		}
		return exitCodeIssuesFound
	}
	return exitCodeOK
}
