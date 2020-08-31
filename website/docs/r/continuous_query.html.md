---
layout: "influxdb"
page_title: "InfluxDB: influxdb_continuous_query"
sidebar_current: "docs-influxdb-resource-continuous_query"
description: |-
  The influxdb_continuous_query resource allows an InfluxDB continuous query to be managed.
---

# influxdb\_continuous\_query

The continuous_query resource allows a continuous query to be created on an InfluxDB server.

## Example Usage

```hcl
resource "influxdb_database" "test" {
    name = "terraform-test"
}

resource "influxdb_continuous_query" "minnie" {
    name = "minnie"
    database = "${influxdb_database.test.name}"
    query = "SELECT min(mouse) INTO min_mouse FROM zoo GROUP BY time(30m)"
}

resource "influxdb_continuous_query" "minnie_resample" {
    name = "minnie_resample"
    database = "${influxdb_database.test.name}"
    query = "SELECT min(mouse) INTO min_mouse_resample FROM zoo GROUP BY time(30m)"
    resample = "EVERY 30m FOR 2h"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name for the continuous_query. This must be unique on the InfluxDB server.
* `database` - (Required) The database for the continuous_query. This must be an existing influxdb database.
* `query` - (Required) The query for the continuous_query.
* `resample` - (Optional) The body of the query's RESAMPLE clause. The format is detailed in the InfluxDB documentation.

## Attributes Reference

This resource exports no further attributes.
