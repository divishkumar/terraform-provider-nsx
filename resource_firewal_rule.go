package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/gonsx"
	"github.com/sky-uk/gonsx/api/distributedfw/fwrules"
	"log"
	"github.com/sky-uk/gonsx/api/distributedfw/sections"
)


func resourceGetSectionTimestamp(sectionID int, sectionType string, m interface{}) *sections.Section {
	sectionNSXClient := m.(*gonsx.NSXClient)
	sectionTimestamp := sections.GetSectionTimestamp(sectionID, sectionType)
	sectErr := sectionNSXClient.Do(sectionTimestamp)
	if sectErr != nil {
		log.Println("could not get timestamp")
	}
	log.Println(sectionTimestamp.GetResponse())
	return sectionTimestamp.GetResponse()
}

func resourceFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Update: resourceFirewallRuleUpdate,
		Delete: resourceFirewallRuleDelete,
		Schema: map[string]*schema.Schema{
			"ruleid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "A name for the  rule",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				ForceNew:    false,
				Description: "determines if the rule is disabled or not",
			},
			"ruletype": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of rule, valid values are LAYER2 and LAYER3",
			},
			"logged": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Should this rule be logged",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "What to do with the packets that match this rule, allow,drop, etc",
			},
			"source": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "Source of traffic for the firewall rule, it could be, CDIR, IP Set, IPv4 addresses, Virtual Machine, Vnic, Security Group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the source",
						},
						"type": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Type of source",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the source",
						},
						"isvalid": {
							Optional:    true,
							Type:        schema.TypeBool,
							ForceNew:    false,
							Description: "Is the source valid",
						},
					},
				},
			},
			"destination": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "Source of traffic for the firewall rule, it could be, CDIR, IP Set, IPv4 addresses, Virtual Machine, Vnic, Security Group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the source",
						},
						"type": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Type of source",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the source",
						},
						"isvalid": {
							Optional:    true,
							Type:        schema.TypeBool,
							ForceNew:    false,
							Description: "Is the source valid",
						},
					},
				},
			},
			"service": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Name of the service",
						},
						"value": {
							Optional:    true,
							Type:        schema.TypeString,
							ForceNew:    false,
							Description: "Value of the service",
						},
						"dstport": {
							Optional:    true,
							Type:        schema.TypeInt,
							ForceNew:    false,
							Description: "Destination port for the service",
						},
						"protocol": {
							Optional:    true,
							Type:        schema.TypeInt,
							ForceNew:    false,
							Description: "Protocol id ",
						},
						"subprotocol": {
							Optional:    true,
							Type:        schema.TypeInt,
							ForceNew:    false,
							Description: "SubProtocol id ",
						},
					},
				},
			},
			"sectionid": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Id for the section where the rule bellongs",
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Direction for the traffic",
			},
			"packettype": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Type of packets",
			},
		},
	}

}

