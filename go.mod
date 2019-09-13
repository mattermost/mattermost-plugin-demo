module github.com/mattermost/mattermost-plugin-demo

go 1.12

require (
	github.com/blang/semver v3.6.1+incompatible
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/mattermost/mattermost-server v1.4.1-0.20190911153151-98489b9e67d9
	github.com/minio/minio-go v0.0.0-20190422205105-a8704b60278f // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	// Workaround for https://github.com/golang/go/issues/30831 and fallout.
	github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
)
