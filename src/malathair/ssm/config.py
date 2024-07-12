import os
import shutil
import toml

USER_CONFIG_PATH = os.path.expanduser("~/.config/ssm.conf")
GLOBAL_CONFIG_PATH = "/usr/local/etc/ssm.conf"


class Config:
    def __init__(self):
        self.ssh_port: str = "22"
        self.jump_host: str | None = None
        self.tunnel_port: str = "6060"
        self.domains: list[str] | None = None

        self.sshpass: bool = self._sshpass_available()

        self.load_config_file()

    def _sshpass_available(self) -> bool:
        if shutil.which("sshpass") is None:
            return False

        if os.getenv("SSHPASS") is None:
            return False

        return True

    def _locate_config_file(self) -> str | None:
        if os.path.isfile(USER_CONFIG_PATH):
            return USER_CONFIG_PATH

        if os.path.isfile(GLOBAL_CONFIG_PATH):
            return GLOBAL_CONFIG_PATH

        return None

    def load_config_file(self) -> None:
        path = self._locate_config_file()

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
