# Kanban Board

Rules:
1. Three columns: PENDING, DOING, BLOCKED, DONE
2. Tasks have urgency labels [NEXT], [HIGH], [MEDIUM], [LOW]
3. Only one task (and its subtasks) in DOING
4. Periodically clear DONE tasks.
  4.1. Determine what warrants a mention in the `checklist.md` file, and add it there as completed.
  4.2. Clean the DONE section.

## PENDING (ordered by urgency)

- [MEDIUM] Test in macOS
- [MEDIUM] Set up GitHub Actions for automatic releases
- [LOW] Create Homebrew formula for macOS installation
- [LOW] Create package configurations for common Linux package managers
- [LOW] Develop minimal unit tests
  - Update existing tests as it's likely outdated
  - Write tests for updater package
  - Write tests for xdgdirs package
- [LOW] Create integration tests for different platforms (Linux, macOS, Raspberry Pi)
- [LOW] Implement automated testing in CI/CD pipeline

## DOING

- [HIGH] Create a streamlined process for building on different architectures

## BLOCKED (until)
- (24/08/01) Test in a Raspberry Pi aarch64

## DONE

