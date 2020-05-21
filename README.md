# Plug 

## Installation

You'll need to have the `go` toolchain setup correctly.
Once you've done this, clone this repo and run `go install`. 

## Usage

Build your plugin using `gradle releaseBundle`. This command will
place your bundled plugin and a JSON metadata file into `build/distributions`
(relative to your plugin directory root). 

This tool has one command: `plug serve --plugin-dir <path-to-build/distributions>`.
This will start a server for your plugin.

Inside your service config, you'll need to point to this server:

```
spinnaker:
  extensibility:
    plugins:
      Armory.RandomWaitPlugin:
        enabled: true
        extensions:
          armory.randomWaitStage:
            enabled: true
            config:
              maxWaitTime: 10
    repositories:
      examplePluginRepo:
        url: http://localhost:9001/plugins.json
```
