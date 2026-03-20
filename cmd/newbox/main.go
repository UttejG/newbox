package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui"
	"github.com/uttejg/newbox/internal/adapter/output/catalogprovider"
	"github.com/uttejg/newbox/internal/adapter/output/checker"
	"github.com/uttejg/newbox/internal/adapter/output/detector"
	"github.com/uttejg/newbox/internal/adapter/output/pkgmgr"
	"github.com/uttejg/newbox/internal/adapter/output/runner"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
	"github.com/uttejg/newbox/internal/core/service"
)

func main() {
	// Find if 'list' is the first non-flag argument, scanning past any flags so
	// that e.g. "newbox --dry-run list" is handled correctly.
	listIdx := -1
	for i, arg := range os.Args[1:] {
		if arg == "list" {
			listIdx = i + 1 // index in os.Args
			break
		}
		if !strings.HasPrefix(arg, "-") {
			break // first non-flag arg is not "list"
		}
	}
	if listIdx >= 0 {
		d := &detector.SystemDetector{}
		platform, err := d.Detect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
			os.Exit(1)
		}
		if err := runList(platform, os.Args[listIdx+1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	dryRun := flag.Bool("dry-run", false, "Simulate installation without executing commands")
	summary := flag.Bool("summary", false, "Print a text summary of the install plan (requires --dry-run)")
	flag.Parse()

	if *summary && !*dryRun {
		fmt.Fprintln(os.Stderr, "Error: --summary requires --dry-run")
		os.Exit(2)
	}

	d := &detector.SystemDetector{}
	platform, err := d.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
		os.Exit(1)
	}

	// Choose CommandRunner based on dry-run flag.
	var cmdRunner port.CommandRunner
	if *dryRun {
		cmdRunner = &runner.DryRunRunner{}
	} else {
		cmdRunner = &runner.ExecRunner{}
	}

	// Wire adapters.
	syschecker := &checker.SystemChecker{Runner: cmdRunner}
	catalogSvc := service.NewCatalogService(&catalogprovider.EmbeddedProvider{})

	// Select package manager based on the detected platform; leave installSvc nil
	// (disabling the install flow) for platforms not yet supported.
	var installSvc port.InstallService
	switch platform.OS {
	case domain.OSMacOS:
		brew := pkgmgr.NewBrew(cmdRunner)
		installSvc = service.NewInstallService(brew, syschecker, *dryRun)
	}

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
		catIDs := make([]string, 0, len(sel.ToolsByCategory))
		for catID := range sel.ToolsByCategory {
			catIDs = append(catIDs, catID)
		}
		sort.Strings(catIDs)
		for _, catID := range catIDs {
			fmt.Printf("  %s:\n", catID)
			for _, t := range sel.ToolsByCategory[catID] {
				fmt.Printf("    • %s\n", t.Name)
			}
		}
	}
}

func runList(platform *domain.Platform, args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	categoryFlag := fs.String("category", "", "Filter to a specific category ID")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	svc := service.NewCatalogService(&catalogprovider.EmbeddedProvider{})
	categories, err := svc.GetCategories(platform.OS)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
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
	return nil
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
