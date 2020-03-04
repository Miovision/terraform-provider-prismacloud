package prismacloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-prismacloud/prismacloud_client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"prismacloud_account":       resourceAccount(),
			"prismacloud_account_group": resourceAccountGroup(),
		},
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	baseUrl := d.Get("base_url").(string)
	creds := prismacloud_client.LoadCredentials()
	c := prismacloud_client.MakePrismaCloudClient(baseUrl, creds.AccessKeyId, creds.SecretAccessKey)
	return c, nil
}
