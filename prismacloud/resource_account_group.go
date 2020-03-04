package prismacloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-prismacloud/prismacloud_client"
	"sort"
)

func resourceAccountGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountGroupCreate,
		Read:   resourceAccountGroupRead,
		Update: resourceAccountGroupUpdate,
		Delete: resourceAccountGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAccountGroupCreate(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	accountIds := []string{}

	req := prismacloud_client.AccountGroupRequest{name, description, accountIds}

	if attr, ok := d.GetOk("account_ids"); ok {
		req.AccountIds = attr.([]string)
	}

	err := prismaClient.Post("/cloud/group", req, nil)

	if err != nil {
		return fmt.Errorf("Failure creating account group %s.  %+v", name, err)
	}

	//Once we create the group we need to query all the groups to get its ID
	retval := []prismacloud_client.AccountGroup{}
	err = prismaClient.Get("/cloud/group", &retval)
	if err != nil {
		return fmt.Errorf("Failure reading account group id.  %+v", err)
	}
	sort.Sort(prismacloud_client.AccountGroupByName(retval))
	idx := sort.Search(len(retval), func(i int) bool { return retval[i].Name >= name })

	d.SetId(retval[idx].Id)

	return resourceAccountGroupRead(d, m)
}

func resourceAccountGroupRead(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)

	accountGroupId := d.Id()

	accountGroup := new(prismacloud_client.AccountGroup)
	err := prismaClient.Get(fmt.Sprintf("/cloud/group/%s", accountGroupId), accountGroup)
	if err != nil {
		return fmt.Errorf("Failure reading account group %s, %+v", accountGroupId, err)
	}

	d.Set("name", accountGroup.Name)
	d.Set("description", accountGroup.Description)
	d.Set("account_ids", accountGroup.AccountIds)

	return nil
}

func resourceAccountGroupUpdate(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)
	d.Partial(true)

	requiresChange := false
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	accountIds := expandStringList(d.Get("account_ids").([]interface{}))

	req := prismacloud_client.AccountGroupRequest{name, description, accountIds}

	if d.HasChange("name") {
		d.SetPartial("name")
		requiresChange = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		requiresChange = true
	}

	if requiresChange {
		accountGroupId := d.Id()
		err := prismaClient.Put(fmt.Sprintf("/cloud/group/%s", accountGroupId), req)
		if err != nil {
			return err
		}
	}
	d.Partial(false)
	return resourceAccountRead(d, m)
}

func resourceAccountGroupDelete(d *schema.ResourceData, m interface{}) error {
	prismaClient := m.(*prismacloud_client.PrismaCloudClient)
	accountGroupId := d.Id()
	err := prismaClient.Delete(fmt.Sprintf("/cloud/group/%s", accountGroupId))
	if err != nil {
		return err
	}
	return nil
}
