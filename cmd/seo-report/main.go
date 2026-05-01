package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jaychinthrajah/seo-report/internal/config"
	"github.com/jaychinthrajah/seo-report/internal/fetcher"
	"github.com/jaychinthrajah/seo-report/internal/processor"
	"github.com/jaychinthrajah/seo-report/internal/renderer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seo-report",
	Short: "Generate daily AEO/SEO news digest as a GitHub Pages static site",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate today's report and regenerate the index",
	RunE:  runGenerate,
}

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "List configured RSS sources",
	RunE:  runSources,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	dateStr, _ := cmd.Flags().GetString("date")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	var reportDate time.Time
	if dateStr != "" {
		reportDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("parsing date %q: %w", dateStr, err)
		}
	} else {
		reportDate = time.Now().UTC().Truncate(24 * time.Hour)
	}

	since := reportDate.Add(-24 * time.Hour)
	fmt.Printf("Fetching items since %s...\n", since.Format("2006-01-02 15:04 UTC"))

	items, errs := fetcher.FetchAll(cfg, since)
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "fetch warning: %v\n", e)
	}
	fmt.Printf("Fetched %d raw items from %d sources\n", len(items), len(cfg.Sources.RSS))

	fmt.Println("Classifying items...")
	classified := processor.Classify(items, cfg)

	// Build ordered categories for report
	cats := renderer.BuildReport(classified, cfg.Categories)


	// Collect source names
	sourceNames := make([]string, 0, len(cfg.Sources.RSS))
	for _, s := range cfg.Sources.RSS {
		sourceNames = append(sourceNames, s.Name)
	}

	generatedAt := time.Now().UTC()
	reportHTML := renderer.RenderReport(reportDate, cats, sourceNames, generatedAt)

	if dryRun {
		fmt.Println("--- DRY RUN OUTPUT ---")
		fmt.Println(reportHTML)
		return nil
	}

	outDir := cfg.Output.Dir
	if outDir == "" {
		outDir = "docs"
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Write per-day report
	reportFilename := reportDate.Format("2006-01-02") + ".html"
	reportPath := filepath.Join(outDir, reportFilename)
	if err := os.WriteFile(reportPath, []byte(reportHTML), 0644); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}
	fmt.Printf("Report written to %s\n", reportPath)

	// Regenerate index
	indexMetas, err := buildIndexMetas(outDir)
	if err != nil {
		return fmt.Errorf("building index: %w", err)
	}
	indexHTML := renderer.RenderIndex(indexMetas, generatedAt)
	indexPath := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(indexHTML), 0644); err != nil {
		return fmt.Errorf("writing index: %w", err)
	}
	fmt.Printf("Index written to %s\n", indexPath)

	return nil
}

// buildIndexMetas scans the output dir for YYYY-MM-DD.html files and builds report metadata.
func buildIndexMetas(outDir string) ([]renderer.ReportMeta, error) {
	entries, err := os.ReadDir(outDir)
	if err != nil {
		return nil, err
	}

	var metas []renderer.ReportMeta
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) != len("2006-01-02.html") || name[10:] != ".html" {
			continue
		}
		dateStr := name[:10]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		// We don't parse item/category counts from the file; use 0 as placeholder.
		// The full pipeline populates this for the current day report; for history we use the file.
		metas = append(metas, renderer.ReportMeta{
			Date:          date,
			Filename:      name,
			ItemCount:     0,
			CategoryCount: 0,
		})
	}

	// Sort newest first
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Date.After(metas[j].Date)
	})

	// Populate current day counts if available (last generated report)
	return metas, nil
}

func runSources(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	fmt.Printf("%-40s %s\n", "Name", "Weight")
	fmt.Printf("%-40s %s\n", "----", "------")
	for _, s := range cfg.Sources.RSS {
		fmt.Printf("%-40s %d\n", s.Name, s.Weight)
	}
	fmt.Printf("\nTotal: %d sources\n", len(cfg.Sources.RSS))
	return nil
}

func init() {
	generateCmd.Flags().String("config", "config.yaml", "Path to config file")
	generateCmd.Flags().String("date", "", "Generate report for specific date (YYYY-MM-DD)")
	generateCmd.Flags().Bool("dry-run", false, "Print HTML to stdout, no file written")

	sourcesCmd.Flags().String("config", "config.yaml", "Path to config file")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(sourcesCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
