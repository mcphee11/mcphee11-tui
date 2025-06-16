# mcphee11-tui

A TUI to consolidate many of the different CLI capabilities I have built up in individual repos.

## Features

- Beautiful TUI built with [Charm Bubbletea](https://github.com/charmbracelet/bubbletea)
- Search the [Genesys Cloud Release Notes](https://help.mypurecloud.com/release-notes-home/monthly-archive/)
- Search the extended release notes for [Embedded Clients](https://help.mypurecloud.com/articles/release-notes-for-the-genesys-cloud-embedded-clients/), [SCIM](https://help.mypurecloud.com/articles/release-notes-for-genesys-cloud-scim-identity-management/), [DataActions](https://help.mypurecloud.com/articles/release-notes-for-the-data-actions-integrations/), Desktop Apps: [MAC](https://help.mypurecloud.com/release-notes-home/genesys-cloud-for-mac-desktop-app-release-notes/), [WIN](https://help.mypurecloud.com/release-notes-home/genesys-cloud-for-windows-desktop-app-release-notes/), [GCBA](https://help.mypurecloud.com/articles/genesys-cloud-background-assistant-gcba-release-notes/)
- Build Banking PWA (Demo)
- Update TTS across all flows
- Common Modules update dependencies
- Google Bot Migration -> Digital BOT Flow
- Backup Flows locally to YAML
- Export lists to CSV

## Installation

### Prerequisites

- Go 1.24.2 or higher

### Install

```
go install github.com/mcphee11/mcphee11-tui@latest
```

Then in a new terminal run:

```
mcphee11-tui
```

- NOTE: depending on your GO install you may need to add this to your PATH. If GO is configured correctly you will not need to do this manually as the install will add it to the PATH.

## Genesys Cloud ORG OAuth

To leverage the parts of the TUI that leverage Genesys Cloud you will need to supply `region clientId secret` from your ORG that is a `client credentials` these are required to be set a environment variables

```
MCPHEE11_TUI_REGION=YOUR_REGION    eg:mypurecloud.com.au
MCPHEE11_TUI_CLIENT_ID=YOUR_ID
MCPHEE11_TUI_SECRET=YOUR_Secret
```

If your like me then you have the [Genesys Cloud CLI](https://developer.genesys.cloud/devapps/cli/) running, because of this I have also allowed you to use this CLI profiles configuration for the TUI login. If you pass nothing it will attempt to use the `default` profile in the config.toml file you can also specify a profile with the environment variable

```
MCPHEE11_TUI_PROFILE=YOUR_PROFILE
```

If you wish to use a specific profile from the gc cli configuration you already have setup. If there is neither no gc cli setup or specific `MCPHEE11_TUI....` configuration then you will still be able to do things like search the release notes etc but not leverage the components of the TUI that require access to a Genesys Cloud ORG.

The region needs to be the URL like for example `mypurecloud.com.au` if an ORG is set it will display the ORG Name in the home page Title.

## Navigation

When at the main menu press `?` to see the additional help menu

## Debugging

While I do try to pass status update messages to the TUI there are at times longer messages and more debug messages that can help you. To enable this set the environment variable

```
MCPHEE11_TUI_DEBUG=true
```

This will then create a file `mcphee11-tui-log.log` in the location that you are running the TUI. If you have issues please enable the log file have a look in there and if your still not sure then raise an issue in the repo with log details. Check the log for any flow names etc that you need to mask out first before submitting.

As this is an example TUI and not an official product of any kind I will do my best to assist, I'm also open to PR's if people find issues and want to contribute to this code as well.
