package godaddy

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"terraform-provider-st-godaddy/api"
)

var domainInfo api.RegisterDomainInfo

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_API_KEY", nil),
				Description: "GoDaddy API Key.",
			},

			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_API_SECRET", nil),
				Description: "GoDaddy API Secret.",
			},

			"baseurl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.godaddy.com",
				Description: "GoDaddy Base Url(defaults to production).",
			},
			//admin
			"admin_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_fax": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_jobtitle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_namelast": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_namefirst": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_namemiddle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"admin_postcode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},
			//billing
			"billing_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_fax": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_jobtitle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_namelast": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_namefirst": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_namemiddle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"billing_postcode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},
			//reg
			"reg_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_fax": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_jobtitle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_namelast": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_namefirst": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_namemiddle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"reg_postcode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			//tech
			"tech_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_fax": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_jobtitle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_namelast": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_namefirst": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_namemiddle": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_city": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_country": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},

			"tech_postcode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "GoDaddy Base Url(defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"godaddy_domain_record": resourceDomainRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Key:     d.Get("key").(string),
		Secret:  d.Get("secret").(string),
		BaseURL: d.Get("baseurl").(string),
	}

	//admin
	domainInfo.ContactAdmin.Email = d.Get("admin_email").(string)
	domainInfo.ContactAdmin.Fax = d.Get("admin_fax").(string)
	domainInfo.ContactAdmin.JobTitle = d.Get("admin_jobtitle").(string)
	domainInfo.ContactAdmin.NameLast = d.Get("admin_namelast").(string)
	domainInfo.ContactAdmin.NameFirst = d.Get("admin_namefirst").(string)
	domainInfo.ContactAdmin.NameMiddle = d.Get("admin_namemiddle").(string)
	domainInfo.ContactAdmin.Organization = d.Get("admin_organization").(string)
	domainInfo.ContactAdmin.Phone = d.Get("admin_phone").(string)
	domainInfo.ContactAdmin.AddressMailing.Address1 = d.Get("admin_address").(string)
	domainInfo.ContactAdmin.AddressMailing.City = d.Get("admin_city").(string)
	domainInfo.ContactAdmin.AddressMailing.Country = d.Get("admin_country").(string)
	domainInfo.ContactAdmin.AddressMailing.State = d.Get("admin_state").(string)
	domainInfo.ContactAdmin.AddressMailing.PostalCode = d.Get("admin_postcode").(string)

	//ContactBilling
	domainInfo.ContactBilling.Email = d.Get("billing_email").(string)
	domainInfo.ContactBilling.Fax = d.Get("billing_fax").(string)
	domainInfo.ContactBilling.JobTitle = d.Get("billing_jobtitle").(string)
	domainInfo.ContactBilling.NameLast = d.Get("billing_namelast").(string)
	domainInfo.ContactBilling.NameFirst = d.Get("billing_namefirst").(string)
	domainInfo.ContactBilling.NameMiddle = d.Get("billing_namemiddle").(string)
	domainInfo.ContactBilling.Organization = d.Get("billing_organization").(string)
	domainInfo.ContactBilling.Phone = d.Get("billing_phone").(string)
	domainInfo.ContactBilling.AddressMailing.Address1 = d.Get("billing_address").(string)
	domainInfo.ContactBilling.AddressMailing.City = d.Get("billing_city").(string)
	domainInfo.ContactBilling.AddressMailing.Country = d.Get("billing_country").(string)
	domainInfo.ContactBilling.AddressMailing.State = d.Get("billing_state").(string)
	domainInfo.ContactBilling.AddressMailing.PostalCode = d.Get("billing_postcode").(string)

	//reg
	domainInfo.ContactRegistrant.Email = d.Get("reg_email").(string)
	domainInfo.ContactRegistrant.Fax = d.Get("reg_fax").(string)
	domainInfo.ContactRegistrant.JobTitle = d.Get("reg_jobtitle").(string)
	domainInfo.ContactRegistrant.NameLast = d.Get("reg_namelast").(string)
	domainInfo.ContactRegistrant.NameFirst = d.Get("reg_namefirst").(string)
	domainInfo.ContactRegistrant.NameMiddle = d.Get("reg_namemiddle").(string)
	domainInfo.ContactRegistrant.Organization = d.Get("reg_organization").(string)
	domainInfo.ContactRegistrant.Phone = d.Get("reg_phone").(string)
	domainInfo.ContactRegistrant.AddressMailing.Address1 = d.Get("reg_address").(string)
	domainInfo.ContactRegistrant.AddressMailing.City = d.Get("reg_city").(string)
	domainInfo.ContactRegistrant.AddressMailing.Country = d.Get("reg_country").(string)
	domainInfo.ContactRegistrant.AddressMailing.State = d.Get("reg_state").(string)
	domainInfo.ContactRegistrant.AddressMailing.PostalCode = d.Get("reg_postcode").(string)

	//tech
	domainInfo.ContactTech.Email = d.Get("tech_email").(string)
	domainInfo.ContactTech.Fax = d.Get("tech_fax").(string)
	domainInfo.ContactTech.JobTitle = d.Get("tech_jobtitle").(string)
	domainInfo.ContactTech.NameLast = d.Get("tech_namelast").(string)
	domainInfo.ContactTech.NameFirst = d.Get("tech_namefirst").(string)
	domainInfo.ContactTech.NameMiddle = d.Get("tech_namemiddle").(string)
	domainInfo.ContactTech.Organization = d.Get("tech_organization").(string)
	domainInfo.ContactTech.Phone = d.Get("tech_phone").(string)
	domainInfo.ContactTech.AddressMailing.Address1 = d.Get("tech_address").(string)
	domainInfo.ContactTech.AddressMailing.City = d.Get("tech_city").(string)
	domainInfo.ContactTech.AddressMailing.Country = d.Get("tech_country").(string)
	domainInfo.ContactTech.AddressMailing.State = d.Get("tech_state").(string)
	domainInfo.ContactTech.AddressMailing.PostalCode = d.Get("tech_postcode").(string)

	log.Println(json.Marshal(domainInfo))

	return config.Client()
}
