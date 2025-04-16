package format

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"make_dataset/internal/types"
)

type InstructionData struct {
	Instruction string `json:"instruction"`
	Input       string `json:"input"`
	Output      string `json:"output"`
}

// ConvertToLLMEntries reads a JSONL file and returns instruction-formatted entries
func ConvertToLLMEntries(inputFile string, limit int) ([]InstructionData, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("could not open input file: %v", err)
	}
	defer file.Close()

	var dataset []InstructionData
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		if count >= limit {
			break
		}
		var ioc types.IOC
		if err := json.Unmarshal(scanner.Bytes(), &ioc); err != nil {
			continue
		}

		inputStruct := map[string]interface{}{
			"ioc":              ioc.IOC,
			"ioc_type":         ioc.IOCType,
			"threat_type":      ioc.ThreatType,
			"malware_family":   ioc.MalwareFamily,
			"tags":             ioc.Tags,
			"confidence_level": ioc.ConfidenceLevel,
			"reference":        ioc.Reference,
		}
		inputJSON, _ := json.MarshalIndent(inputStruct, "", "  ")

		tags := "unspecified"
		if len(ioc.Tags) > 0 {
			tags = fmt.Sprintf("%v", ioc.Tags)
		}
		malware := ioc.MalwareFamily
		if malware == "" {
			malware = "an unknown malware family"
		}

		outputText := fmt.Sprintf(
			"This indicator '%s' is classified as a %s.\nIt is potentially associated with %s, based on the tags: %s.\nThe confidence level is %d%%.",
			ioc.IOC, ioc.ThreatType, malware, tags, ioc.ConfidenceLevel)

		dataset = append(dataset, InstructionData{
			Instruction: "Classify this threat indicator and provide context.",
			Input:       string(inputJSON),
			Output:      outputText,
		})
		count++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dataset, nil
}

// WriteJSON writes the instruction data slice to a JSON file
func WriteJSON(outputFile string, data []InstructionData) error {
	outData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode output JSON: %v", err)
	}

	if err := os.WriteFile(outputFile, outData, 0644); err != nil {
		return fmt.Errorf("could not write output file: %v", err)
	}

	return nil
}

// ConvertToCSV writes a limited number of entries as CSV
func ConvertToCSV(inputFile, outputFile string, limit int) error {
	entries, err := ConvertToLLMEntries(inputFile, limit)
	if err != nil {
		return err
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("could not create output CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	writer.Write([]string{"instruction", "input", "output"})

	for _, entry := range entries {
		writer.Write([]string{entry.Instruction, entry.Input, entry.Output})
	}

	return nil
}
