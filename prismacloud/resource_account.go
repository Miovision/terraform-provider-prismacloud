package prismacloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-prismacloud/prismacloud_client"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"external_id": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"group_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"role_arn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)

	accountId := d.Get("account_id").(string)
	enabled := true
	external_id := d.Get("external_id").(string)
	groupIds := expandStringList(d.Get("group_ids").([]interface{}))
	name := d.Get("name").(string)
	roleArn := d.Get("role_arn").(string)

	cloud_type := d.Get("cloud_type").(string)

	req := prismacloud_client.AccountRequest{accountId, enabled, external_id, groupIds, name, roleArn}
	err := prismaClient.Post(fmt.Sprintf("/cloud/%s", cloud_type), req, nil)

	if err != nil {
		return err
	}

	d.SetId(accountId)

	return resourceAccountRead(d, m)
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)

	accountId := d.Id()
	cloud_type := d.Get("cloud_type").(string)
	account := new(prismacloud_client.Account)
	err := prismaClient.Get(fmt.Sprintf("/cloud/%s/%s", cloud_type, accountId), account)

	if err != nil {
		return err
	}

	d.Set("account_id", account.AccountId)
	d.Set("group_ids", account.GroupIds)
	d.Set("name", account.Name)
	d.Set("cloud_type", account.CloudType)
	d.Set("external_id", account.ExternalId)
	d.Set("role_arn", account.RoleArn)
	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)

	accountId := d.Get("account_id").(string)
	enabled := true
	external_id := d.Get("external_id").(string)
	groupIds := expandStringList(d.Get("group_ids").([]interface{}))
	name := d.Get("name").(string)
	roleArn := d.Get("role_arn").(string)

	cloud_type := d.Get("cloud_type").(string)

	req := prismacloud_client.AccountRequest{accountId, enabled, external_id, groupIds, name, roleArn}

	err := prismaClient.Put(fmt.Sprintf("/cloud/%s/%s", cloud_type, accountId), req)

	if err != nil {
		return err
	}

	return resourceAccountRead(d, m)
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)
	accountId := d.Id()
	cloud_type := d.Get("cloud_type").(string)
	err := prismaClient.Delete(fmt.Sprintf("/cloud/%s/%s", cloud_type, accountId))
	if err != nil {
		return err
	}
	return nil
}
