# mcphee11-tui

A TUI to consolidate many of the different CLI capabilities I have built up in individual repos.

## Features

- Beautiful TUI built with [Charm Bubbletea](https://github.com/charmbracelet/bubbletea)
- Search the [Genesys Cloud Release Notes](https://help.mypurecloud.com/release-notes-home/monthly-archive/)
- Build Banking PWA
- Update TTS
- Common Modules
- Google Bot Migration
- Backup Flows

## Installation

### Prerequisites

- Go 1.24.2 or higher

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
MCPHEE11_TUI_REGION=YOUR_REGION    eg:mypurecloud.com.au
MCPHEE11_TUI_CLIENT_ID=YOUR_ID
MCPHEE11_TUI_SECRET=YOUR_Secret
```

The region needs to be the URL like for example `mypurecloud.com.au` if an ORG is set it will display the ORG Name in the home page Title.

## Navigation

When at the main menu press `?` to see the additional help menu

## Debugging

While I do try to pass status update messages to the TUI there are at times longer messages and more debug messages that can help you. To enable this set the environment variable

```
MCPHEE_TUI_DEBUG=true
```

This will then create a file `mcphee11-tui-log.log` in the location that you are running the TUI. If you have issues please enable the log file have a look in there and if your still not sure then raise an issue in the repo with log details. Check the log for any flow names etc that you need to mask out first before submitting.

As this is an example TUI and not an official product of any kind I will do my best to assist, I'm also open to PR's if people find issues and want to contribute to this code as well.
