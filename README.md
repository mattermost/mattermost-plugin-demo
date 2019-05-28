# Demo Plugin ![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-demo/master.svg)

This plugin demonstrates the capabilities of a Mattermost plugin. It includes the same build scripts as [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample), but implements all supported server-side hooks and registers a component for each supported client-side integration point. See [server/README.md](server/README.md) and [webapp/README.md](webapp/README.md) for more details. The plugin also doubles as a testbed for verifying plugin functionality during release testing.

Feel free to base your own plugin off this repository, removing or modifying components as needed. If you're already familiar with what plugins can do, consider starting from [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample) instead, which includes the same build framework but omits the demo implementations.

This plugin is compatible with Mattermost server version 5.10 and later. For Mattermost versions 5.2 - 5.9, you may download [v0.0.5 of the Demo Plugin](https://github.com/mattermost/mattermost-plugin-demo/releases/tag/v0.0.5).

For details on getting started, see [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample).
