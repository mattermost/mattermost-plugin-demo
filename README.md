# Demo Plugin

This plugin demonstrates the capabilities of a Mattermost plugin. It includes the same build scripts as [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample), but implements all supported server-side hooks and registers a component for each supported client-side integration point. See [server/README.md](server/README.md) and [webapp/README.md](webapp/README.md) for more details. The plugin also doubles as a testbed for verifying plugin functionality during release testing.

Feel free to base your own plugin off this repository, removing or modifying components as needed. If you're already familiar with what plugins can do, consider starting from [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample) instead, which includes the same build framework but omits the demo implementations.

Note that this plugin is authored for Mattermost 5.2 and later, and is not compatible with earlier releases of Mattermost.

For details on getting started, see [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample).

## System-wide setting implementation example

If you want to propagate server-side settings changes to the webapp-side of your plugin, you can implement this by emitting a custom WebSocket event from server-side and handle the event on the webapp-side.

There is an example implementation in `server/configuration.go` and `webapp/plugin.jsx`, using a WebSocket event named `system_wide_setting_changed`. 
It displays a modal with the setting's new value when it is changed from `admin_console`.   