func resourceFirewallRuleCreate(d *schema.ResourceData, m interface{}) error {
	nsxclient := m.(*gonsx.NSXClient)
	var fwRule fwrules.Rule

	if v, ok := d.GetOk("name"); ok {
		fwRule.Name = v.(string)
	} else {
		return fmt.Errorf("Name parameter is required")
	}

	if v, ok := d.GetOk("disabled"); ok {
		fwRule.Disabled = v.(bool)
	}

	if v, ok := d.GetOk("ruletype"); ok {
		fwRule.RuleType = v.(string)
	} else {
		return fmt.Errorf("Rule Type is required")
	}

	if v, ok := d.GetOk("logged"); ok {
		fwRule.Logged = v.(string)
	}

	if v, ok := d.GetOk("action"); ok {
		fwRule.Action = v.(string)
	} else {
		return fmt.Errorf("Error needs to be set")
	}


	if v, ok := d.GetOk("source"); ok {
		if sources, ok := v.(*schema.Set); ok {
			for _, source := range sources.List() {
				sourceObject := source.(map[string]interface{})
				newSource := fwrules.Source{}

				if sourceNameValue, ok := sourceObject["name"]; ok {
					newSource.Name = sourceNameValue.(string)
				}

				if sourceTypeValue, ok := sourceObject["type"]; ok {
					newSource.Type = sourceTypeValue.(string)
				}

				if sourceValue, ok := sourceObject["value"]; ok {
					newSource.Value = sourceValue.(string)
				}

				if sourceValidValue, ok := sourceObject["isvalid"]; ok {
					newSource.IsValid = sourceValidValue.(bool)
				}

				fwRule.Sources = append(fwRule.Sources, newSource)
			}
		}
	} else {
		return fmt.Errorf("Source  is required")
	}

	if v, ok := d.GetOk("destination"); ok {
		if destinations, ok := v.(*schema.Set); ok {
			for _, destination := range destinations.List() {
				destinationObject := destination.(map[string]interface{})
				newDestination := fwrules.Destination{}

				if destinationNameValue, ok := destinationObject["name"]; ok {
					newDestination.Name = destinationNameValue.(string)
				}

				if destinationTypeValue, ok := destinationObject["type"]; ok {
					newDestination.Type = destinationTypeValue.(string)
				}

				if destinationValue, ok := destinationObject["value"]; ok {
					newDestination.Value = destinationValue.(string)
				}

				if destinationValidValue, ok := destinationObject["isvalid"]; ok {
					newDestination.IsValid = destinationValidValue.(bool)
				}

				fwRule.Destinations = append(fwRule.Destinations, newDestination)
			}
		}
	}

	if v, ok := d.GetOk("service"); ok {
		if services, ok := v.(*schema.Set); ok {
			for _, service := range services.List() {
				serviceObject := service.(map[string]interface{})
				newService := fwrules.Service{}

				if serviceNameValue := serviceObject["name"]; ok {
					newService.Name = serviceNameValue.(string)

				}

				if serviceValue := serviceObject["value"]; ok {
					newService.Value = serviceValue.(string)
				}

				if serviceDstPortValue := serviceObject["dstport"]; ok {
					newService.DestinationPort = serviceDstPortValue.(int)
				}

				if serviceProtocolValue := serviceObject["protocol"]; ok {
					newService.Protocol = serviceProtocolValue.(int)
				}

				if serviceSubProtocolValue := serviceObject["subprotocol"]; ok {
					newService.SubProtocol = serviceSubProtocolValue.(int)
				}
				fwRule.Services = append(fwRule.Services, newService)

			}
		}

	}

	if v, ok := d.GetOk("sectionid"); ok {

		fwRule.SectionID = v.(int)

	} else {
		return fmt.Errorf("Section ID is required")
	}

	if v, ok := d.GetOk("direction"); ok {
		fwRule.Direction = v.(string)
	} else {
		return fmt.Errorf("Direction is required")
	}
	/*nsxMutexKV.Lock(fwRule.EdgeID)
	defer nsxMutexKV.Unlock(fwRule.EdgeID)*/
	createFWRuleAPI := fwrules.NewCreate(fwRule)
	timeStampCall := resourceGetSectionTimestamp(fwRule.SectionID, fwRule.RuleType, m )
	nsxclient.SetHeader("If-Match", timeStampCall.Timestamp)
	log.Println(timeStampCall.Timestamp)
	log.Println(fwRule.SectionID)
	log.Println(fwRule.RuleType)
	createErr := nsxclient.Do(createFWRuleAPI)
	if createErr != nil {
		return fmt.Errorf("Could not create firewall rule")
	}

	if createFWRuleAPI.StatusCode() != 201 {
		log.Println(createFWRuleAPI.StatusCode())
		log.Println("COULD NOT CREATE !!!!!!!")
		return fmt.Errorf("COULD NOT CREATE !!!!!!!")
	}

	d.SetId(d.Get("ruleid").(string))
	return nil

}

func resourceFirewallRuleRead(d *schema.ResourceData, m interface{}) error {
	/*nsxclient := m.(*gonsx.NSXClient)
	ReadRule = fwrules.NewGetSingle()*/
	return nil

}

func resourceFirewallRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return nil

}

func resourceFirewallRuleDelete(d *schema.ResourceData, m interface{}) error {
	return nil

}
