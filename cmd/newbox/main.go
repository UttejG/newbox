package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui"
	"github.com/uttejg/newbox/internal/adapter/output/catalogprovider"
	"github.com/uttejg/newbox/internal/adapter/output/checker"
	"github.com/uttejg/newbox/internal/adapter/output/detector"
	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/adapter/output/runner"
	"github.com/uttejg/newbox/internal/adapter/output/statestore"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/core/service"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Simulate installation without executing commands")
	summary := flag.Bool("summary", false, "Print a text summary of the install plan (requires --dry-run)")
	flag.Parse()

	d := &detector.SystemDetector{}
	platform, err := d.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) > 0 && args[0] == "list" {
		runList(platform, args[1:])
		return
	}

	// Choose CommandRunner based on dry-run flag.
	var cmdRunner port.CommandRunner
	dryRunner := &runner.DryRunRunner{}
	if *dryRun {
		cmdRunner = dryRunner
	} else {
		cmdRunner = &runner.ExecRunner{}
	}

	// Wire state store for resume support.
	store, err := statestore.NewFileStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not init state store: %v\n", err)
	}

	// Offer resume if a previous install was interrupted.
	if store != nil && store.Exists() {
		savedState, _ := store.Load()
		if savedState != nil && len(savedState.CompletedIDs) > 0 {
			fmt.Fprintf(os.Stderr, "\nPrevious install found (%d tools completed). Resume? [Y/n]: ",
				len(savedState.CompletedIDs))
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := strings.TrimSpace(scanner.Text())
			if strings.EqualFold(answer, "n") {
				_ = store.Clear()
			}
		}
	}

	// Wire adapters.
	pkgManager := pkgmgr.NewForPlatform(platform, cmdRunner)
	syschecker := &checker.SystemChecker{Runner: cmdRunner}
	installSvc := service.NewInstallService(pkgManager, syschecker, store, *dryRun)
	catalogSvc := service.NewCatalogService(&catalogprovider.EmbeddedProvider{})

	// Non-interactive summary mode.
	if *dryRun && *summary {
		runSummary(platform, catalogSvc, installSvc)
		return
	}

	// Launch TUI.
	app := tui.NewApp(platform, catalogSvc, installSvc, *dryRun)

	p := tea.NewProgram(app, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	appModel := finalModel.(*tui.AppModel)
	sel := appModel.FinalSelection()
	if sel != nil && sel.TotalCount() > 0 {
		fmt.Printf("\nSelected %d tools:\n", sel.TotalCount())
		for catID, tools := range sel.ToolsByCategory {
			fmt.Printf("  %s:\n", catID)
			for _, t := range tools {
				fmt.Printf("    • %s\n", t.Name)
			}
		}
	}
}

// runSummary prints a dry-run install plan for all tools on the current platform.
func runSummary(platform *domain.Platform, catalogSvc port.CatalogService, installSvc port.InstallService) {
	cats, err := catalogSvc.GetCategories(platform.OS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading catalog: %v\n", err)
		os.Exit(1)
	}

	byCategory := make(map[string][]domain.Tool, len(cats))
	for _, cat := range cats {
		byCategory[cat.ID] = cat.Tools
	}
	selection := &domain.UserSelection{
		Platform:        platform,
		ToolsByCategory: byCategory,
	}

	plan, err := installSvc.Plan(context.Background(), selection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building plan: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Dry-run install plan for %s\n\n", platform.Summary())
	for _, step := range plan.Steps {
		switch step.Status {
		case domain.StatusDryRun:
			fmt.Printf("  [ dry-run ] %s\n", step.Command)
		case domain.StatusSkipped:
			fmt.Printf("  [ skip    ] %s  (already installed)\n", step.Command)
		}
	}
	pending := plan.PendingSteps()
	skipped := plan.SkippedSteps()
	fmt.Printf("\n  Total: %d would install, %d would skip\n", len(pending), len(skipped))
}

func runList(platform *domain.Platform, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	categoryFlag := fs.String("category", "", "Filter to a specific category ID")
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	svc := service.NewCatalogService(&catalogprovider.EmbeddedProvider{})
	categories, err := svc.GetCategories(platform.OS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading catalog: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Available tools on %s\n\n", platform.Summary())

	for _, cat := range categories {
		if *categoryFlag != "" && cat.ID != *categoryFlag {
			continue
		}

		fmt.Printf("%s\n", cat.Name)
		fmt.Printf("%s\n", strings.Repeat("─", 40))
		for _, tool := range cat.Tools {
			marker := ""
			if tool.DotfilesDefault {
				marker = " ★"
			}
			fmt.Printf("  %-20s %s%s\n", tool.Name, tool.Description, marker)
		}
		fmt.Println()
	}
}



