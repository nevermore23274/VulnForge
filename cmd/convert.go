package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"make_dataset/internal/format"

	"github.com/spf13/cobra"
)

var (
	infile    string
	outfile   string
	limit     int
	formatOpt string
	dryRun    bool
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert scraped IOCs into instruction-tuning dataset format",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := loadConfig() // safe even if config is missing

		if infile == "" {
			fmt.Println("❌ --infile is required.")
			cmd.Help()
			os.Exit(1)
		}

		if formatOpt == "" && cfg.DefaultFormat != "" {
			formatOpt = cfg.DefaultFormat
		}

		if limit == 30 && cfg.DefaultLimit > 0 {
			limit = cfg.DefaultLimit
		}

		if outfile == "" && !dryRun && formatOpt != "" {
			outfile = fmt.Sprintf("output/formatted_output.%s", formatOpt)
		}

		switch strings.ToLower(formatOpt) {
		case "json":
			entries, err := format.ConvertToLLMEntries(infile, limit)
			if err != nil {
				fmt.Printf("❌ Conversion failed: %v\n", err)
				os.Exit(1)
			}

			if dryRun {
				preview, _ := json.MarshalIndent(entries[:min(3, len(entries))], "", "  ")
				fmt.Println("🧪 Dry Run (showing up to 3 samples):")
				fmt.Println(string(preview))
				return
			}

			if err := format.WriteJSON(outfile, entries); err != nil {
				fmt.Printf("❌ Failed to write JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ JSON dataset written to: %s\n", outfile)

		case "csv":
			if dryRun {
				fmt.Println("🧪 Dry Run: CSV format selected (output will not be written).")
				return
			}

			if err := format.ConvertToCSV(infile, outfile, limit); err != nil {
				fmt.Printf("❌ Failed to write CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ CSV dataset written to: %s\n", outfile)

		default:
			fmt.Println("❌ --format must be 'json' or 'csv'")
			os.Exit(1)
		}
	},
}

func init() {
	convertCmd.Flags().StringVar(&infile, "infile", "", "Path to input .jsonl file")
	convertCmd.Flags().StringVar(&outfile, "outfile", "", "Path to save output file (ignored in dry-run)")
	convertCmd.Flags().IntVar(&limit, "limit", 30, "Limit number of samples to convert")
	convertCmd.Flags().StringVar(&formatOpt, "format", "json", "Output format: json or csv")
	convertCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show preview without writing file")

	AddCommandFunc(convertCmd)
}

// utility: ensure we don’t slice beyond list length
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
