# xdg-user-dirs-cross

Manage your XDG directories.

Meant to replace `xdg-user-dirs` and `xdg-user-dirs-update` but adds:

- Support for the missing non-user XDG paths, namely:
   - `XDG_DATA_HOME`
   - `XDG_CONFIG_HOME`
   - `XDG_STATE_HOME`
   - `XDG_CACHE_HOME`
   - `XDG_RUNTIME_DIR`
- Cross-platform: macOS and Linux, perfect for cross-platform dotfiles.

It works very similarly, but instead of using `~/.config/user-dirs.dirs` which is both a subjectively awful name, and objectively non XDG compliant directory (feel free to correct me if I'm wrong), we use `~/.config/xdg/` both for the user editable `user.dirs` and the `generated.dirs` which is a representation of what is being applied.

## Features

- Non-destructive directory updates, when a new XDG folder is set, the previous folder is not modified in any way
- Customizable user directory locations on `~/.config/xdg/user.dirs`
- Automatic generation of `~/.config/xdg/generated.dirs`, which will be a merge of `~/.config/xdg/user.dirs` and default XDG standards as per [this XDG go library](https://github.com/adrg/xdg)

## Installation

<details>
<summary>Compiling from source</summary>

To compile `xdg-user-dirs-cross` for macOS and aarch64 Linux (Raspberry Pi), follow these steps:

1. Ensure you have Go installed on your system. You can download it from https://golang.org/dl/

2. Clone the repository:
   ```
   git clone https://github.com/yourusername/xdg-user-dirs-cross.git
   cd xdg-user-dirs-cross
   ```

3. Compile for your current system:
   ```
   go build -o xdg-user-dirs-cross
   ```

4. Move the binary to a directory in your PATH:
   ```
   sudo mv xdg-user-dirs-cross /usr/local/bin/
   ```

</details>

## Usage

The tool is designed to be evaluated by the shell. This means that the only output is the exported variables:

```
$ ./xdg_user_dirs_update_cross
export XDG_CONFIG_HOME="/home/adrian/.config"
export XDG_DATA_HOME="/home/adrian/.local/share"
export XDG_RUNTIME_DIR="/run/user/1000"
export XDG_DOCUMENTS_DIR="/home/adrian/Documents"
export XDG_MUSIC_DIR="/home/adrian/Music"
export XDG_VIDEOS_DIR="/home/adrian/Videos"
export XDG_TEMPLATES_DIR="/home/adrian/Templates"
export XDG_CACHE_HOME="/home/adrian/.local/cache"
export XDG_STATE_HOME="/home/adrian/.local/state"
export XDG_DESKTOP_DIR="/home/adrian/Desktop"
export XDG_DOWNLOAD_DIR="/home/adrian/Downloads"
export XDG_PICTURES_DIR="/home/adrian/Pictures"
export XDG_PUBLICSHARE_DIR="/home/adrian/Public"
```

So it's meant to be used like this:

```
eval "$(xdg-user-dirs-cross)"
```

In order for the env variables to be set on the user session, you can do it in any way you like. Ex: `~/.zshenv`, `~/.profile`, `~/.zshrc`, `~/.bashrc`...


1. (Optional) Edit `~/.config/xdg/user.dirs` with your desired XDG directory locations. The tool will generate a `~/.config/xdg/generated.dirs` file based on this configuration.

2. To ensure XDG environment variables are set, add the following line to your shell's startup file (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.profile`):
   ```
   eval "$(xdg-user-dirs-cross)"
   ```

   This command evaluates the output of `xdg-user-dirs-cross`, which consists of export statements. This is necessary because a Go program cannot directly modify the environment of the shell that calls it.

3. Restart your shell or log out and log back in for the changes to take effect.

The tool will generate a `~/.config/xdg/generated.dirs` file, which is a combination of user-specified directories in `user.dirs` and platform-specific defaults for directories not specified in `user.dirs`. All modifications should be done in `user.dirs`.

### Command-line Options

- `-h, --help`: Show help message
- `-d, --debug`: Enable verbose output
- `-n, --dry-run`: Simulate changes without applying them
- `-c, --create-dirs`: Create directories if they don't exist
- `-l, --log-file`: Specify the log file path (default: $HOME/.local/state/xdg-user-dirs-cross/xdg-user-dirs-cross.log)

Example usage with log file specification:
```
xdg-user-dirs-cross -l ~/xdg-update.log
```

## Configuration

- `~/.config/xdg/user.dirs`: User-defined configuration (edit this file)
- `~/.config/xdg/generated.dirs`: Generated configuration file (do not edit this file directly)

### Example Configuration

Here's an example of what your `~/.config/xdg/user.dirs` file might look like:

```
XDG_CACHE_HOME="$HOME/.local/cache"
```

For instance, I prefer to have the cache folder in `~/.local/cache` rather than `~/.cache` simply because I prefer a clutter-free `~` home :)

You can omit any directories you don't want to customize, and the tool will use platform-specific defaults.

For instance, in this case, `XDG_VIDEOS_DIR` will be `~/Videos` on Linux and `~/Movies` on macOS.

## Default Behavior

- Uses XDG-style paths for core directories on both macOS and Linux
- Applies platform-specific defaults for user directories when not specified
- Preserves exact configurations from `user.dirs`
- Generates `~/.config/xdg/generated.dirs` based on `user.dirs` and defaults

For more information on XDG Base Directory Specification, check [XDG](https://github.com/adrg/xdg).

## Command-line Options

- `-h`: Show help message
- `-d, --debug`: Enable verbose output
- `-n, --dry-run`: Simulate changes without applying them
- `-c, --create-dirs`: Create directories if they don't exist

## FAQ

Why do you remove the `~/.config/user-dirs.dirs`?

The [XDG](https://github.com/adrg/xdg) library which this program relies on:

> XDG user directories environment variables are usually not set on most operating systems. However, if they are present in the environment, they take precedence.

> On Unix-like operating systems (except macOS and Plan 9), the package reads the user-dirs.dirs config file.

[source](https://github.com/adrg/xdg?tab=readme-ov-file#xdg-user-directories)

So in order to read the actual defaults for your system, we need to first remove the file and unset the vars.

Note that the program attempts a backup of `~/.config/user-dirs.dirs` on `~/.config/xdg/user-dirs.dirs-backup` rather than just deleting it, but it may override the `~/.config/xdg/user-dirs.dirs-backup` if it already exists.

Why is `XDG_DATA_DIRS` and `XDG_CONFIG_DIRS`?

Because it's missing from the [XDG lib](https://github.com/adrg/xdg) we use.
