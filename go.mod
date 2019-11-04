module github.com/mattermost/mattermost-plugin-demo

go 1.12

require (
	github.com/a8m/mark v0.1.1-0.20170507133748-44f2db618845 // indirect
	github.com/blang/semver v3.6.1+incompatible
	github.com/gernest/wow v0.1.0 // indirect
	github.com/go-ldap/ldap v3.0.3+incompatible // indirect
	github.com/go-redis/redis v6.15.2+incompatible // indirect
	github.com/mattermost/mattermost-server v0.0.0-20191030173614-4111ffdb4db7
	github.com/minio/cli v1.20.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
)

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	// Workaround for https://github.com/golang/go/issues/30831 and fallout.
	github.com/golang/lint => github.com/golang/lint v0.0.0-20190227174305-8f45f776aaf1
)
