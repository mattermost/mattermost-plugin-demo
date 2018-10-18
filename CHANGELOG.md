# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## 0.0.5 - 2018-10-18
### Added
- Upgraded to Mattermost 5.4
- Partial unit test coverage

## 0.0.4 - 2018-10-01
### Added
- Example configuration usage
- Idiomatic error usage, leveraging support in plugin framework to marshal unregistered error types as strings.
### Removed
- Dependence on configuration auto-unmarshalling provided by Mattermost server.

## 0.0.3 - 2018-09-21
## Added
- Inject plugin version into bundle name
- Expose plugin version from manifest in server and client code

## 0.0.1 - 2018-08-16
### Added
- Initial release
