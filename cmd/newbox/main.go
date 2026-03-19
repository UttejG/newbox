package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/uttejg/newbox/internal/adapter/output/catalogprovider"
	"github.com/uttejg/newbox/internal/adapter/output/detector"
	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/service"
)

func main() {
	d := &detector.SystemDetector{}
	platform, err := d.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting platform: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 && os.Args[1] == "list" {
		if err := runList(platform, os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println(detector.FormatDetectionInfo(platform))
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
