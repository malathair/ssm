package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	version = "v0.0.0"

	config RuntimeConfig
)

var rootCmd = &cobra.Command{
	Version: version,
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if config.dryRun {
			printConfig()
		}

		err := ssm(args[0])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	},
}

type RuntimeConfig struct {
	remoteCmd  string
	debugLvl   int
	domains    []string
	dryRun     bool
	jump       bool
	jumphost   string
	port       int
	tunnel     bool
	tunnelPort int
}

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {
	// Set custom help/usage and version templates
	cobra.AddTemplateFunc(
		"wrappedFlagUsage",
		func(cmd *pflag.FlagSet) string {
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				width = 80
			}

			return cmd.FlagUsagesWrapped(width - 1)
		},
	)
	rootCmd.SetUsageTemplate(
		"Usage: ssm [OPTION]... [user@]host\n\n" +
			"{{if .HasAvailableLocalFlags}}" +
			"{{wrappedFlagUsage .LocalFlags | trimTrailingWhitespaces}}" +
			"{{end}}\n",
	)
	rootCmd.SetVersionTemplate(`{{ printf "malathair-ssm v%s\n" .Version }}`)

	// Load default values from config file so we can pass them as defaults to flags
	loadConfig()

	// Flags
	rootCmd.Flags().StringVarP(
		&config.remoteCmd,
		"command",
		"c",
		"",
		"Execute the specified command on the remote system without opening an interactive shell. "+
			"The connection will be terminated immediately after command executes",
	)
	rootCmd.Flags().BoolVar(
		&config.dryRun,
		"dry-run",
		viper.GetBool("flags.dry-run"),
		"Do a dry run and print the raw SSH command that would be executed for debugging",
	)
	rootCmd.Flags().BoolVarP(
		&config.jump,
		"jump",
		"j",
		viper.GetBool("flags.jump"),
		"Use a jumphost to reach the remote host",
	)
	rootCmd.Flags().StringVarP(
		&config.jumphost,
		"jumphost",
		"J",
		viper.GetString("flags.jumphost"),
		"Overrides the jump host specified in the configuration file",
	)
	rootCmd.Flags().IntVarP(
		&config.port,
		"port",
		"p",
		viper.GetInt("flags.port"),
		"Specifies the port to use for the SSH session",
	)
	rootCmd.Flags().BoolVarP(
		&config.tunnel,
		"tunnel",
		"t",
		false,
		"Start a SOCKS5 tunnel on the port defined in the configuration file. "+
			"You may then use a SOCKS5 proxy config in your browser, or a SOCKS5 proxy client like "+
			"tsocks or proxychains to proxy tcp traffic through the SOCKS5 tunnel. Has no effect if "+
			"the --command flag is specified",
	)
	rootCmd.Flags().IntVarP(
		&config.tunnelPort,
		"tunnel-port",
		"T",
		viper.GetInt("flags.tunnel-port"),
		"Specifies the local port to bind the SOCKS5 tunnel to",
	)
	rootCmd.Flags().CountVarP(
		&config.debugLvl,
		"verbose",
		"v",
		"Print verbose debug messages about the SSH connection. "+
			"Multiple -v options increase the verbosity up to a maximum of 3 (-vvv)",
	)

	rootCmd.MarkFlagsMutuallyExclusive("jump", "jumphost")
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
	viper.SetDefault("domains", []string{})
	viper.SetDefault("debug-level", 0)
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

	// Load domains manually since there isn't a flag for this
	config.domains = viper.GetStringSlice("domains")
	// Load debug level manually since the flag can't bind a default
	config.debugLvl = viper.GetInt("debug-level")
}

func printConfig() {
	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)

	if config.remoteCmd == "" {
		fmt.Fprintf(writer, "Command\tn/a\n")
	} else {
		fmt.Fprintf(writer, "Command\t%v\n", config.remoteCmd)
	}
	fmt.Fprintf(writer, "Debug Level\t%v\n", config.debugLvl)
	fmt.Fprintf(writer, "Domains\t%v\n", config.domains)
	fmt.Fprintf(writer, "DryRun\t%v\n", config.dryRun)
	fmt.Fprintf(writer, "Jump\t%v\n", config.jump)
	if config.jumphost == "" {
		fmt.Fprintf(writer, "Jumphost\tn/a\n")
	} else {
		fmt.Fprintf(writer, "Jumphost\t%v\n", config.jumphost)
	}
	fmt.Fprintf(writer, "Port\t%v\n", config.port)
	fmt.Fprintf(writer, "Tunnel\t%v\n", config.tunnel)
	fmt.Fprintf(writer, "Tunnel Port\t%v\n", config.tunnelPort)

	writer.Flush()

	fmt.Println()
}
