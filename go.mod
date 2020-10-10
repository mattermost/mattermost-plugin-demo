module github.com/mattermost/mattermost-plugin-demo

go 1.12

require (
	github.com/blang/semver v3.6.1+incompatible
	github.com/mattermost/mattermost-plugin-api v0.0.9
	github.com/mattermost/mattermost-server/v5 v5.25.1
	github.com/mholt/archiver/v3 v3.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
)

replace github.com/mattermost/mattermost-server/v5 v5.25.1 => github.com/shieldsjared/mattermost-server/v5 v5.0.0-20201008141002-56ee0d27d54b
