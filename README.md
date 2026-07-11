# afk

> For remote only teams, profile status on collaboration tools like slack is an important indicator of your availability, which means those statuses require a well-thought out approach, `afk` is a small step towards taking this small thing a bit too seriously!

## Why?

- **Support for multiple (slack) workspaces**
  - No brainer for consultants! 👀
- **Compliance friendly**.
  - `afk` is completely local. All secrets reamin locked in your system's keychain.
- **The afk config file is dotfiles friendly**.
  - Build up as many status presets as you like!
  - The config file is privacy friendly, afk never exposes which workspace(s) you are a part of.

## Why not?

- Initial friction while setting up the app (see [setup](#setup)) in workspaces you want to operate on.
- Possible disapproval by slack workspace admins.
  - Admins can kick any app, I recommend not telling them in advance, if they find it, explain why you need it (check readme header).
- Workspace App Limit, For FREE plan, slack has a limit of 10 apps, if you are in luck, you can setup afk. For Pro/Enterprise plans, there's no limit. However, I recommend not talking about `afk` to your colleagues (imagine 100 afk-bots on a slack workspace 😬)


## Installation

## Setup

TODO: Instructions on how to create a slack app

## Usage

### Going away

> Note: you can build your own status presets, `lunch` below is just an example.

Set lunch status, afk will read expiration default duration from your config

```
afk lunch
```

afk with custom expiration time, this will be **automatically saved for future `afk lunch` invocations**.

```
afk lunch 67m
```

Going on a random break?

```
afk break 67m
```

### Back online

AFK has a `-c\--clear` flag that can be used if you are back early or didn't set an expiration time earlier.

```
afk -c
```

## Future

- Support for Microsoft Teams (bleh)
- Support for Mattermost