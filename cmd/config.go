package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config RuntimeConfig

type RuntimeConfig struct {
	remoteCmd  string
	debugLvl   int
	domains    []string
	dryRun     bool
	jump       bool
	jumphost   string
	port       int
	sshpass    bool
	tunnel     bool
	tunnelPort int
}

func loadConfig() {
	// Set config paths
	viper.AddConfigPath(".")
	if runtime.GOOS == "windows" {
		// Windows is weird so we'll use %USERPROFILE%
		// (possibly should be %USERPROFILE%\AppData\Local\ssm\ instead)
		homeDir, _ := os.UserHomeDir()
		viper.AddConfigPath(homeDir)
	} else {
		// Otherwise adhere to the XDG spec and use XDG_CONFIG_HOME
		configDir, _ := os.UserConfigDir()
		viper.AddConfigPath(configDir)
	}

	viper.SetConfigName("ssm.conf")
	viper.SetConfigType("toml")

	// Set some sensible default values in case loading a config fails
	viper.SetDefault("debug-level", 0)
	viper.SetDefault("domains", []string{})
	viper.SetDefault("sshpass", false)
	viper.SetDefault("flags.dry-run", false)
	viper.SetDefault("flags.jump", false)
	viper.SetDefault("flags.jumphost", "")
	viper.SetDefault("flags.port", 22)
	viper.SetDefault("flags.tunnel", false)
	viper.SetDefault("flags.tunnel-port", 6060)

	err := viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Do nothing here. We'll just fall back to defaults
		} else {
			cobra.CheckErr(err)
		}
	}

	// Load defaults that don't have flags or for flags that can't bind defaults
	config.debugLvl = viper.GetInt("debug-level")
	config.domains = viper.GetStringSlice("domains")
	config.sshpass = viper.GetBool("sshpass")
}

func printConfig() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

	{
		fmt.Println()
		fmt.Fprintf(writer, "Parameter\tValue\n")
		fmt.Fprintf(writer, "---------\t-----\n")

		fmt.Fprintf(writer, "Debug Level\t%v\n", config.debugLvl)
		fmt.Fprintf(writer, "Domains\t%v\n", config.domains)
		fmt.Fprintf(writer, "DryRun\t%v\n", config.dryRun)
		fmt.Fprintf(writer, "Jump\t%v\n", config.jump)
		fmt.Fprintf(writer, "Jumphost\t%v\n", config.jumphost)
		fmt.Fprintf(writer, "Port\t%v\n", config.port)
		fmt.Fprintf(writer, "Remote Cmd\t%v\n", config.remoteCmd)
		fmt.Fprintf(writer, "Sshpass\t%v\n", config.sshpass)
		fmt.Fprintf(writer, "Tunnel\t%v\n", config.tunnel)
		fmt.Fprintf(writer, "Tunnel Port\t%v\n", config.tunnelPort)
	}

	writer.Flush()

	fmt.Println()
}
