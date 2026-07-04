# afk

CLI to update slack user status across workspaces

## Installation

## Setup

Instructions on how to create a slack app

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

If you didn't set an expiration time while setting status

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