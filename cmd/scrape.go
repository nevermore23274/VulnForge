package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"make_dataset/internal/types"

	"github.com/spf13/cobra"
)

type Response struct {
	QueryStatus string        `json:"query_status"`
	Data        []types.IOC   `json:"data"`
}

var scrapeTag string

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape IOCs from ThreatFox and save to a JSONL file",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := loadConfig()

		// Fallback to config if --tag not provided
		if scrapeTag == "" && cfg.DefaultTag != "" {
			scrapeTag = cfg.DefaultTag
		}

		date := time.Now().Format("010206") // MMDDYY
		filename := fmt.Sprintf("output/threatfox_iocs_%s", date)
		if scrapeTag != "" {
			filename += "_" + scrapeTag
		}
		filename += ".jsonl"

		if err := scrapeThreatFox(filename); err != nil {
			fmt.Printf("❌ Failed to scrape: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Scraped IOCs saved to: %s\n", filename)
	},
}

func init() {
	scrapeCmd.Flags().StringVar(&scrapeTag, "tag", "", "Optional tag to append to the output filename")
	AddCommandFunc(scrapeCmd)
}

func scrapeThreatFox(savePath string) error {
	body := []byte(`{"query": "get_iocs"}`)
	resp, err := http.Post("https://threatfox-api.abuse.ch/api/v1/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to request ThreatFox API: %v", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var result Response
	if err := json.Unmarshal(resBody, &result); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	if result.QueryStatus != "ok" {
		return fmt.Errorf("API error: %s", result.QueryStatus)
	}

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	for _, entry := range result.Data {
		j, err := json.Marshal(entry)
		if err != nil {
			continue
		}
		file.Write(j)
		file.Write([]byte("\n"))
	}

	return nil
}
