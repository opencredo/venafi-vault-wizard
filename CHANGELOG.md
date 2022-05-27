# Changelog
All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

## 0.1.3 (2022/05/27)

### Fixed

* Updates go version from 1.16 to 1.18 within release workflow [#82](https://github.com/opencredo/venafi-vault-wizard/pull/82)

## 0.1.2 (2022/05/27)

### Changed
* Updates golang to 1.18
* Updates mockery to 2.10.4
* Updates the IP range used within Vagrant example files
* Bump github.com/hashicorp/vault/sdk from 0.4.1 to 0.5.0 [#79](https://github.com/opencredo/venafi-vault-wizard/pull/79)
* Bump github.com/hashicorp/vault/api from 1.5.0 to 1.6.0 [#78](https://github.com/opencredo/venafi-vault-wizard/pull/78)
* Updates to various build related dependencies

## 0.1.1 (2021/09/09)

### Changed
* Example Vault instances have UI enabled.

### Fixed
* Release binary name changed from venafi-vault-wizard_x.x.x to vvw to match documentation.

## 0.1.0 (2021/09/08)

* Initial release of the Venafi Vault Wizard that can be used to 
  install the [Venafi PKI Monitor](https://github.com/Venafi/vault-pki-monitor-venafi) 
  and/or [Venafi PKI Backend ](https://github.com/Venafi/vault-pki-backend-venafi) Vault Plugins.
* Provides single and multi node Vault examples along with configuration for 
  both Venafi Trust Protection Platform, (TPP) and Venafi as a Service, (VaaS).
* Documentation is provided covering installation and configuration of the Venafi Vault Wizard.