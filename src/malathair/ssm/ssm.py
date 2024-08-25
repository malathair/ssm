#!/usr/bin/env python3
# Simple SSH Manager is a program to simplify sshing to various networking devices
import argparse
import ipaddress
import socket
import subprocess

from .config import Config


# Class to override the default argparse help formatting (makes things a bit cleaner)
class HelpFormatter(argparse.HelpFormatter):
    def _format_action_invocation(self, action):
        if not action.option_strings:
            default = self._get_default_metavar_for_positional(action)
            (metavar,) = self._metavar_formatter(action, default)(1)
            return metavar

        parts = []
        parts.extend(action.option_strings)
        return ", ".join(parts)


# Define CLI arguments for the program
def arg_parser(config):
    parser = argparse.ArgumentParser(
        description="An SSH wrapper to simplify life",
        formatter_class=HelpFormatter,
    )

    session = parser.add_argument_group()
    jump = session.add_mutually_exclusive_group()
    tunnels = parser.add_argument_group()

    parser.add_argument(
        "host", type=str, help="Subdomain of the host's url or the host's IP address"
    )

    jump.add_argument(
        "-j",
        "--jump",
        action="store_true",
        help="SSHs via the jump host specified in the configuration file",
    )
    jump.add_argument(
        "-J",
        "--jumphost",
        default=config.jump_host,
        type=str,
        help="Overrides the jump host specified in the configuration file",
    )

    session.add_argument(
        "-p",
        "--port",
        default=config.ssh_port,
        type=str,
        help="Specifies the port to use for the SSH session",
    )

    tunnels.add_argument(
        "-t",
        "--tunnel",
        action="store_true",
        help="Start a SOCKS5 tunnel on the port defined in the configuration file",
    )

    return parser.parse_args()


# Return a valid IP address or hostname to connect to
def build_domain(host_arg, config):
    host_index = host_arg.find("@") + 1
    host = host_arg[host_index::]

    try:
        ipaddress.IPv4Network(host)
        return host_arg
    except ipaddress.AddressValueError:
        pass

    # If we are confident the host is not an FQDN skip this test. DNS lookups are incredibly slow
    # if the lookup value isn't an FQDN. This optimization only works with public DNS servers. If
    # using a private DNS server that contains records that are just the host portion of the FQDN
    # then this will make it impossible to resolve those
    if "." in host:
        try:
            socket.gethostbyname(host)
            return host_arg
        except socket.gaierror:
            pass

    for domain in config.domains:
        try:
            socket.gethostbyname(host + "." + domain)
            return host_arg + "." + domain
        except socket.gaierror:
            pass

    raise Exception(f'Host "{host}" is unreachable')


# Initiate the SSH session using the defined parameters
def ssh(args, config, domain):
    alt_user = False if domain.find("@") == -1 else True
    command = "ssh -p " + args.port + " " + domain

    # Jumphosting causes problems with sshpass. So only use sshpass
    # if we are not jumphosting
    if args.jump or (args.jumphost != config.jump_host):
        command = command + " -J " + args.jumphost
    elif config.sshpass and not alt_user:
        command = "sshpass -e " + command

    if args.tunnel:
        command = command + " -D " + config.tunnel_port

    return subprocess.run(command.split(), check=True)


def main():
    config = Config()
    args = arg_parser(config)

    try:
        domain = build_domain(args.host, config)
        ssh(args, config, domain)
    except Exception as e:
        print(e)


if __name__ == "__main__":
    main()
