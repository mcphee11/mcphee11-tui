# mcphee11-tui

A TUI to consolidate many of the different CLI capabilities I have built up in individual repos.

## Features

- Beautiful TUI built with [Charm Bubbletea](https://github.com/charmbracelet/bubbletea)
- Search the [Genesys Cloud Release Notes](https://help.mypurecloud.com/release-notes-home/monthly-archive/)
- Build Banking PWA
- Update to Genesys Enhanced TTS
- Google Bot Migration
- Backup Flows

## Installation

### Prerequisites

- Go 1.18 or higher

### Install

```bash
go install github.com/mcphee11/mcphee11-tui@latest
```

Then in a new terminal run:

```bash
mcphee11-tui
```

## Genesys Cloud ORG OAuth

To leverage the parts of the TUI that leverage Genesys Cloud you will need to supply `region clientId secret` from your ORG that is a `client credentials` these are required to be set a environment variables

```
MCPHEE11_TUI_REGION
MCPHEE11_TUI_CLIENT_ID
MCPHEE11_TUI_SECRET
```

The region needs to be the URL like for example `mypurecloud.com.au` if an ORG is set it will display the ORG Name in the home page Title.

## Navigation

When at the main menu press `?` to see the additional help menu
