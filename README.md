# afk

CLI to update slack user status across workspaces

> For remote friendly devs teams, profile status on collaboration software like slack is an important indicator of your availability, which means those statuses require a well-thought out approach, `afk` is a small step towards taking this small thing a bit seriously!

## Why?

- **Compliance friendly**.
  - `afk` is completely local. All secrets reamin locked in your system's keychain.
- **The afk config file is dotfiles friendly**.
  - Build up as many status presets as you like!
- **Support for multiple (slack) workspaces**
  - No brainer for consultants! 👀

## Why not?

- Initial friction while setting up the app (see [below](#setup)) in workspaces you want to operate on.
- Possible disapproval by slack workspace admins.
  - For FREE plan, slack has a limit of 10 apps, if you are in luck, you can setup afk. For Pro/Enterprise plans, slack has no limit still, I recommend not talking about `afk` to your colleagues (image a 100 afk-bots on a slack workspace 😬)
  - Slack Admins can kick any app they want, I recommend not telling them in advance, if they find it, they will reach out, than explain why (check readme header).
<!-- - You have a better workflow setup. -->
<!-- - You are not lazy to update your status. -->

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