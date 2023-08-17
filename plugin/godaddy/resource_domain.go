package godaddy

import (
	"context"
	"fmt"
	"github.com/forease/gotld"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"terraform-provider-st-godaddy/api"
)

const (
	attrMode  = "mode"
	attrYears = "years"
)
const MODE_CREATE = "create"
const MODE_RENEW = "renew"

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		UpdateContext: resourceDomainUpdate,
		ReadContext:   resourceDomainRead,
		DeleteContext: resourceDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainImport,
		},

		Schema: map[string]*schema.Schema{
			attrDomain: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Purchased available domain name on your account",
			},
			attrMode: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "domain operation type, include create, renew",
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_MODE", "create"),
			},
			attrYears: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of years to register",
				Default:     "2",
			},
		},
	}
}

func resourceDomainImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	fmtlog(ctx, "[resourceRecordImport!]")
	if err := data.Set("domain", data.Id()); err != nil {
		return nil, err
	}
	if err := data.Set("mode", MODE_CREATE); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}

func resourceDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	fmtlog(ctx, "[resourceDomainCreate!]")
	client := meta.(*api.Client)

	domain := strings.ToLower(data.Get(attrDomain).(string))
	mode := strings.ToLower(data.Get(attrMode).(string))
	years := data.Get(attrYears).(int)

	switch mode {
	case MODE_CREATE:
		diags := createDomain(ctx, client, domain, years)
		if diags.HasError() {
			return diags
		}
	case MODE_RENEW:
		diags := renewDomain(ctx, client, domain, years)
		if diags.HasError() {
			return diags
		}
	default:
		return diag.Errorf("unsupported mode:%s, mode can only be create or renew", mode)
	}

	data.SetId(domain)

	return nil
}

func resourceDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	fmtlog(ctx, "[resourceDomainRead!]")
	client := meta.(*api.Client)

	domainName := strings.ToLower(data.Get(attrDomain).(string))

	_, err := client.GetDomain(domainName)

	if err == nil {
		_ = data.Set("domain", domainName)
	}

	return nil
}

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	fmtlog(ctx, "[resourceRecordUpdate!]")
	client := meta.(*api.Client)

	//we can do nothing on old name,year and mode
	oldDomainRaw, newDomainRaw := data.GetChange("domain")
	newDomain := newDomainRaw.(string)
	oldDomain := oldDomainRaw.(string)

	_, newYearRaw := data.GetChange("years")
	newYear := newYearRaw.(int)

	_, newModeRaw := data.GetChange("mode")
	newMode := newModeRaw.(string)

	switch newMode {
	case MODE_CREATE:
		//delete old domain first
		diags := deleteDomain(ctx, client, oldDomain)
		if diags.HasError() {
			return diags
		}

		//create new domain then
		diags = createDomain(ctx, client, newDomain, newYear)
		if diags.HasError() {
			return diags
		}

	case MODE_RENEW:
		//can't do anything about old domain
		diags := renewDomain(ctx, client, newDomain, newYear)
		if diags.HasError() {
			return diags
		}
	default:
		return diag.Errorf("unsupported mode:%s, mode can only be create or renew", newMode)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {

	fmtlog(ctx, "[resourceRecordDelete!]")
	client := meta.(*api.Client)
	domainName := strings.ToLower(data.Get(attrDomain).(string))

	diags := deleteDomain(ctx, client, domainName)
	if diags.HasError() {
		return diags
	}

	return nil
}

func createDomain(cxt context.Context, client *api.Client, domainName string, year int) diag.Diagnostics {

	var domains []string
	domains = append(domains, domainName)
	log.Println("domain", domainName, "do not exist, check whether it's available to purchase....")
	available, err := client.DomainAvailable(domains)
	if err != nil {
		return diag.FromErr(err)
	}
	if !available {
		return diag.Errorf("domain %s is not available,please try to another one", domainName)
	}

	//extract tld
	tld, _, err := gotld.GetTld(domainName)
	agreement, err := client.GetAgreement(tld.Tld, false)
	if err != nil {
		return diag.FromErr(err)
	}
	//construct agreement keys
	var agreementKeys []string
	for _, v := range agreement {
		agreementKeys = append(agreementKeys, v.AgreementKey)
	}

	err = client.Purchase(domainName, agreementKeys, _domainInfo)
	_domainInfo.Period = year
	if err != nil {
		fmtlog(cxt, "Creating [%s] failed!", domainName)
		return diag.FromErr(err)
	}
	fmtlog(cxt, "Creating [%s] success!", domainName)
	return nil
}

func renewDomain(cxt context.Context, client *api.Client, domainName string, year int) diag.Diagnostics {

	err := client.DomainRenew(domainName, year)
	if err != nil {
		fmtlog(cxt, "Renew [%s] failed!", domainName)
		return diag.FromErr(err)
	}
	fmtlog(cxt, "Renew [%s] success!", domainName)
	return nil
}

func deleteDomain(cxt context.Context, client *api.Client, domainName string) diag.Diagnostics {

	err := client.DomainCancel(domainName)
	if err != nil {
		fmtlog(cxt, "Delete [%s] failed!", domainName)
		return diag.FromErr(err)
	}
	fmtlog(cxt, "Delete [%s] success!", domainName)
	return nil
}

func fmtlog(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a)
	tflog.Info(ctx, msg)
}
