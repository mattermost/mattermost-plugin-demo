# ====================================================================================
# Build Targets
# ====================================================================================

## Checks the code style, tests, builds and bundles the plugin.
.PHONY: all
all: check-style test dist

## Ensures the plugin manifest is valid
.PHONY: manifest-check
manifest-check:
	pluginctl manifest check

## Cleans the server build artifacts.
.PHONY: clean-server
clean-server:
ifneq ($(HAS_SERVER),)
	rm -rf server/dist
endif

## Builds the server, if it exists, for all supported architectures, unless MM_SERVICESETTINGS_ENABLEDEVELOPER is set.
.PHONY: server
server: clean-server
server:
ifneq ($(HAS_SERVER),)
ifneq ($(MM_DEBUG),)
	$(info DEBUG mode is on; to disable, unset MM_DEBUG)
endif
	mkdir -p server/dist;
ifneq ($(MM_SERVICESETTINGS_ENABLEDEVELOPER),)
	@echo Building plugin only for $(DEFAULT_GOOS)-$(DEFAULT_GOARCH) because MM_SERVICESETTINGS_ENABLEDEVELOPER is enabled
	cd server && env CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-$(DEFAULT_GOOS)-$(DEFAULT_GOARCH);
else
	cd server && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-linux-amd64;
	cd server && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-linux-arm64;
	cd server && env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-darwin-amd64;
	cd server && env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-darwin-arm64;
	cd server && env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-windows-amd64.exe;
endif
endif

## Ensures NPM dependencies are installed without having to run this all the time.
webapp/node_modules: $(wildcard webapp/package.json)
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) install
	touch $@
endif

## Builds the webapp, if it exists.
.PHONY: webapp
webapp: webapp/node_modules
ifneq ($(HAS_WEBAPP),)
ifeq ($(MM_DEBUG),)
	cd webapp && $(NPM) run build;
else
	cd webapp && $(NPM) run debug;
endif
endif

## Generates a tar bundle of the plugin for install.
.PHONY: bundle
bundle:
	rm -rf dist/
	mkdir -p dist/$(PLUGIN_ID)
	cp plugin.json dist/$(PLUGIN_ID)/plugin.json
ifneq ($(wildcard $(ASSETS_DIR)/.),)
	cp -r $(ASSETS_DIR) dist/$(PLUGIN_ID)/
endif
ifneq ($(HAS_PUBLIC),)
	cp -r public dist/$(PLUGIN_ID)/
endif
ifneq ($(HAS_SERVER),)
	mkdir -p dist/$(PLUGIN_ID)/server
	cp -r server/dist dist/$(PLUGIN_ID)/server/
endif
ifneq ($(HAS_WEBAPP),)
	mkdir -p dist/$(PLUGIN_ID)/webapp
	cp -r webapp/dist dist/$(PLUGIN_ID)/webapp/
endif
ifeq ($(shell uname),Darwin)
	cd dist && tar --disable-copyfile -cvzf $(BUNDLE_NAME) $(PLUGIN_ID)
else
	cd dist && tar -cvzf $(BUNDLE_NAME) $(PLUGIN_ID)
endif

	@echo plugin built at: dist/$(BUNDLE_NAME)

## Builds and bundles the plugin.
.PHONY: dist
dist: apply server webapp bundle
