package connect

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	kc "github.com/tkmtwoprojects/go-kafka/connect"
)

const testConnectConnector_basic = `
resource "connect_connector" "datagenZero" {
  name = "DatagenUsersZero"
  config = {
    "name": "DatagenUsersZero"
		"connector.class": "io.confluent.kafka.connect.datagen.DatagenConnector",
		"key.converter":   "org.apache.kafka.connect.storage.StringConverter",
		"kafka.topic":     "datagen.users.Zero",
		"max.interval":    "5000",
		"quickstart":      "users"
    }
}
`

const testResourceConnector_updateConfig = `
resource "connect_connector" "datagenZero" {
  name = "DatagenUsersZero"
  config = {
    "name": "DatagenUsersZero"
		"connector.class": "io.confluent.kafka.connect.datagen.DatagenConnector",
		"key.converter":   "org.apache.kafka.connect.storage.StringConverter",
		"kafka.topic":     "datagen.users.Zero",
		"max.interval":    "5150",
		"quickstart":      "users"
    }
}
`

func TestAccConnectConnectorBasicCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testConnectConnectorCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testConnectConnector_basic,
				Check:  testConnectConnectorCheckCreate,
			},
		},
	})
}

func testConnectConnectorCheckCreate(s *terraform.State) error {
	log.Print("[DEBUG] Running testConnectConnectorCheckCreate")

	resourceState := s.Modules[0].Resources["connect_connector.datagenZero"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}
	
	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}
	
	name := instanceState.ID
	if name != instanceState.Attributes["name"] {
		return fmt.Errorf("id does not match name")
	}
	
	client := testAccProvider.Meta().(*kc.Client)
	c, err := client.Connectors.Get("DatagenUsersZero")
	if err != nil {
		return err
	}

	maxInterval := c.Config["max.interval"]
	expected := "5000"
	if maxInterval != expected {
		return fmt.Errorf("max.interval should be %s, got %s connector not updated. \n %v", expected, maxInterval, c.Config)
	}

	return nil
}





func testConnectConnectorCheckDestroy(s *terraform.State) error {
	log.Print("[DEBUG] Running testConnectConnectorCheckDestroy")

	resourceState := s.Modules[0].Resources["connect_connector.datagenZero"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}


	client := testAccProvider.Meta().(*kc.Client)
	c, err := client.Connectors.Get("DatagenUsersZero")
	if err != nil {
		return err
	}

	if c.Name != "" {
		return fmt.Errorf("Connector %q still exists", instanceState.ID)
	}
	
	return nil
}

/*
func TestAccConnectorConfigUpdate(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: testResourceConnector_initialConfig,
				Check:  testResourceConnector_initialCheck,
			},
			{
				Config:            testResourceConnector_initialConfig,
				ResourceName:      "connect_connector.datagenZero",
				ImportStateVerify: true,
				ImportState:       true,
			},
			{
				Config: testResourceConnector_updateConfig,
				Check:  testResourceConnector_updateCheck,
			},
		},
	})
}
*/

/*
func testResourceConnector_initialCheck(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["connect_connector.datagenZero"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	name := instanceState.ID

	if name != instanceState.Attributes["name"] {
		return fmt.Errorf("id doesn't match name")
	}

	client := testProvider.Meta().(kc.Client)

	c, err := client.Get("datagenZero")
	if err != nil {
		return err
	}

	maxInterval := c.Config["max.interval"]
	expected := "5000"
	if maxInterval != expected {
		return fmt.Errorf("max.interval should be %s, got %s connector not updated. \n %v", expected, maxInterval, c.Config)
	}

	return nil
}
*/

/*
func testResourceConnector_updateCheck(s *terraform.State) error {
	client := testProvider.Meta().(kc.Client)

	c, err := client.Get("datagenZero")
	if err != nil {
		return err
	}

	maxInterval := c.Config["max.interval"]
	expected := "5150"
	if maxInterval != expected {
		return fmt.Errorf("max.interval should be %s, got %s connector not updated. \n %v", expected, maxInterval, c.Config)
	}

	return nil
}
*/
