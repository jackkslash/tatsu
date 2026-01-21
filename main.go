package main

import (
	"fmt"
	"os"

	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/harness"
)

const Version = "0.1.0"

func main() {
	// Parse args
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("❌ Error: task description required")
			printUsage()
			os.Exit(1)
		}
		task := os.Args[2]
		runTask(task)
	case "version", "--version", "-v":
		fmt.Printf("tatsu v%s\n", Version)
	default:
		fmt.Printf("❌ Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runTask(task string) {
	fmt.Printf("🎯 Task: %s\n\n", task)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Configuration loaded successfully")
	fmt.Printf("   Agent: %s\n", cfg.Agent.Command)
	fmt.Printf("   Validate: %s\n\n", cfg.Validate.Command)

	// Check harness availability
	h := harness.NewOpenCodeHarness()
	if !h.IsAvailable() {
		fmt.Printf("❌ %s is not installed or not in PATH\n", h.Name())
		fmt.Println("   Install from: https://github.com/EmbeddedLLM/opencode")
		os.Exit(1)
	}

	fmt.Printf("✅ %s is available\n\n", h.Name())

	fmt.Println("⚠️  Task execution not yet implemented")
}

func printUsage() {
	fmt.Println("tatsu v" + Version)
	fmt.Println("\nUsage:")
	fmt.Println("  tatsu run \"task description\"")
	fmt.Println("  tatsu version")
	fmt.Println("\nExample:")
	fmt.Println("  tatsu run \"add unit tests to the parser\"")
	fmt.Println("\nRequires tatsu.yaml in current directory")
}
