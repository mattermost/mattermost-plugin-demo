# Ensure that go is installed. Note that this is independent of whether or not a server is being
# built, since the build script itself uses go.
ifeq ($(GO),)
    $(error "go is not available: see https://golang.org/doc/install")
endif

# Gather build variables to inject into the manifest tool
BUILD_HASH_SHORT = $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_TAG_LATEST = $(shell git describe --tags --match 'v*' --abbrev=0 2>/dev/null)
BUILD_TAG_CURRENT = $(shell git tag --points-at HEAD 2>/dev/null)

# Extract the plugin id from the manifest.
PLUGIN_ID ?= $(shell pluginctl manifest get '{{.Id}}')
ifeq ($(PLUGIN_ID),)
    $(error "Cannot parse id from $(MANIFEST_FILE)")
endif

# Extract the plugin version from the manifest.
PLUGIN_VERSION ?= $(shell pluginctl manifest get '{{.Version}}')
ifeq ($(PLUGIN_VERSION),)
    $(error "Cannot parse version from $(MANIFEST_FILE)")
endif

# Determine if a server is defined in the manifest.
HAS_SERVER ?= $(shell pluginctl manifest get '{{.HasServer}}')

# Determine if a webapp is defined in the manifest.
HAS_WEBAPP ?= $(shell pluginctl manifest get '{{.HasWebapp}}')

# Determine if a /public folder is in use
HAS_PUBLIC ?= $(wildcard public/.)

# Determine if the mattermost-utilities repo is present
HAS_MM_UTILITIES ?= $(wildcard $(MM_UTILITIES_DIR)/.)

# Store the current path for later use
PWD ?= $(shell pwd)

# Ensure that npm (and thus node) is installed.
ifneq ($(HAS_WEBAPP),)
ifeq ($(NPM),)
    $(error "npm is not available: see https://www.npmjs.com/get-npm")
endif
endif

BUNDLE_NAME ?= $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz
