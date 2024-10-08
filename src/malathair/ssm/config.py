import getpass
import os
import shutil
import socket
import toml


USER_CONFIG_PATH = os.path.expanduser("~/.config/ssm.conf")
GLOBAL_CONFIG_PATH = "/usr/local/etc/ssm.conf"


class Config:
    def __init__(self):
        self.ssh_port: str = "22"
        self.jump_host: str | None = None
        self.tunnel_port: str = "6060"
        self.domains: list[str] | None = None

        self.sshpass: bool = self.__sshpass_available()

        self.__load_config_file()

    def __sshpass_available(self) -> bool:
        if shutil.which("sshpass") is None:
            return False

        if os.getenv("SSHPASS") is None:
            return False

        return True

    def __locate_config_file(self) -> str | None:
        if os.path.isfile(USER_CONFIG_PATH):
            return USER_CONFIG_PATH

        if os.path.isfile(GLOBAL_CONFIG_PATH):
            return GLOBAL_CONFIG_PATH

        return None

    def __load_config_file(self) -> None:
        path = self.__locate_config_file()

        if path is None:
            return

        with open(path, "r", encoding="utf-8") as file:
            config = toml.load(file)

        try:
            self.ssh_port = str(config["ssh"]["port"])
        except KeyError:
            pass

        try:
            self.jump_host = config["ssh"]["jump"]
        except KeyError:
            pass

        try:
            self.tunnel_port = str(config["tunnel"]["port"])
        except KeyError:
            pass

        try:
            self.domains = config["domains"]
        except KeyError:
            pass

    def get_config_dict(self) -> dict:
        config = {}

        config["domains"] = [] if self.domains is None else self.domains
        config["ssh"] = {}
        config["ssh"]["port"] = self.ssh_port
        config["ssh"]["jump"] = "" if self.jump_host is None else self.jump_host
        config["ssh"]["sshpass"] = self.sshpass
        config["tunnel"] = {}
        config["tunnel"]["port"] = self.tunnel_port

        return config


def ask_yes_no_question(question: str) -> bool:
    valid_response = False
    while not valid_response:
        response = input(question + " [y/N]: ")

        match response.lower():
            case "y":
                return True
            case "n":
                return False
            case "":
                return False
            case _:
                print('Invalid response: Please enter "y" or "n"')


def validator(validator_func):
    def wrapper(current_setting):
        print(current_setting)
        response = input("Enter a new value or leave blank to keep the current setting: ")

        return validator_func(response)

    return wrapper


@validator
def port_validator(port: str) -> str | None:
    try:
        if port != "":
            port_int = int(port)
            assert port_int > 0 and port_int < 65536

        return port
    except Exception:
        return None


@validator
def jumphost_validator(jump_host: str) -> str | None:
    try:
        if jump_host != "":
            assert "." in jump_host
            assert jump_host[-1] != "." and jump_host[0] != "."
            socket.gethostbyname(jump_host)
        return jump_host
    except Exception:
        question = f"The entered value ({jump_host}) does not appear to be a valid FQDN. Would you like to use it anyways?"

        if ask_yes_no_question(question):
            print("Forcing use of entered value...")

            return jump_host

        return None


def domain_validator(domain: str) -> str | None:
    try:
        assert "." in domain
        assert domain[-1] != "." and domain[0] != "."
        return domain
    except Exception:
        question = (
            f"The entered value ({domain}) appears to be invalid. Would you like to use it anyways?"
        )

        if ask_yes_no_question(question):
            print("Forcing use of entered value...")
            return domain

        return None


def domain_editor_input_handler(domains: list[str]) -> tuple[list[str], bool]:
    while True:
        selection = input("Enter a menu option: ")

        match selection:
            case "1":
                new_domain = domain_validator(input("Enter a new domain: "))
                if new_domain is not None:
                    domains.append(new_domain)
                return domains, True
            case "2":
                domain_to_remove = input("Enter a domain to remove: ")
                if domain_to_remove in domains:
                    domains.remove(domain_to_remove)
                return domains, True
            case "3":
                return None, False
            case "4":
                return domains, False
            case _:
                print("Invalid selection. Please enter the number of a menu item to proceed")


def edit_domains(domains: list[str]) -> list[str]:
    old_domains = domains.copy()
    edit = True
    while edit:
        print_domain_editor_menu(domains)
        new_domains, edit = domain_editor_input_handler(domains.copy())

        if new_domains is None:
            return old_domains

        domains = new_domains

    return domains


def print_domain_editor_menu(domains: list[str]) -> None:
    os.system("clear")
    print("Domain Editor")
    print("===============================================\n")
    print("Current domain list: \n")
    for i in range(len(domains)):
        print(f"{domains[i]}")
    print("\n===============================================\n")
    print("\t1. Add")
    print("\t2. Remove")
    print("\t3. Exit domain editor without saving")
    print("\t4. Exit domain editor and save\n")


def print_current_configuration(config: dict) -> None:
    print("===============================================\n")
    print("The current configuration is:\n")
    print(toml.dumps(config))
    print("===============================================\n")


def configure() -> None:
    path = GLOBAL_CONFIG_PATH if getpass.getuser() == "root" else USER_CONFIG_PATH
    config = Config()
    config_toml = config.get_config_dict()

    os.system("clear")

    print("Welcome to the configuration utility for ssm!")
    print_current_configuration(config_toml)

    valid_response = False
    while not valid_response:
        if ask_yes_no_question("Would you like to edit the current configuration?"):
            valid_response = True
        else:
            print("\nExiting the configuration utility\n")
            return

    os.system("clear")

    # Set ssh port value
    valid_response = False
    while not valid_response:
        port = port_validator(f'The current default ssh port is {config_toml["ssh"]["port"]}.')

        if port is None:
            print("Invalid port detected. Please enter a valid port number")
        else:
            if port != "":
                config_toml["ssh"]["port"] = port
            valid_response = True

    os.system("clear")

    # Set jumphost value
    valid_response = False
    while not valid_response:
        jump_host = jumphost_validator(
            f'The current default jumphost is {config_toml["ssh"]["jump"]}.'
        )

        if jump_host is not None:
            if jump_host != "":
                config_toml["ssh"]["jump"] = jump_host
            valid_response = True

    os.system("clear")

    # Enable/Disable sshpass usage
    response = ask_yes_no_question(
        "Would you like to enable the use of sshpass to automate password logins?"
    )

    if response and shutil.which("sshpass") is None:
        print(
            "\nWarning:\n"
            + "The sshpass utility does not appear to be installed, or is not located on the\n"
            + "users PATH. sshpass usage will be enabled but won't have any effect until the\n"
            + "utility has been installed on your system or made available to the current user.\n"
        )
        input("Press any key to continue")

    config_toml["ssh"]["sshpass"] = response

    os.system("clear")

    # Set Tunnel Port
    valid_response = False
    while not valid_response:
        port = port_validator(
            f'The current default proxy tunnel port is {config_toml["tunnel"]["port"]}.'
        )

        if port is None:
            print("Invalid port detected. Please enter a valid port number")
        else:
            if port != "":
                config_toml["tunnel"]["port"] = port
            valid_response = True

    os.system("clear")

    # Set Domains
    config_toml["domains"] = edit_domains(config_toml["domains"])

    os.system("clear")

    print_current_configuration(config_toml)
    if ask_yes_no_question("Would you like to save the current configuration?"):
        print("Saving...")

        with open(path, "w") as file:
            file.write(toml.dumps(config_toml))

        print(f"Configuration saved at {path}")
    else:
        print("Exiting without saving!")
