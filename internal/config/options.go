package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Infile   string
	Outfile  string
	Convert  bool
	Limit    int
}

func ParseFlags() Config {
	now := time.Now()
	dateSuffix := fmt.Sprintf("%02d%02d%02d", now.Month(), now.Day(), now.Year()%100)

	tag := flag.String("tag", "", "Optional tag to append to output files")
	convert := flag.Bool("convert", false, "Skip scraping, convert existing file")
	limit := flag.Int("limit", 30, "Number of samples to include in the output dataset")

	flag.Parse()

	suffix := dateSuffix
	if *tag != "" {
		suffix += "_" + *tag
	}

	infile := fmt.Sprintf("data/threatfox_iocs_%s.jsonl", suffix)
	outfile := fmt.Sprintf("data/formatted_threatfox_instruction_dataset_%s.json", suffix)

	// Ensure the data directory exists
	os.MkdirAll("data", 0755)

	return Config{
		Infile:  infile,
		Outfile: outfile,
		Convert: *convert,
		Limit:   *limit,
	}
}
