package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jack/tatsu/config"
	"github.com/jack/tatsu/harness"
	"github.com/jack/tatsu/prd"
	"github.com/jack/tatsu/runner"
	"github.com/jack/tatsu/tui"
)

const Version = "0.1.0"

const (
	maxIterationsLimit = 100 // Maximum allowed iterations for safety
)

func main() {
	// Parse flags
	maxIterFlag := flag.Int("max-iterations", runner.DefaultMaxIterations, "Maximum number of retry iterations")
	dirFlag := flag.String("C", "", "Run in directory (load tatsu.yaml and run commands there)")
	flag.Parse()

	// Change to target directory if set
	if *dirFlag != "" {
		if err := os.Chdir(*dirFlag); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Validate max iterations
	if *maxIterFlag < 1 {
		fmt.Printf("❌ Error: max-iterations must be at least 1\n")
		os.Exit(1)
	}
	if *maxIterFlag > maxIterationsLimit {
		fmt.Printf("❌ Error: max-iterations cannot exceed %d (got %d)\n", maxIterationsLimit, *maxIterFlag)
		os.Exit(1)
	}

	// Get remaining args after flag parsing
	args := flag.Args()

	// No args: open TUI (everything runs inside TUI; q to quit)
	if len(args) < 1 {
		if _, err := os.Stat("tatsu.yaml"); os.IsNotExist(err) {
			fmt.Println("📝 No tatsu.yaml found. Generating configuration...")
			if err := config.Generate(false); err != nil {
				fmt.Printf("❌ Failed to generate config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✅ Created tatsu.yaml")
		}
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			os.Exit(1)
		}
		h := harness.NewOpenCodeHarness()
		if !h.IsAvailable() {
			fmt.Printf("❌ %s is not installed or not in PATH\n", h.Name())
			os.Exit(1)
		}
		if err := tui.Run(cfg, *maxIterFlag); err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	command := args[0]

	switch command {
	case "run":
		if len(args) < 2 {
			fmt.Println("❌ Error: task description required")
			printUsage()
			os.Exit(1)
		}
		task := args[1]
		runTask(task, *maxIterFlag)
	case "generate", "gen":
		force := false
		if len(args) > 1 {
			for _, arg := range args[1:] {
				if arg == "--force" || arg == "-f" {
					force = true
					break
				}
			}
		}
		generateConfig(force)
	case "prd":
		if len(args) < 2 {
			fmt.Println("❌ Error: PRD file required")
			printUsage()
			os.Exit(1)
		}
		prdFile := args[1]
		runPRD(prdFile, *maxIterFlag)
	case "version", "--version", "-v":
		fmt.Printf("tatsu v%s\n", Version)
	default:
		fmt.Printf("❌ Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runTask(task string, maxIter int) {
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
	fmt.Printf("   Validate: %s\n", cfg.Validate.Command)
	if maxIter != runner.DefaultMaxIterations {
		fmt.Printf("   Max iterations: %d\n", maxIter)
	}
	fmt.Println()

	// Check harness availability
	h := harness.NewOpenCodeHarness()
	if !h.IsAvailable() {
		fmt.Printf("❌ %s is not installed or not in PATH\n", h.Name())
		fmt.Println("   Install from: https://github.com/EmbeddedLLM/opencode")
		os.Exit(1)
	}

	fmt.Printf("✅ %s is available\n\n", h.Name())

	// Run task with runner
	r := runner.NewWithMaxIterations(cfg, h, maxIter)
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

func runPRD(prdFile string, maxIter int) {
	fmt.Printf("📄 Loading PRD: %s\n\n", prdFile)

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

	// Check harness availability
	h := harness.NewOpenCodeHarness()
	if !h.IsAvailable() {
		fmt.Printf("❌ %s is not installed or not in PATH\n", h.Name())
		fmt.Println("   Install from: https://github.com/EmbeddedLLM/opencode")
		os.Exit(1)
	}

	// Load PRD
	prdDoc, err := prd.LoadPRD(prdFile)
	if err != nil {
		fmt.Printf("❌ Failed to load PRD: %v\n", err)
		os.Exit(1)
	}

	// Execute PRD
	r := runner.NewWithMaxIterations(cfg, h, maxIter)
	executor := prd.NewExecutor(r)
	if err := executor.ExecutePRD(prdDoc, prdFile); err != nil {
		fmt.Printf("⚠️  %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("tatsu v" + Version)
	fmt.Println("\nUsage:")
	fmt.Println("  tatsu                          Open TUI (Tab to switch Task / PRD)")
	fmt.Println("  tatsu run \"task description\"  Run a task")
	fmt.Println("  tatsu prd <file>                Execute tasks from PRD file")
	fmt.Println("  tatsu generate [--force]       Generate tatsu.yaml")
	fmt.Println("  tatsu version                  Show version")
	fmt.Println("\nFlags:")
	fmt.Printf("  -max-iterations N              Maximum retry iterations (default: %d, max: %d)\n", runner.DefaultMaxIterations, maxIterationsLimit)
	fmt.Println("\nExamples:")
	fmt.Println("  tatsu run \"add unit tests to the parser\"")
	fmt.Println("  tatsu run -max-iterations 5 \"quick test\"")
	fmt.Println("  tatsu prd PRD.example.md")
	fmt.Println("  tatsu prd -max-iterations 10 PRD.example.md")
	fmt.Println("  tatsu generate")
	fmt.Println("  tatsu generate --force")
	fmt.Println("\nNote: tatsu.yaml will be auto-generated on first run if it doesn't exist")
}
