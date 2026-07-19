# afk

For remote only teams, profile status on collaboration tools like slack is an important indicator of your availability, which means those statuses require a well-thought out approach, `afk` takes this small thing a bit too seriously!

> [!IMPORTANT]
> `afk` is experimental and has only been tested on a modern MacOS system. Please report issues.

## Why?

- **Support for multiple (slack) workspaces**
  - No brainer for consultants! 👀
- **No secrets management & Compliance friendly**
  - `afk` is completely local. All secrets reamin locked in your system's keychain.
- **The afk config file is [dotfiles friendly](https://github.com/Bhupesh-V/.Varshney/blob/master/.config/afk/config.toml)**
  - Build up as many status presets as you like!
  - The config is also privacy friendly, afk never exposes which workspace(s) you are a part of.
- **Unified availability status**
  - `afk` config club things like `status`, `presence` and `dnd` elements under [one status preset](https://github.com/Bhupesh-V/afk/blob/main/internal/config/sample.toml#L31-L42) config.

## Why not?

- Initial friction while setting up the app (see [setup](#setup)) in workspaces you want to operate on.
- Possible disapproval by slack workspace admins.
  - Slack Admins can kick any app, I recommend not telling them in advance, if they find it, explain why you need it (check readme header).
- Slack Workspace App Limit
  - For FREE plan, slack has a limit of 10 apps, if you are in luck, you can setup afk. Although you can bypass this by creating your own free workspace.
  - For Pro/Enterprise plans, there's no limit. However, I recommend not talking about `afk` to your colleagues (imagine 100 afk-bots on a slack workspace 😬)

## Install

Grab the binary from [releases](https://github.com/Bhupesh-V/afk/releases).

## Setup

### Slack

1. Go to [Slack API Dashboard](https://api.slack.com/apps) on a browser where you are logged in to any slack workspace.
2. Under Your **App Configuration Tokens**, click *Generate Token*. Assign the only available initial development workspace (don't worry about this).
3. Copy the **Access Token** (we'll need it later on).
4. Run `afk --setup`.
5. Provide the access token you just copied.
6. Follow on-screen instructions.

#### Enable Multi-Workspace Setup

1. In Slack App Dashboard, choose the automatically created "*AFK*" app and click **Manage Distribution** in the left sidebar.
2. Complete the Distribute App Checklist.
3. Under the section **Share Your App with Other Workspaces**, click **Activate Public Distribution**. Do not submit your app to the Slack Marketplace.
4. Re-run `afk -s`.

## Usage

### Going away

> Note: you can build your own status presets, `lunch` below is just an example.

afk with custom expiration time, this will be **automatically saved for future `afk lunch` invocations**.

```
afk lunch 67m
```

Set lunch status, afk will read expiration default duration from your config

```
afk lunch
```

### Back online

afk has a `-c\--clear` flag that can be used if you are back early or didn't set an expiration time earlier.

```
afk -c
```

## Config

afk ships with some default presets, modify them to your liking by editing the config file:

- Linux/MacOS: `.config/afk/config.toml`
- Windows: `C:\Users\<UserName>\AppData\Roaming\afk\config.toml`