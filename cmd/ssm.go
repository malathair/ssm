package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func buildFqdn(host string) string {
	ip := net.ParseIP(host)
	if ip != nil {
		if config.dryRun {
			fmt.Printf("%s is valid IP address\n", host)
			fmt.Println()
		}
		return ip.String()
	} else if config.dryRun {
		fmt.Printf("%s is not a valid IP Address\n", host)
	}

	_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", host)
	if err == nil {
		if config.dryRun {
			fmt.Printf("%s is a valid FQDN\n", host)
			fmt.Println()
		}
		return host
	} else if config.dryRun {
		fmt.Printf("%s is not a valid FQDN (%s)\n", host, err)
	}

	if len(config.domains) < 1 {
		fmt.Println("No domains found for FQDN autocompletion")
		return ""
	}

	if config.dryRun {
		fmt.Println("Attempting to build FQDN from domain list:")
	}
	for _, domain := range config.domains {
		fqdn := host + "." + domain
		_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", fqdn)
		if err == nil {
			if config.dryRun {
				fmt.Printf("%s is a valid FQDN\n", fqdn)
				fmt.Println()
			}
			return fqdn
		} else if config.dryRun {
			fmt.Printf("%s is not a valid FQDN (%s)\n", fqdn, err)
		}
	}

	return ""
}

func buildSshArgs(user, fqdn string) ([]string, error) {
	var sshArgs []string

	if user != "" {
		sshArgs = append(sshArgs, user+"@"+fqdn)
	} else {
		sshArgs = append(sshArgs, fqdn)
	}

	sshArgs = append(sshArgs, "-p", strconv.Itoa(config.port))

	// Debug verbosity caps at 3 since that's all OpenSSH supports
	if config.debugLvl > 0 {
		sshArgs = append(sshArgs, "-"+strings.Repeat("v", min(config.debugLvl, 3)))
	}

	// TODO: Need to deal with sshpass here.
	// 	sshpass breaks down when dealing with multiple password prompts which is a common occurrence
	//  when jumphosting. Either need to disable the use of sshpass when jumphosting or implement
	//  a version of it that can handle multiple password prompts in succession
	if config.jump || config.jumphost != viper.GetString("flags.jumphost") {
		jumphost := buildFqdn(config.jumphost)
		if jumphost == "" {
			return nil, errors.New(fmt.Sprintf("Failed to find valid FQDN for jumphost %s\n", config.jumphost))
		}
		sshArgs = append(sshArgs, "-J", jumphost)

	}

	// TODO: Implement SSH options

	if config.remoteCmd != "" {
		sshArgs = append(sshArgs, config.remoteCmd)
	} else if config.tunnel {
		sshArgs = append(sshArgs, "-D", strconv.Itoa(config.tunnelPort))
	}

	return sshArgs, nil
}

func ssm(userAndHost string) error {
	var (
		err      error
		fqdn     string
		host     string
		parts    []string
		shellCmd *exec.Cmd
		sshArgs  []string
		user     string
	)

	// Handle splitting user and host if an alternate username is supplied
	// Users are expected with the standard <user>@<host> format
	parts = strings.Split(userAndHost, "@")
	if len(parts) > 2 {
		log.Fatalf("Invalid user and host combination: %s", userAndHost)
	} else if len(parts) == 2 {
		user = parts[0]
		host = parts[1]
	} else {
		user = ""
		host = parts[0]
	}

	fqdn = buildFqdn(host)
	if fqdn == "" {
		return errors.New(fmt.Sprintf("Failed to find valid FQDN for %s", host))
	}

	sshArgs, err = buildSshArgs(user, fqdn)
	if err != nil {
		return err
	}

	if config.dryRun {
		fmt.Printf("Generated SSH command: \n  %s\n", "ssh "+strings.Join(sshArgs, " "))
	} else {
		shellCmd = exec.Command("ssh", sshArgs...)

		shellCmd.Stdin = os.Stdin
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr

		err = shellCmd.Start()
		if err != nil {
			return err
		}

		err = shellCmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
