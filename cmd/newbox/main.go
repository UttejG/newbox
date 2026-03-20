package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uttejg/newbox/internal/adapter/input/tui"
	"github.com/uttejg/newbox/internal/adapter/output/catalogprovider"
	"github.com/uttejg/newbox/internal/adapter/output/detector"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/service"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Simulate installation without executing commands")
	flag.Parse()

	d := &detector.SystemDetector{}
	platform, err := d.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) > 0 && args[0] == "list" {
		if err := runList(platform, args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Launch TUI
	catalogSvc := service.NewCatalogService(&catalogprovider.EmbeddedProvider{})
	app := tui.NewApp(platform, catalogSvc, *dryRun)

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
