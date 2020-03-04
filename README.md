# terraform-provider-prismcloud

Prisma Cloud is a tool owned by Palo Alto networks that provides insights into your cloud compliance and security 
posture.  This terraform provider leverages the [api](https://api.docs.prismacloud.io/reference).  Because Prismacloud 
does not have a go client, we choose to maintain our own.

## Setup

Installation involves putting the binary for your system in the `~/.terraform.d/plugins/` directory.  See the
[documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for more details. 

In addition to installing the executable, this provider looks in `~/.prismacloud/credentials` for credentials to access 
the rest API.  That file should look like the below: 

```json 
{
  "AccessKeyId": "fill-me-in",
  "SecretAccessKey": "fill-me-in"
}

```

## Usage

Configure this provider like any other.  Make sure to configure the `base_url` to the correct url.  A list of urls
can be found [here](https://api.docs.prismacloud.io/reference#try-the-apis).

```hcl-terraform
provider "prismacloud" {
  base_url = "https://api.prismacloud.io"
}
```

### Resources

#### Account Group
Provides a Prisma Cloud account group.

##### Example
```hcl-terraform
resource "prismacloud_account_group" "tester" {
  name = "tf-test"
  description = "My example account group!"
}
```

##### Argument Reference
The following arguments are supported by this resource:
* `name` - (Required) The name of the account group
* `description` - (Optional) Short description of what the account group is for.

##### Attribute References
In addition to all of the arguments, the following are also exported:
* `account_ids` - Ids of accounts that are attached to this account group


#### Account
Provides a Prisma Cloud account that represents a cloud account.

##### Example

An example AWS account.
```hcl-terraform
resource "prismacloud_account" "tester" {
  name = "tf-test-2"
  cloud_type = "aws"
  account_id = "111111111111"
  external_id = "some-secret"
  group_ids = [prismacloud_account_group.tester.id]
  role_arn = "arn:aws:iam::111111111111:role/testing"
}
```

##### Argument Reference
The following arguments are supported by this resource:
* `name` - (Required) The name to give the account
* `account_id` - (Required) The ID of your account given to you by you cloud provider
* `cloud_type` - (Required) One of AWS, GCP, Alibaba, or Azure
* `group_ids` - (Optional) The ids of the account groups that this account should belong to.

For AWS accounts:
* `external_id` - (Required) The external ID of the role prisma cloud asssumes
* `role_arn` - (Required) The ARN of the role that prisma cloud should assume to collect information from your account. 

##### Attribute References
All arguments are exported as attributes.  No additional attributes are exported In addition to all of the arguments.

## Development

The makefile will test and build the provider.  Here is a synopsis of the make targets

`make fmt` - uses `go fmt` to format the go files.
`make test` - uses `go test` to run all of the tests.
`make provider` - builds a linux and a mac binary and places them in the build directory.
`make clean` - removes any temporary or built files.
