# Demo Plugin

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-demo/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-demo)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-demo/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-demo)


This plugin demonstrates the capabilities of a Mattermost plugin. It includes the same build scripts as [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template), but implements all supported server-side hooks and registers a component for each supported client-side integration point. See [server/README.md](server/README.md) and [webapp/README.md](webapp/README.md) for more details. The plugin also doubles as a testbed for verifying plugin functionality during release testing.

Feel free to base your own plugin off this repository, removing or modifying components as needed. If you're already familiar with what plugins can do, consider starting from [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template) instead, which includes the same build framework but omits the demo implementations.

Note that this plugin is authored for the Mattermost version indicated in the `min_server_version` within the [plugin.json](https://github.com/mattermost/mattermost-plugin-demo/blob/2461499b06453b7a37d9ca4aedd4d23d24089047/plugin.json#L6), and is not compatible with earlier releases of Mattermost.

For details on getting started, see [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template).
