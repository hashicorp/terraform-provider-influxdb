package influxdb

import (
	"fmt"
	"strconv"

	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/influxdata/influxdb/client"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: createDatabase,
		Read:   readDatabase,
		Delete: deleteDatabase,
		Update: updateDatabase,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_policies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Required: true,
						},
						"replication": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"shardgroupduration": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"default": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func createDatabase(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	name := d.Get("name").(string)
	queryStr := fmt.Sprintf("CREATE DATABASE %s", quoteIdentifier(name))
	query := client.Query{
		Command: queryStr,
	}

	resp, err := conn.Query(query)
	if err != nil {
		return err
	}
	if resp.Err != nil {
		return resp.Err
	}

	d.SetId(name)

	if v, ok := d.GetOk("retention_policies"); ok {
		retentionPolicies := v.([]interface{})
		for _, vv := range retentionPolicies {
			retentionPolicy := vv.(map[string]interface{})
			if err := createRetentionPolicy(conn, retentionPolicy["name"].(string), retentionPolicy["duration"].(string), retentionPolicy["replication"].(int), retentionPolicy["shardgroupduration"].(string), retentionPolicy["default"].(bool), name); err != nil {
				return err
			}
		}
	}

	return nil
}

func createRetentionPolicy(conn *client.Client, policyName string, duration string, replication int, shardGroupDuration string, defaultPolicy bool, database string) error {
	var shardDuration string

	if shardGroupDuration != "" {
		shardDuration = fmt.Sprintf("SHARD DURATION %s ", shardGroupDuration)
	}

	if defaultPolicy {
		return exec(conn, fmt.Sprintf("CREATE RETENTION POLICY %s ON %s DURATION %s REPLICATION %d %s DEFAULT", quoteIdentifier(policyName), quoteIdentifier(database), duration, replication, shardDuration))
	} else {
		return exec(conn, fmt.Sprintf("CREATE RETENTION POLICY %s ON %s DURATION %s REPLICATION %d %s", quoteIdentifier(policyName), quoteIdentifier(database), duration, replication, shardDuration))
	}
}

func updateRetentionPolicy(conn *client.Client, policyName string, duration string, replication int, shardGroupDuration string, defaultPolicy bool, database string) error {
	var shardDuration string

	if shardGroupDuration != "" {
		shardDuration = fmt.Sprintf("SHARD DURATION %s ", shardGroupDuration)
	}

	if defaultPolicy {
		return exec(conn, fmt.Sprintf("ALTER RETENTION POLICY %s ON %s DURATION %s REPLICATION %d %s DEFAULT", quoteIdentifier(policyName), quoteIdentifier(database), duration, replication, shardDuration))
	} else {
		return exec(conn, fmt.Sprintf("ALTER RETENTION POLICY %s ON %s DURATION %s REPLICATION %d %s", quoteIdentifier(policyName), quoteIdentifier(database), duration, replication, shardDuration))
	}
}

func deleteRetentionPolicy(conn *client.Client, policyName string, database string) error {
	return exec(conn, fmt.Sprintf("DROP RETENTION POLICY %s ON %s", quoteIdentifier(policyName), quoteIdentifier(database)))
}

func readDatabase(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	name := d.Id()

	// InfluxDB doesn't have a command to check the existence of a single
	// database, so we instead must read the list of all databases and see
	// if ours is present in it.
	query := client.Query{
		Command: "SHOW DATABASES",
	}

	resp, err := conn.Query(query)
	if err != nil {
		return err
	}
	if resp.Err != nil {
		return resp.Err
	}

	for _, result := range resp.Results[0].Series[0].Values {
		if result[0] == name {
			readRetentionPolicies(d, meta)
			return nil
		}
	}

	// If we fell out here then we didn't find our database in the list.
	d.SetId("")

	return nil
}

func readRetentionPolicies(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	name := d.Get("name").(string)

	query := client.Query{
		Command: fmt.Sprintf("SHOW RETENTION POLICIES ON %s", name),
	}

	resp, err := conn.Query(query)
	if err != nil {
		return err
	}

	if resp.Err != nil {
		return resp.Err
	}

	var retentionPolicies = []map[string]string{}

	if resp.Results[0].Err == nil {
		for _, result := range resp.Results[0].Series[0].Values {
			var retentionPolicy = map[string]string{
				"name":               result[0].(string),
				"duration":           result[1].(string),
				"shardGroupDuration": result[2].(string),
				"replicaN":           result[3].(json.Number).String(),
				"default":            strconv.FormatBool(result[4].(bool)),
			}
			retentionPolicies = append(retentionPolicies, retentionPolicy)
		}
	}

	d.Set("retention_policies", retentionPolicies)
	return nil
}

func deleteDatabase(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	name := d.Id()

	queryStr := fmt.Sprintf("DROP DATABASE %s", quoteIdentifier(name))
	query := client.Query{
		Command: queryStr,
	}

	resp, err := conn.Query(query)
	if err != nil {
		return err
	}
	if resp.Err != nil {
		return resp.Err
	}

	d.SetId("")

	return nil
}

func updateDatabase(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	name := d.Get("name").(string)

	if d.HasChange("retention_policies") {
		oldRPVal, newRPVal := d.GetChange("retention_policies")
		oldRPs := oldRPVal.([]interface{})
		newRPs := newRPVal.([]interface{})

		newRPMap := make(map[string]bool)
		oldRPMap := make(map[string]bool)

		for _, newRP := range newRPs {
			newPolicy := newRP.(map[string]interface{})
			policyName := newPolicy["name"].(string)
			newRPMap[policyName] = true
		}

		// RPs in old map but not in new map should be deleted, we'll also create old policies while we are at it
		for _, oldRP := range oldRPs {
			oldPolicy := oldRP.(map[string]interface{})
			policyName := oldPolicy["name"].(string)
			oldRPMap[policyName] = true

			if !newRPMap[policyName] {
				if err := deleteRetentionPolicy(conn, policyName, name); err != nil {
					return err
				}
			}
		}

		for _, newRP := range newRPs {
			newPolicy := newRP.(map[string]interface{})
			policyName := newPolicy["name"].(string)

			// If policy is not in old map, it has to be created newly, otherwise it has to be updated
			if !oldRPMap[policyName] {
				if err := createRetentionPolicy(conn, policyName, newPolicy["duration"].(string), newPolicy["replication"].(int), newPolicy["shardgroupduration"].(string), newPolicy["default"].(bool), name); err != nil {
					return err
				}
			} else {
				if err := updateRetentionPolicy(conn, policyName, newPolicy["duration"].(string), newPolicy["replication"].(int), newPolicy["shardgroupduration"].(string), newPolicy["default"].(bool), name); err != nil {
					return err
				}
			}
		}
	}

	return readDatabase(d, meta)
}
