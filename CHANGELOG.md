# Changelog

## [0.46.0](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.45.1...v0.46.0) (2026-04-16)


### Features

* add pfsense_rest_api_key resource for REST API v2 key management ([13c5db1](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/13c5db16b772c680fdf245540b79cbdf6c71e1d5))

## [0.45.1](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.45.0...v0.45.1) (2026-04-08)


### Bug Fixes

* harden PHP command execution with session recovery, login detection, and pkg database fallback ([01bdc44](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/01bdc44053e028973bcc4c87f9a25be21340cc44))

## [0.45.0](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.44.0...v0.45.0) (2026-04-08)


### Features

* add pfsense_certificate_authority resource and data sources for native CA management ([ad20690](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/ad20690ad0462dc0a4f0271aae4f6f27a205bc0e))

## [0.44.0](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.43.3...v0.44.0) (2026-04-08)


### Features

* add pfsense_package resource and data sources for native package management ([ee40956](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/ee40956f7ff8706ec27fd2cebb8634bb540fd4a1))

## [0.43.3](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.43.2...v0.43.3) (2026-04-02)


### Bug Fixes

* escape user values in PHP commands to prevent parse errors ([5724393](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/57243937bcc56ae49ec5cfab80c67c848222a0f5))

## [0.43.2](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.43.1...v0.43.2) (2026-04-02)


### Bug Fixes

* improve IPsec Phase 2 reliability and fix keepalive handling ([50613e8](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/50613e846535f3d54835f08881945a09b341fc30))

## [0.43.1](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.43.0...v0.43.1) (2026-04-02)


### Bug Fixes

* handle SSH config as both object and string in system advanced admin response ([53b451a](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/53b451a7ae669ee624542b3fbc8c44ee4d176ec0))

## [0.43.0](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.42.0...v0.43.0) (2026-04-02)


### Features

* add IPsec Phase 1 and Phase 2 resources and data sources ([4e87512](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/4e87512cfb54bc9baae47a28756b2e142d3cb312))

## [0.42.0](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.41.0...v0.42.0) (2026-04-02)


### Features

* add DNS resolver general and advanced singleton resources and data sources ([b933cc8](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/b933cc8628685b57cd8c8a37f261f84eecec59dd))

## [0.41.0](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.40.1...v0.41.0) (2026-04-01)


### Features

* add DHCPv4 server resource and data source for per-interface DHCP configuration ([21bad8e](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/21bad8ebb04b7f7c3d579c07daa809d2561048d3))

## [0.40.1](https://github.com/emanuelbesliu/terraform-provider-pfsense/compare/v0.40.0...v0.40.1) (2026-04-01)


### Bug Fixes

* make GPG passphrase optional in release workflows ([4cd98ac](https://github.com/emanuelbesliu/terraform-provider-pfsense/commit/4cd98aceeb596a715b89cd43fb67fa3c97f86841))

## Changelog
