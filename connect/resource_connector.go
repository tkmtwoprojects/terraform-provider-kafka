package connect

import (
	//"errors"
	//"fmt"
	"log"

	kc "github.com/tkmtwoprojects/go-kafka/connect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func connectConnectorResource() *schema.Resource {
	return &schema.Resource{
		Create: connectorCreate,
		Read:   connectorRead,
		Update: connectorUpdate,
		Delete: connectorDelete,
		Importer: &schema.ResourceImporter{
			State: setNameFromID,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the connector",
			},
			"config": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Sensitive:   true,
				Description: "A map of string k/v attributes all treated as sensitive.",
			},
		},
	}
}

func setNameFromID(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	connectorName := d.Id()
	log.Printf("Import connector with name: %s", connectorName)
	d.Set("name", connectorName)

	return []*schema.ResourceData{d}, nil
}


//
// Create a connector
//
func connectorCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*kc.Client)
	name := d.Get("name").(string)
	config := d.Get("config").(map[string]interface{})
	
	connectorResponse, err := c.Connectors.Create(name, config)
	
	if err == nil {
		d.SetId(connectorResponse.Name)
		d.Set("config", connectorResponse.Config)
		log.Printf("Created the connector '%s'", connectorResponse.Name)
		return connectorRead(d, meta)
	}
	
	return err
}

//
// Read a connector
//
// Not used much outside of verifying or importing.
//
func connectorRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*kc.Client)
	name := d.Get("name").(string)
	
	connectorResponse, err := c.Connectors.Get(name)

	if err == nil {
		log.Printf("Read the connector '%s'", connectorResponse.Name)
		return nil
	}

	return err
}


func connectorUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*kc.Client)
	name := d.Get("name").(string)
	config := d.Get("config").(map[string]interface{})

	connectorResponse, err := c.Connectors.Update(name, config)

	if err == nil {
		d.SetId(connectorResponse.Name)
		d.Set("config", connectorResponse.Config)
		log.Printf("Updated the connector '%s'", connectorResponse.Name)
		return connectorRead(d, meta)
	}
	
	return err
}


func connectorDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*kc.Client)
	name := d.Get("name").(string)

	err := c.Connectors.Delete(name)
	
	if err == nil {
		d.SetId("")
		log.Printf("Deleted the connector '%s'", name)
	}

	return err
}
