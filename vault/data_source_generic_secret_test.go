package vault

import (
	"fmt"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestDataSourceGenericSecret(t *testing.T) {
	r.Test(t, r.TestCase{
		Providers: testProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []r.TestStep{
			{
				Config: testDataSourceGenericSecret_config,
				Check:  testDataSourceGenericSecret_check,
			},
		},
	})
}

func TestV2Secret(t *testing.T) {
	r.Test(t, r.TestCase{
		Providers: testProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []r.TestStep{
			{
				Config: testv2DataSourceGenericSecret_config,
				Check:  testDataSourceGenericSecret_check,
			},
			{
				Config: testv2DataSourceGenericSecretUpdated_config,
				Check:  testDataSourceGenericSecret_check,
			},
		},
	})
}

var testv2DataSourceGenericSecret_config = `

resource "vault_generic_secret" "test" {
    path = "secret/foo"
    data_json = <<EOT
{
    "zip": "zap"
}
EOT
}

data "vault_generic_secret" "test" {
    path = "${vault_generic_secret.test.path}"
		version = -1
}

`
var testv2DataSourceGenericSecretUpdated_config = `

resource "vault_generic_secret" "test" {
    path = "secret/foo"
    data_json = <<EOT
{
    "zip": "kablamo"
}
EOT
}

data "vault_generic_secret" "test" {
    path = "${vault_generic_secret.test.path}"
		version = 1
}

`

var testDataSourceGenericSecret_config = `

resource "vault_mount" "v1" {
	  path = "secretsv1"
	  type = "kv"
	  options = {
		  version = "1"
	  }
}

resource "vault_generic_secret" "test" {
    path = "${vault_mount.v1.path}/foo"
    data_json = <<EOT
{
    "zip": "zap"
}
EOT
}

data "vault_generic_secret" "test" {
    path = "${vault_generic_secret.test.path}"
}

`

func testDataSourceGenericSecret_check(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["data.vault_generic_secret.test"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state %v", s.Modules[0].Resources)
	}

	iState := resourceState.Primary
	if iState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	wantJson := `{"zip":"zap"}`
	if got, want := iState.Attributes["data_json"], wantJson; got != want {
		return fmt.Errorf("data_json contains %s; want %s", got, want)
	}

	if got, want := iState.Attributes["data.zip"], "zap"; got != want {
		return fmt.Errorf("data[\"zip\"] contains %s; want %s", got, want)
	}

	return nil
}
