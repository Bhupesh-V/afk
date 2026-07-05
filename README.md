# afk

CLI to update slack user status across workspaces

## Why

- **Compliance friendly**.
  - AFK is completely local. All slack secrets reamin locked in your system's keychain.
- **The afk config file is dotfiles friendly**.
  - Build up as many status presets as you like!

## Installation

## Setup

TODO: Instructions on how to create a slack app

## Usage

### Going away

Set lunch status, afk will read expiration time from your config

```
afk lunch
```

afk with custom expiration time, this will be automatically saved for future `afk lunch` invocations.

```
afk lunch 67m
```

Going on a random break?

```
afk break 67m
```

### Back online

AFK has a special `clear` preset that can be used if you are back early or didn't set an expiration time earlier.

```
afk clear
```

<!-- ## Modes

1. Breakfast
2. Lunch
3. Dinner
4. Chore
5. OOO
6. Sick
7. Rest -->

## Future

- Support for Microsoft Teams (bleh)
- Support for Mattermost (💪🏽)