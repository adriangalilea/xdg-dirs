# XDG User Dirs Update Cross - Project Checklist

This checklist serves as an index of all tasks for the project, both completed and pending. It provides a comprehensive overview of the project's progress and remaining work. For day-to-day task management, we use the checklist-current.md file as a Kanban board, which focuses on immediate and upcoming tasks.

## 1. Core Functionality
- [x] Implement reading from ~/.config/user.dirs (user-defined configuration)
- [x] Implement generation of ~/.config/xdg/generated.dirs (system-wide configuration)
- [x] Implement support for core XDG directories with cross-platform consistency
- [x] Use XDG library for default values and cross-platform compatibility
- [x] Implement checks for existence of user-specific XDG directories
- [x] Implement directory creation for non-existent directories (controlled by a flag)
- [x] Ensure non-destructive updates (don't move existing content)
- [x] Set environment variables based on the combined user.dirs and default values

## 2. Configuration and Customization
- [x] Never automatically generate user.dirs file
- [x] Maintain XDG-compatible syntax (XDG_*_DIR="$HOME/Directory") in both files
- [x] Allow omission of user-specific directories in user.dirs
- [x] Ensure generated.dirs contains all directories, combining user choices from user.dirs with defaults
- [x] Add a clear warning in generated.dirs that it's generated and how to update it
- [x] Allow users to easily override any directory location in user.dirs

## 3. Cross-Platform Compatibility
- [x] Use XDG library to handle platform-specific defaults (e.g., "Movies" on macOS vs "Videos" on Linux)
- [x] Ensure consistent behavior for core XDG directories across platforms

## 4. Path Handling and Security
- [x] Implement path normalization
- [x] Handle paths relative to $HOME
- [x] Clarify and implement policy on expanding environment variables in paths
- [x] Implement best practices for handling file permissions when creating directories or files
- [x] Ensure secure handling of user-provided paths

## 5. Command-Line Interface
- [x] Implement basic run command to generate generated.dirs and set up directories
- [x] Add --debug (-d) flag for verbose output
- [x] Add --dry-run (-n) option for simulating changes without applying them
- [x] Add option to show current settings (via --debug flag)
- [x] Implement flags for controlling directory creation behavior (-c)
- [x] Add flag for specifying custom log file path (-l)
- [x] Add help flag (-h, --help) to display usage information

## 6. Error Management and Logging
- [x] Implement graceful handling of permission issues
- [x] Create clear, informative error messages
- [x] Implement proper error logging system
- [x] Handle conflicts between default paths and user-specified paths (if any)
- [x] Implement logging system for actions taken (only when --debug is used)
- [x] Ensure low verbosity by default, with option for increased verbosity
- [x] Create dedicated logger package
- [x] Implement file-based logging for non-debug mode
- [x] Implement logging system with file output and debug capabilities
- [x] Implement log rotation to prevent large log files

## 7. Code Structure and Documentation
- [x] Implement core functionality in separate package
- [x] Update method names to match actual XDGDirs struct methods
- [x] Review and update code comments to clarify design decisions
- [x] Ensure code readability without over-commenting

## 8. Testing and Quality Assurance
- [x] Implement minimal testing to ensure basic functionality
- [x] Review and potentially expand existing minimal tests
- [ ] Initial real-world testing
- [ ] Develop minimal unit tests
- [ ] Create integration tests for different platforms (Linux, macOS, Raspberry Pi)
- [ ] Implement automated testing in CI/CD pipeline

## 10. Build Process & Distribution
- [x] Create a streamlined process for building on different architectures
  - [x] `goreleaser`
  - [x] GitHub actions
- [x] Add compilation instructions for macOS and aarch64 Linux (Raspberry Pi)
- [ ] Create Homebrew formula for macOS installation
- [ ] Create package configurations for common Linux package managers

## 12. Documentation and README Updates
- [x] Update README with:
  - [x] Describe all completed tasks translated into user-facing language as features
  - [x] Provide clear documentation on customization options
  - [x] Clear explanation of user.dirs vs user-dirs.dirs
  - [x] Explanation of how user-dirs.dirs is generated and its purpose
  - [x] Usage instructions for maintaining consistent directories across systems
  - [x] Examples of user.dirs configuration
  - [x] Explanation of command-line options
  - [x] Rationale for design decisions
- [x] Update README.md with new installation and usage instructions
- [x] Instruct users to add the binary to their PATH
- [x] Provide a command to evaluate the script output in a file that executes early (for both macOS and Linux)
- [x] Review and potentially expand existing documentation
- [x] Update all documentation to reflect the final state of the project
- [x] Add information about the logging system and its behavior

## 13. Final Review and Cleanup
- [x] Perform comprehensive code review
- [x] Address any remaining TODOs or FIXMEs in the code
- [x] Ensure consistent coding style throughout the project
- [x] Conduct thorough code review to ensure quality and consistency

## 14. Attempted Service Implementation Strategies (Failed)
- [x] Implement systemd user service for Linux
- [x] Create launchd user agent for macOS
- [x] Develop custom PAM module for Linux
- [x] Utilize /etc/environment.d/ for systemd-based systems
- [x] Modify /etc/pam.d/sshd to include custom PAM configuration
- [x] Create shell script to be executed by PAM
- [x] Attempt to use pam_exec.so in PAM configuration
- [x] Try setting environment variables through systemd service
- [x] Experiment with different service types (oneshot, simple, forking)
- [x] Attempt to use LoginHook on macOS
- [x] Try modifying /etc/profile and /etc/zprofile
- [x] Experiment with creating a custom PAM session module

## 15. Rebrand

- [x] rename to xdg-dirs since it's not only covering xdg-user-dirs as the original tool, but all, and the -cross part makes it very long, either way name doesn't conflict now since we remove the `user` part
