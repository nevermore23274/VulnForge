package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	DefaultFormat string `json:"default_format"`
	DefaultTag    string `json:"default_tag"`
	DefaultLimit  int    `json:"default_limit"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI preferences like default format or tag",
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := loadConfig()
		if err != nil {
			fmt.Println("⚠️ No config found or failed to load.")
			return
		}
		out, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println(string(out))
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a config value (e.g., default_format csv)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		cfg, _ := loadConfig()

		switch key {
		case "default_format":
			cfg.DefaultFormat = value
		case "default_tag":
			cfg.DefaultTag = value
		case "default_limit":
			fmt.Sscanf(value, "%d", &cfg.DefaultLimit)
		default:
			fmt.Println("❌ Unknown config key.")
			return
		}

		if err := saveConfig(cfg); err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			return
		}
		fmt.Println("✅ Config updated.")
	},
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".threatfox", "config.json")
}

func loadConfig() (Config, error) {
	var cfg Config
	path := configPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func saveConfig(cfg Config) error {
	path := configPath()
	os.MkdirAll(filepath.Dir(path), 0755)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clear all saved preferences",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(configPath())
		if err != nil {
			fmt.Println("⚠️ Could not remove config file:", err)
			return
		}
		fmt.Println("✅ Config reset.")
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	AddCommandFunc(configCmd)
}
