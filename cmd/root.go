package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var version = ""

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
	rootCmd.SetVersionTemplate(`{{ printf "malathair-ssm %s\n" .Version }}`)

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
