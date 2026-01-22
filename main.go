package main

import (
	"fmt"
	"os"

	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/harness"
	"github.com/jack/tatsu/runner"
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
	case "generate", "gen":
		force := false
		if len(os.Args) > 2 {
			for _, arg := range os.Args[2:] {
				if arg == "--force" || arg == "-f" {
					force = true
					break
				}
			}
		}
		generateConfig(force)
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

	// Check if config exists, generate if not
	if _, err := os.Stat("tatsu.yaml"); os.IsNotExist(err) {
		fmt.Println("📝 No tatsu.yaml found. Generating configuration...")
		if err := config.Generate(false); err != nil {
			fmt.Printf("❌ Failed to generate config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Created tatsu.yaml")
		fmt.Println("   (Run 'tatsu generate --force' to regenerate)")
	}

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

	// Run task with runner
	r := runner.New(cfg, h)
	if err := r.Run(task); err != nil {
		fmt.Printf("⚠️  %v\n", err)
		os.Exit(1)
	}
}

func generateConfig(force bool) {
	if force {
		fmt.Println("🔧 Generating tatsu.yaml (overwriting existing file)...")
	} else {
		fmt.Println("🔧 Generating tatsu.yaml...")
	}

	if err := config.Generate(force); err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Created tatsu.yaml")
	fmt.Println("\n📝 Review and update the configuration as needed:")
	fmt.Println("   - agent.command: Your AI agent command")
	fmt.Println("   - validate.command: Your test/validation command")
}

func printUsage() {
	fmt.Println("tatsu v" + Version)
	fmt.Println("\nUsage:")
	fmt.Println("  tatsu run \"task description\"  Run a task")
	fmt.Println("  tatsu generate [--force]       Generate tatsu.yaml")
	fmt.Println("  tatsu version                  Show version")
	fmt.Println("\nExamples:")
	fmt.Println("  tatsu run \"add unit tests to the parser\"")
	fmt.Println("  tatsu generate")
	fmt.Println("  tatsu generate --force")
	fmt.Println("\nNote: tatsu.yaml will be auto-generated on first run if it doesn't exist")
}
