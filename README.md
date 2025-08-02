# Malathair's Simple SSH Manager

<!-- PROJECT SHIELDS -->
[![Version][version-shield]][version-url]
[![Go Reference][reference-shield]][reference-url]
[![Go Report Card][reportcard-shield]][reportcard-url]
[![MIT License][license-shield]][license-url]

<br/>

## About the Project

SSM is a CLI utility that attempts to provide a better SSH experience by providing a wrapper
for OpenSSH's ssh command. It does so through a simplified interface with some sensible
defaults that automatically generates and executes OpenSSH's ssh command in the background.
The defaults chosen by SSM can be overridden for all hosts using a TOML based configuration
file, or on a per-host basis through flags.

This tool also provides some additional functionality that is difficult or impossible to
achieve through a typical SSH config file. This includes functionality such as dynamic FQDN
completion and password autofill through the use of sshpass.

This is a port/re-imagining of my original Python implementation which can be found here:
[Malathair's Python Simple SSH Manager](https://github.com/malathair/ssm-python)

<br/>

## Getting Started

### Prerequisites
Additional software needed to use SSM

- OpenSSH SSH client

### Installation
SSM can be installed by downloading the [latest precompiled binary][version-url]

### Configuration
SSM expects the configuration file called `ssm.conf` to exist in the current working
directory or in the user's config home dir.

An example configuration file can be found here: [examples/ssm.conf](examples/ssm.conf)

<br/>

## Updating
TBD...

<br/>

## Uninstalling
SSM can be uninstalled by running:

```bash
rm $(which ssm)
```


<!-- MARKDOWN LINKS & IMAGES -->
[license-shield]: https://img.shields.io/github/license/malathair/ssm.svg
[license-url]: https://github.com/malathair/ssm/blob/main/LICENSE
[reference-shield]: https://pkg.go.dev/badge/github.com/malathair/ssm.svg
[reference-url]: https://pkg.go.dev/github.com/malathair/ssm
[reportcard-shield]: https://goreportcard.com/badge/github.com/malathair/ssm
[reportcard-url]: https://goreportcard.com/report/github.com/malathair/ssm
[version-shield]: https://img.shields.io/github/release/malathair/ssm.svg
[version-url]: https://github.com/malathair/ssm/releases/latest