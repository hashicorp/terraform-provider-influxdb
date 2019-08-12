## 1.3.1 (Unreleased)
## 1.3.0 (August 12, 2019)

IMPROVEMENTS:

* Support for optional shard group duration in database configuration

## 1.2.0 (May 22, 2019)

IMPROVEMENTS:

* The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

## 1.1.0 (March 07, 2019)

IMPROVEMENTS:

* **Quote-safe usernames:** Provider will now work fine with username having quotes

## 1.0.3 (January 29, 2019)

IMPROVEMENTS:

* **Password hashing:** Provider will no longer store passwords in state files as plaintext

DOCUMENTATION UPDATES:

* Updated "Skip SSL Verify" documentation to match actual implementation

## 1.0.2 (September 05, 2018)

IMPROVEMENTS:

* **Provider Configuration:** Added support for skip TLS check on provider init ([#12](https://github.com/terraform-providers/terraform-provider-influxdb/issues/11))

## 1.0.1 (July 19, 2018)

Bug Fixes:

* Unapplied list of grants changes detected based on ordering ([#6](https://github.com/terraform-providers/terraform-provider-influxdb/issues/6))

## 1.0.0 (February 09, 2018)

IMPROVEMENTS:

* **resource_database:** Added support for retention policy configuration ([#3](https://github.com/terraform-providers/terraform-provider-influxdb/issues/3))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
