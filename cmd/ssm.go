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

func buildFqdn(userAndHost string) (string, string) {
	var (
		host string
		user string
	)

	// Handle splitting user and host if an alternate username is supplied
	// Users are expected with the standard <user>@<host> format
	parts := strings.Split(userAndHost, "@")
	if len(parts) > 2 {
		log.Fatalf("Invalid user and host combination: %s", userAndHost)
	} else if len(parts) == 2 {
		user = parts[0]
		host = parts[1]
	} else {
		host = parts[0]
	}

	ip := net.ParseIP(host)
	if ip != nil {
		if config.dryRun {
			fmt.Printf("%s is valid IP address\n", host)
			fmt.Println()
		}
		return user, ip.String()
	} else if config.dryRun {
		fmt.Printf("%s is not a valid IP Address\n", host)
	}

	_, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", host)
	if err == nil {
		if config.dryRun {
			fmt.Printf("%s is a valid FQDN\n", host)
			fmt.Println()
		}
		return user, host
	} else if config.dryRun {
		fmt.Printf("%s is not a valid FQDN (%s)\n", host, err)
	}

	if len(config.domains) < 1 {
		fmt.Println("No domains found for FQDN autocompletion")
		return user, ""
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
			return user, fqdn
		} else if config.dryRun {
			fmt.Printf("%s is not a valid FQDN (%s)\n", fqdn, err)
		}
	}

	return user, ""
}

func buildSshArgs(user, fqdn string) ([]string, error) {
	var sshArgs = []string{"ssh"}

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

	if config.jump || config.jumphost != viper.GetString("flags.jumphost") {
		jumpuser, jumphost := buildFqdn(config.jumphost)
		if jumphost == "" {
			return nil, errors.New(fmt.Sprintf("Failed to find valid FQDN for jumphost %s\n", config.jumphost))
		}
		if jumpuser != "" {
			sshArgs = append(sshArgs, "-J", jumpuser+"@"+jumphost)
		} else {
			sshArgs = append(sshArgs, "-J", jumphost)
		}
	} else if config.sshpass && user == "" {
		sshArgs = append([]string{"sshpass", "-e"}, sshArgs...)
	}

	// TODO: Implement SSH options

	if config.remoteCmd != "" {
		sshArgs = append(sshArgs, config.remoteCmd)
	} else if config.tunnel {
		sshArgs = append(sshArgs, "-D", strconv.Itoa(config.tunnelPort))
	}

	return sshArgs, nil
}

func ssm(host string) error {
	var (
		err      error
		fqdn     string
		shellCmd *exec.Cmd
		sshArgs  []string
		user     string
	)

	user, fqdn = buildFqdn(host)
	if fqdn == "" {
		return errors.New(fmt.Sprintf("Failed to find valid FQDN for %s", host))
	}

	sshArgs, err = buildSshArgs(user, fqdn)
	if err != nil {
		return err
	}

	if config.dryRun {
		fmt.Printf("Generated SSH command: \n  %s\n", strings.Join(sshArgs, " "))
	} else {
		shellCmd = exec.Command(sshArgs[0], sshArgs[1:]...)

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
