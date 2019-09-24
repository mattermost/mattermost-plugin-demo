module github.com/mattermost/mattermost-plugin-demo

go 1.12

require (
	github.com/blang/semver v3.6.1+incompatible
	github.com/mattermost/mattermost-server v0.0.0-20190911153551-9c3fc4ab9f30
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	// Workaround for https://github.com/golang/go/issues/30831 and fallout.
	github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
)
