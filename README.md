# afk

For remote only teams, profile status on collaboration tools like slack is an important indicator of your availability, which means those statuses require a well-thought out approach, `afk` takes this small thing a bit too seriously!

> [!IMPORTANT]
> `afk` is experimental and has only been tested on a modern MacOS system. Please report issues.

## Why?

- **Support for multiple (slack) workspaces**
  - No brainer for consultants! 👀
- **No secrets management & Compliance friendly**.
  - `afk` is completely local. All secrets reamin locked in your system's keychain.
- **The afk config file is dotfiles friendly**.
  - Build up as many status presets as you like!
  - The config file is privacy friendly, afk never exposes which workspace(s) you are a part of.
- **Unified availability status**
  - `afk` forces the user to club things like, notifications, presence and status elements under [one status preset](https://github.com/Bhupesh-V/afk/blob/main/config.sample.toml#L35-L40) config.

## Why not?

- Initial friction while setting up the app (see [setup](#setup)) in workspaces you want to operate on.
- Possible disapproval by slack workspace admins.
  - Admins can kick any app, I recommend not telling them in advance, if they find it, explain why you need it (check readme header).
- Workspace App Limit, For FREE plan, slack has a limit of 10 apps, if you are in luck, you can setup afk. For Pro/Enterprise plans, there's no limit. However, I recommend not talking about `afk` to your colleagues (imagine 100 afk-bots on a slack workspace 😬)

## Setup

1. Setup your preferred client.
2. Run `afk -s/--setup` and follow on-screen instructions.

### Slack

1. Go to the [Slack API Dashboard](https://api.slack.com/apps) and click Create New App (choose From Scratch).
2. Assign it an initial development workspace (choose whichever).
3. Copy the **Client ID** and **Client Secret** (we'll need them later on).
4. Navigate to **OAuth & Permissions** in the left sidebar:
   1. Scroll down to **Scopes**.
   2. Under **User Token Scopes**, add the specific user scopes that `afk` needs:
      - `users.profile:write`
        - Reason: `Gives the users the ability to update user profile status message and emoji across workspaces`
      - `users:write`
        - Reason: `Gives users the ability to update their presence (active, away) across workspaces`
      - `dnd:write`
        - Reason: `Gives users the ability to update their notification settings across workspaces`
5. Scroll up to Redirect URLs and add the following OAuth callback endpoint
   ```
   https://localhost:8080/oauth/callback
   ```
6. In your Slack App Dashboard, click **Manage Distribution** in the left sidebar.
7. Complete the Distribute App Checklist.
8. Under the section **Share Your App with Other Workspaces**, click **Activate Public Distribution**. Do not submit your app to the Slack Marketplace.

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