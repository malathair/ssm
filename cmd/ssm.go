package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func buildFqdn(host string) (string, error) {
	// TODO: build domain
	// For now just return the provided host as is

	// TODO: test if the provided host is a valid IPv4 address
	// TODO: If not a valid IPv4 address, check if the provided host is a valid FQDN
	// TODO: If not a valid FQDN, attempt to build a valid FQDN from the config's domain list

	//	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", args[0])
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	return host, nil
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
		jumphost, err := buildFqdn(config.jumphost)
		if err != nil {
			return nil, err
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

	fqdn, err = buildFqdn(host)
	if err != nil {
		return err
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
