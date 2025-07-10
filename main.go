// Simple SSH Manager (SSM) is a program to simplify the act of SSH'ing into remote hosts
//
// SSM does this by providing a wrapper interface to the underlying OpenSSH command. This allows
// shorthand specification for options that are frequently used but can otherwise be quite verbose
// when using OpenSSH directly. Where possible, an OpenSSH config file should be used, as this
// tool aims to provide a simplified interface for options that need to be dynamic, or are not
// supported by the OpenSSH config file.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/malathair/ssm/cmd"
)

func main() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		s := <-signals
		fmt.Println(s)
		os.Exit(1)
	}()

	cmd.Execute()
}
