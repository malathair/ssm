# Simple SSH Manager

SSM is a small Python wrapper to simplify some common SSH use cases. It is essentially just a wrapper around OpenSSH's ssh command that builds and executes the ssh command for you from a simplified set of inputs.

## Gettings Started

### Prerequisites

 - Python3.11 or newer
 - Pipx

### Installation

Install the application with pipx using the following command:

```bash
pipx install https://github.com/malathair/ssm/releases/download/v1.0.0/malathair_ssm-1.0.0-py3-none-any.whl
```

And then install the config file:

```bash
curl https://raw.githubusercontent.com/malathair/ssm/v1.0.0/example-conf/ssm.conf | sudo tee /usr/local/etc/ssm.conf 2&>/dev/null
```

### Updating

### Configuration

## Usage

### Basic

SSM’s main purpose is to eliminate the need to type out full URLs when accessing into remote systems. To that end, SSM implements a system to automatically detect and build valid URLs for you from minimal input.

SSM will detect if an IP was entered as the destination host rather than a URL. If so, it will attempt to connect to that IP directly with no further adjustments

If an IP is not detected, SSM will check to see if a valid URL was entered as the destination host

If the destination host is not an IP and does not resolve, then SSM assumes that it is a partial URL and will try to find a domain that resolves using the domains listed in the in the config file. It does this by combining the destination host with each domain in turn and then doing a DNS lookup to see if it resolves

As soon as a valid destination host is found, it then initiates an SSH connection to the destination. Domains are evaluated in the order they are listed in the config file so if the partial url provided would resolve on multiple domains, only the first will be used.

<details>
    <summary>Examples</summary>

    ssm 192.168.0.1: Connects directly to an IP (On the local network in this case, but public destinations work too)

    ssm user@192.168.0.1: Same as the pervious example but with an alternate username

    ssm fake.example.net: Connects directly to fake.example.net (bogus example domain) assuming it resolves

    ssm fake: Connects to fake.example.net again if the "example.net" domain is listed in the config file

    ssm user@fake: Same as the previous example but with an alternate user name

</details>

### Extended

SSM also provides a few additional argument flags for dealing with use cases that are a bit more niche but still happen on occasion. Each argument has a short and long version of the flag and you can use either. The flags SSM implements are as follows:

-h or --help: Prints the usage information of the program on the CLI for a quick reference if you need it. If used, SSM will only print the usage information of the program and exit ignoring all other provided arguments.

<details>
    <summary>Help Example</summary>

    ssm [-h] [-j | -J JUMPHOST] [-p PORT] [-t] host

    An SSH wrapper to simplify life. The config file can be found at /usr/local/etc/ssm.conf

    positional arguments:
        host              Subdomain of the host's url or the host's IP address

    options:
        -h, --help        show this help message and exit

        -j, --jump        SSHs via the jump host specified in the configuration file
        -J, --jumphost    Overrides the jump host specified in the configuration file
        -p, --port        Specifies the port to use for the SSH session

        -t, --tunnel      Start a SOCKS5 tunnel on the port defined in the configuration file

</details>

-j or --jump: Routes the SSH connection to the destination host through an intermediate secondary host defined in the configuration file. This is useful for the odd situation where you need to access a site router that can’t be reached by the VPN for some reason. By default, the intermediate host is the Logan office router.

-J or --jumphost: This option allows you to override the intermediary host defined in the configuration file for a specific connection. This is useful for SSH’ing into equipment behind a site router without needing to interact with the CLI of said site router.

<details>
    <summary>Jumphost Examples</summary>

    Example command(s) using the default jumphost:

    ssm -j fake: Connects to fake.example.net by routing the connection through the jumphost defined in the configuration file

    ssm -j user@fake: Same as the previous command but with an alternate user

    Example command(s) using an alternate jumphost:

    ssm -J jump.example.net fake: Connects to fake.example.net by routing the connection through jump.example.net instead of the default jumphost

</details>

-p or --port: This option allows you to override the default SSH port defined in the configuration file for a specific connection. This is useful for infrequent situations where the port differs from the norm.

-t or --tunnel: This option creates a SOCKS tunnel on the SSH connection which allows you to forward TCP connections & traffic to a remote network as though you were on a local machine. This is useful for accessing services (like a web UI) on a device that would not normally not be accessible. It can also be leveraged to run scripts on remote networks where it would normally be impossible to do so. Running scripts over the SOCKS tunnel requires additional tools (like proxychains or tsocks) to redirect the traffic over the SOCKS tunnel as SSM does not handle that for you.

<details>
    <summary>Tunnel Examples</summary>

    Example command(s):

    ssm -t fake: Connects to fake.example.net and binds a SOCKS5 proxy server to port 6060 on localhost. By configuring your browser to use this proxy server you can access the Web UI of resources on the network behind the destination host. Additionally you can run commands as though you were using a machine on the remote network via a 3rd party command line tool like tsocks or proxychains which will redirect tcp traffic over the proxy tunnel.

</details>
