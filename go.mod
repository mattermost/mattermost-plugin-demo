module github.com/mattermost/mattermost-plugin-demo

go 1.12

require (
	github.com/blang/semver v3.6.1+incompatible
	github.com/mattermost/mattermost-server/v5 v5.3.2-0.20200129194125-99a82ef07ec0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
)

replace github.com/mattermost/mattermost-server/v5 => github.com/ashishbhate/mattermost-server/v5 v5.3.2-0.20200131051305-25037fd7f545
