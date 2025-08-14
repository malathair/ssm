#!/bin/sh

get_install_dir() {

    INSTALL_DIR="/usr/local/bin"

    # Check to see if this was run as root. If it was, then install to /usr/local/bin
    # otherwise install to the user's private bin
    if [ "$(id -u)" -ne 0 ]; then

        INSTALL_DIR="$HOME/bin"

        # Check to see if the user has an existing private bin
        # Prefer ~/bin over ~/.local/bin if both exist
        # If there is not an existing private bin then create it at ~/.local/bin
        if [ ! -d "$HOME/bin" ]; then

            INSTALL_DIR="$HOME/.local/bin"

            if [ ! -d "$HOME/.local/bin" ]; then
                mkdir -p "$HOME/.local/bin"
            fi
        fi

    fi

    echo "$INSTALL_DIR"

}

FILE="$(get_install_dir)/ssm"

curl -fsSL "https://github.com/malathair/ssm/releases/latest/download/linux_amd64_ssm" -o $FILE
chmod 755 $FILE

