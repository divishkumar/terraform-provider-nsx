package main

import (
	"fmt"
        "github.com/sky-uk/gonsx"
        "github.com/sky-uk/gonsx/api/dhcprelay"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func getAllDhcpRelays(edgeId string, nsxclient *gonsx.NSXClient) (*dhcprelay.DhcpRelay, error) {
        //
        // Get All DHCP Relay agents.
        //
        api := dhcprelay.NewGetAll(edgeId)
        // make the api call with nsxclient
        err := nsxclient.Do(api)
        // check if we err otherwise read response.
        if err != nil {
                //fmt.Println("Error:", err)
                //return nil, err
                return nil, err
        } else {
                log.Println("Get All Response: ", api.GetResponse())
                return api.GetResponse(), nil
        }
}

func resourceDHCPRelay() *schema.Resource {
	return &schema.Resource{
		Create: resourceDHCPRelayCreate,
		Read:   resourceDHCPRelayRead,
		Delete: resourceDHCPRelayDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				  Type:     schema.TypeString,
				  Required: true,
				  ForceNew: true,
			},

			"edgeid": {
				  Type:     schema.TypeString,
				  Required: true,
				  ForceNew: true,
			},

			"vnicindex": {
				     Type:     schema.TypeString,
				     Required: true,
				     ForceNew: true,
			},

			"giaddress": {
				     Type:     schema.TypeString,
				     Required: true,
				     ForceNew: true,
			},

			"dhcpserverip": {
				        Type:     schema.TypeString,
				        Required: true,
				        ForceNew: true,
			},
		},
	}
}

func resourceDHCPRelayCreate(d *schema.ResourceData, m interface{}) error {
        nsxclient := m.(*gonsx.NSXClient)
        var name, edgeid, vnicindex, giaddress, dhcpserverip string

        // Gather the attributes for the resource.
        if v, ok := d.GetOk("name"); ok {
                name = v.(string)
        } else {
                return fmt.Errorf("name argument is required")
        }

        if v, ok := d.GetOk("edgeid"); ok {
                edgeid = v.(string)
        } else {
                return fmt.Errorf("edgeid argument is required")
        }

        if v, ok := d.GetOk("vnicindex"); ok {
                vnicindex = v.(string)
        } else {
                return fmt.Errorf("vnicindex argument is required")
        }

        if v, ok := d.GetOk("giaddress"); ok {
                giaddress = v.(string)
        } else {
                return fmt.Errorf("giaddress argument is required")
        }

        if v, ok := d.GetOk("dhcpserverip"); ok {
                dhcpserverip = v.(string)
        } else {
                return fmt.Errorf("dhcpserverip argument is required")
        }

        // Create the API, use it and check for errors.
        log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.getAllDhcpRelays(%s, %s)", edgeid, nsxclient))
        currentDHCPRelay, err := getAllDhcpRelays(edgeid, nsxclient)

        if err != nil {
                return fmt.Errorf("Error:", err)
        }

        log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.RelayAgent(%s, %s)", vnicindex, giaddress))
        newRelayAgent := dhcprelay.RelayAgent{VnicIndex: vnicindex, GiAddress: giaddress}

        log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.append(%s, %s)", currentDHCPRelay.RelayAgents, newRelayAgent))
        newRelayAgentsList := append(currentDHCPRelay.RelayAgents, newRelayAgent)

        log.Printf(fmt.Sprintf("[DEBUG] dhcprelay.NewUpdate(%s, %s, %s)", dhcpserverip, edgeid, newRelayAgentsList))
        update_api := dhcprelay.NewUpdate(dhcpserverip, edgeid, newRelayAgentsList)

        err = nsxclient.Do(update_api)
        if err != nil {
                return fmt.Errorf("Error:", err)
        } else if update_api.StatusCode() != 204 {
                return fmt.Errorf("Failed to update the DHCP relay %s", update_api.GetResponse())
        }

        // If we get here, everything is OK.  Set the ID for the Terraform state
        // and return the response from the READ method.
        d.SetId(name)
        return resourceDHCPRelayRead(d, m)
}

func resourceDHCPRelayRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDHCPRelayDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}