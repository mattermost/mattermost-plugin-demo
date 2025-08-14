# ====================================================================================
# Deployment and Plugin Management
# ====================================================================================

## Applies the plugin manifest to the server and webapp codebase
.PHONY: apply
apply:
	pluginctl manifest apply

## Builds and installs the plugin to a server.
.PHONY: deploy
deploy: dist
	pluginctl deploy --bundle-path dist/$(BUNDLE_NAME)

## Builds and installs the plugin to a server, updating the webapp automatically when changed.
.PHONY: watch
watch: apply server bundle
ifeq ($(MM_DEBUG),)
	cd webapp && $(NPM) run build:watch
else
	cd webapp && $(NPM) run debug:watch
endif

## Installs a previous built plugin with updated webpack assets to a server.
.PHONY: deploy-from-watch
deploy-from-watch: bundle
	pluginctl deploy --bundle-path dist/$(BUNDLE_NAME)

## Disable the plugin.
.PHONY: disable
disable: detach
	pluginctl disable

## Enable the plugin.
.PHONY: enable
enable:
	pluginctl enable

## Reset the plugin, effectively disabling and re-enabling it on the server.
.PHONY: reset
reset: detach
	pluginctl reset

## View plugin logs.
.PHONY: logs
logs:
	pluginctl logs

## Watch plugin logs.
.PHONY: logs-watch
logs-watch:
	pluginctl logs --watch
