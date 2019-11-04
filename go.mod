module github.com/mattermost/mattermost-plugin-demo

go 1.12

require (
	github.com/blang/semver v3.6.1+incompatible
	github.com/mattermost/mattermost-server v0.0.0-20191030173614-4111ffdb4db7
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
