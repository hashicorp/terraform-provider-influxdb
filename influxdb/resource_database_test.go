package influxdb

import (
	"fmt"
	"log"
	"testing"

	"encoding/json"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/influxdata/influxdb/client"
)

func TestAccInfluxDBDatabase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseExists("influxdb_database.test"),
					resource.TestCheckResourceAttr(
						"influxdb_database.test", "name", "terraform-test",
					),
				),
			},
		},
	})
}

func TestAccInfluxDBDatabaseWithRPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatabaseWithRPSConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseExists("influxdb_database.rptest"),
					resource.TestCheckResourceAttr(
						"influxdb_database.rptest", "name", "terraform-rp-test",
					),
					testAccCheckRetentionPolicy("influxdb_database.rptest", "terraform-rp-test", "1day", "24h0m0s", "1", true),
					testAccCheckRetentionPolicy("influxdb_database.rptest", "terraform-rp-test", "52weeks", "8736h0m0s", "1", false),
				),
			},
			resource.TestStep{
				Config: testAccDatabaseWithRPSUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseExists("influxdb_database.rptest"),
					resource.TestCheckResourceAttr(
						"influxdb_database.rptest", "name", "terraform-rp-test",
					),
					testAccCheckRetentionPolicyNonExisting("influxdb_database.rptest", "name", "52weeks"),
					testAccCheckRetentionPolicy("influxdb_database.rptest", "terraform-rp-test", "1day", "24h0m0s", "1", false),
					testAccCheckRetentionPolicy("influxdb_database.rptest", "terraform-rp-test", "12weeks", "2016h0m0s", "1", true),
				),
			},
		},
	})
}

func testAccCheckDatabaseExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No database id set")
		}

		conn := testAccProvider.Meta().(*client.Client)

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
			if result[0] == rs.Primary.Attributes["name"] {
				return nil
			}
		}

		return fmt.Errorf("Database %q does not exist", rs.Primary.Attributes["name"])
	}
}

func testAccCheckRetentionPolicyNonExisting(n, database, policyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user id set")
		}

		conn := testAccProvider.Meta().(*client.Client)

		log.Printf("%s Created command %v", "SIDGOD", fmt.Sprintf("SHOW RETENTION POLICIES ON %s", rs.Primary.Attributes["name"]))

		query := client.Query{
			Command: fmt.Sprintf("SHOW RETENTION POLICIES ON \"%s\"", rs.Primary.Attributes["name"]),
		}

		resp, err := conn.Query(query)
		if err != nil {
			return err
		}

		if resp.Err != nil {
			return resp.Err
		}

		for _, result := range resp.Results[0].Series[0].Values {
			if result[0].(string) == policyName {
				return fmt.Errorf("Retention Policy %q on %q for %q exists", policyName, database, rs.Primary.Attributes["name"])
			}
		}

		return nil
	}
}

func testAccCheckRetentionPolicy(n, database, policyName, duration, replication string, isDefault bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user id set")
		}

		conn := testAccProvider.Meta().(*client.Client)

		log.Printf("%s Created command %v", "SIDGOD", fmt.Sprintf("SHOW RETENTION POLICIES ON %s", rs.Primary.Attributes["name"]))

		query := client.Query{
			Command: fmt.Sprintf("SHOW RETENTION POLICIES ON \"%s\"", rs.Primary.Attributes["name"]),
		}

		resp, err := conn.Query(query)
		if err != nil {
			return err
		}

		if resp.Err != nil {
			return resp.Err
		}

		for _, result := range resp.Results[0].Series[0].Values {
			if result[0].(string) == policyName {
				if result[1].(string) != duration {
					log.Printf("%s Result %v", "SIDGOD", result[1])
					return fmt.Errorf("Duration %q on retention Policy %q on %q for %q does not match", duration, policyName, database, rs.Primary.Attributes["name"])
				} else if result[3].(json.Number).String() != replication {
					return fmt.Errorf("Replication %q on retention Policy %q on %q for %q does not match", replication, policyName, database, rs.Primary.Attributes["name"])
				} else if result[4].(bool) != isDefault {
					return fmt.Errorf("Default %v on retention Policy %q on %q for %q does not match", isDefault, policyName, database, rs.Primary.Attributes["name"])
				}
				return nil
			}
		}

		return fmt.Errorf("Retention Policy %q on %q for %q does not exist", policyName, database, rs.Primary.Attributes["name"])
	}
}

var testAccDatabaseConfig = `

resource "influxdb_database" "test" {
    name = "terraform-test"
}

`

var testAccDatabaseWithRPSConfig = `
resource "influxdb_database" "rptest" {
  name = "terraform-rp-test"
  retention_policies = [
    {
      name = "1day",
      duration = "1d"
      default = "true"
    },
    {
      name = "52weeks",
      duration = "52w",
    }
  ]
}
`

var testAccDatabaseWithRPSUpdateConfig = `
resource "influxdb_database" "rptest" {
  name = "terraform-rp-test"
  retention_policies = [
    {
      name = "1day",
      duration = "1d"
    },
    {
      name = "12weeks",
      duration = "12w",
      default = "true"
    },
  ]
}
`
