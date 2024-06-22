# discord-updater-linux
Tired of manually editing Discord's `build_info.json` to update Discord ? This simple program will do it automatically for you !

## How to build
Install the `go` compiler:
* apt: `sudo apt install golang-go`
* pacman: `sudo pacman -S go`

Build the program:
* Open a terminal in this directory.
* Run `go build -o discord-updater-linux main.go`.

## How to install
* Copy the executable somewhere, for instance in `/usr/bin`.
* Edit the `ExecStart` line of `discord-updater-linux.service` to specify the path of the executable and the path of the build_info.json to operate on.
* Copy `discord-updater-linux.service` in `/etc/systemd/user`.
* Enable the service with `systemctl --user enable discord-updater-linux.service` (no `sudo` !).
* Start the service with `systemctl --user start discord-updater-linux.service` (no `sudo` !).