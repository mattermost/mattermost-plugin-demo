# Demo Plugin

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-demo/master.svg)](https://circleci.com/gh/mattermost/mattermost-plugin-demo)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-demo/master.svg)](https://codecov.io/gh/mattermost/mattermost-plugin-demo)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-demo)](https://github.com/mattermost/mattermost-plugin-demo/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-demo/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-demo/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

**Maintainer:** [@hanzei](https://github.com/hanzei)
**Co-Maintainer:** [@jfrerich](https://github.com/jfrerich)

This plugin demonstrates the capabilities of a Mattermost plugin. It includes the same build scripts as [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template), but implements all supported server-side hooks and registers a component for each supported client-side integration point. See [server/README.md](server/README.md) and [webapp/README.md](webapp/README.md) for more details. The plugin also doubles as a testbed for verifying plugin functionality during release testing.

Once installed and enabled, you can specify both the channel and user for the demo plugin. If the specified channel or user doesn't exist, the plugin creates it for you.

Feel free to base your own plugin off this repository, removing or modifying components as needed. If you're already familiar with what plugins can do, consider starting from [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template) instead, which includes the same build framework but omits the demo implementations.

Note that this plugin is authored for the Mattermost version indicated in the `min_server_version` within [plugin.json](https://github.com/mattermost/mattermost-plugin-demo/blob/master/plugin.json), and is not compatible with earlier releases of Mattermost.

For details on getting started, see [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template).

## Releasing this plugin

A new minor version of this plugin is released with every feature release of Mattermost. The new version should be cut until [Code complete](https://docs.mattermost.com/process/feature-release.html#f-t-minus-14-working-days-code-complete).
