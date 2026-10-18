# i3-auto-mode

A daemon that watches [i3 window manager](https://i3wm.org/) events and
automatically switches the active [i3 mode](https://i3wm.org/docs/userguide.html#binding_modes)
based on which application is currently focused.

Useful for binding-mode setups where you want different keybindings active
depending on the focused app (e.g. workspace with a browser window, an IDE,
etc.), without manually triggering the mode switch yourself.

## Install

Grab a prebuilt package or binary from the
[releases page](../../releases) (linux/amd64 and linux/arm64):

```bash
# Debian/Ubuntu
sudo dpkg -i i3-auto-mode_*_linux_amd64.deb

# Fedora/RHEL
sudo rpm -i i3-auto-mode_*_linux_amd64.rpm

# Alpine
sudo apk add --allow-untrusted i3-auto-mode_*_linux_amd64.apk
```

Packages install the binary to `/usr/bin/i3-auto-mode` and the systemd
unit to `/usr/lib/systemd/user/i3-auto-mode.service`

Prefer a plain binary? Grab the `.tar.gz` instead — it contains the
binary plus the same systemd unit and this README.

## Configuration

Config file locations:
 * `$XDG_CONFIG_HOME/i3-auto-mode/config.yaml`
 * `~/.config/i3-auto-mode/config.yaml` if `XDG_CONFIG_HOME` is unset
 * `-config <filename>`
 * `--config <filename>`

Config file example:
```yaml
default_mode: 'default' # the mode used for empty workspace or unknown apps
modes:
  browser:
    application_classes:
      - Firefox # WM_CLASS values, use `xprop WM_CLASS` to find out a class
  whatever:
    application_classes:
      - foo
      - bar
```
In the example configuration above:
 * `browser`, `whatever` - `i3` modes
 * `Firefox`, `foo`, `bar` - WM_CLASS values for application windows

Each mode name must correspond to a `mode "<name>" { ... }` block already
defined in your i3 config — this tool only switches between modes, it
doesn't define them. An application class may only belong to one mode.

## Running it

### Via i3 config

Add to your i3 config so it starts with your session:

```
exec_always --no-startup-id $HOME/.local/bin/i3-auto-mode
```

### Via systemd (user service)

Runs the daemon under `systemd --user` instead of as a plain i3-launched
process, so it gets restarted on crash (`Restart=on-failure`).

Installed from a `.deb`/`.rpm`/`.apk` [release](../../releases) package?
The unit is already in place at `/usr/lib/systemd/user/i3-auto-mode.service`
— skip straight to the `exec` lines below. Built the binary yourself
instead? See [`systemd/i3-auto-mode.service`](systemd/i3-auto-mode.service)

```bash
mkdir -p ~/.config/systemd/user
cp systemd/i3-auto-mode.service ~/.config/systemd/user/
systemctl --user daemon-reload
```

Add to your i3 config so it starts with your session:

```
exec --no-startup-id systemctl --user import-environment DISPLAY XAUTHORITY
exec --no-startup-id systemctl --user restart i3-auto-mode.service
```

`import-environment` matters regardless of how you start i3 (`startx`, a
display manager, ...): the daemon finds i3's IPC socket by running `i3
--get-socketpath`, which needs `$DISPLAY` — and the systemd `--user`
instance doesn't necessarily have it unless something sets it explicitly.

## Build

```bash
go build -o "$HOME/.local/bin/i3-auto-mode" .
```

## Development

```bash
go build ./...   # build
go vet ./...     # static checks
go test ./...    # run tests
```

## License

[MIT](LICENSE)
