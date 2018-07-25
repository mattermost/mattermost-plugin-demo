# Demo Plugin

This plugin demonstrates the capabilities of a Mattermost plugin. It includes the same build scripts as [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample), but implements all supported server-side hooks and registers a component for each supported client-side integration point. See [server/README.md](server/README.md) and [webapp/README.md](webapp/README.md) for more details. The plugin also doubles as a testbed for verifying plugin functionality during release testing.

Feel free to base your own plugin off this repository, removing or modifying components as needed. If you're already familiar with what plugins can do, consider starting from [mattermost-plugin-sample](https://github.com/mattermost/mattermost-plugin-sample) instead, which includes the same build framework but omits the demo implementations.

Note that this plugin is authored for Mattermost 5.2 and later, and is not compatible with earlier releases of Mattermost.

## Getting Started
Shallow clone the repository to a directory matching your plugin name:
```
git clone --depth 1 https://github.com/mattermost/mattermost-plugin-demo com.example.my-plugin
```

Edit `plugin.json` with your `id`, `name`, and `description`:
```
{
    "id": "com.example.my-plugin",
    "name": "My Plugin",
    "description": "A plugin to enhance Mattermost."
}
```

Build your plugin:
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.example.my-plugin.tar.gz
```

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration and [http](https://httpie.org/) to be installed:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065/
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

Alternatively, if you are running your `mattermost-server` out of a sibling directory by the same name, use the `deploy` target alone to  unpack the files into the right directory. You will need to restart your server and manually enable your plugin.

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

## Q&A

### How do I make a server-only or web app-only plugin?

Simply delete the `server` or `webapp` folders and remove the corresponding sections from `plugin.json`. The build scripts will skip the missing portions automatically.

### How do I remove unwanted hooks from the server?

Simply delete the corresponding implementations (or files). The Mattermost server automatically identifies which hooks have been implemented when the plugin is started.
