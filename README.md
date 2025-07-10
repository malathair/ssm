# Malathair's Simple SSH Manager

> [!CAUTION]
> This tool is still in the process of being ported to Go. Functionality may be missing or broken.

> [!NOTE]
> Please be aware that this tool will not have 100% feature parity with the original Python project
> and that it is more of a spiritual successor rather than a port

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

### Prerequisites

- OpenSSH SSH client

### Installation
TBD...

### Updating
TBD...

### Configuration
TBD...

### Uninstalling
TBD...
