---
layout: "influxdb"
page_title: "InfluxDB: influxdb_database"
sidebar_current: "docs-influxdb-resource-database"
description: |-
  The influxdb_database resource allows an InfluxDB database to be created.
---

# influxdb\_database

The database resource allows a database to be created on an InfluxDB server.

## Example Usage

```hcl
resource "influxdb_database" "metrics" {
    name = "awesome_app"
}

resource "influxdb_database" "metrics_aggregation" {
  name = "testdb11"
  retention_policies = [
    {
      name = "52weeks",
      duration = "52w"
      default = "true"
    },
    {
      name = "104weeks",
      duration = "104w"
      shardgroupduration = "3d"
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the database. This must be unique on the
  InfluxDB server.
* `retention_policies` - (Optional) A list of retention policies for specified database

Each `retention_policies` supports the following:

* `name` - (Required) The name of the retention policy
* `duration` - (Required) The duration for retention policy, format of duration can be found at InfluxDB Documentation.
* `replication` - (Optional) Determines how many copies of data points are stored in a cluster. Not applicable for single node / Open Source version of InfluxDB. Default value of 1.
* `shardgroupduration` - (Optional) Determines how much time each shard group spans. How and why to modify can be found at InfluxDB Documentation.
* `default` - (Optional) Marks current retention policy as default. Default value is false.

## Attributes Reference

This resource exports no further attributes.
