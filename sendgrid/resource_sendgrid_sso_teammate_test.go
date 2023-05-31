package sendgrid_test

import (
	"fmt"
	"testing"

	sendgrid "github.com/SpotOnInc/terraform-provider-sendgrid/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSendgridSsoTeammateBasic(t *testing.T) {
	email := "terraform-user-" + acctest.RandString(10) + "@example.org"
	firstName := "terraform"
	lastName := "user"
	isAdmin := false
	persona := "observer"
	scopes := []string{"mail.send", "alerts.read"}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSendgridSSOTeammateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSendgridSSOTeammateConfigBasic(
					email,
					firstName,
					lastName,
					isAdmin,
					persona,
					scopes,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSendgridSSOTeammateExists("sendgrid_sso_teammate.this"),
				),
			},
		},
	})
}

func testAccCheckSendgridSSOTeammateDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*sendgrid.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sendgrid_sso_teammate" {
			continue
		}

		email := rs.Primary.ID

		if _, request := c.DeleteSSOTeamMate(email); request.Err != nil {
			return request.Err
		}
	}

	return nil
}

func testAccCheckSendgridSSOTeammateConfigBasic(
	email string,
	firstName string,
	lastName string,
	isAdmin bool,
	persona string,
	scopes []string,
) string {
	return fmt.Sprintf(`
resource "sendgrid_sso_teammate" "this" {
	email      = "%s"
	first_name = "%s"
	last_name  = "%s"
	is_admin   = %t
	persona    = "%s"
	scopes     = %q
}
`, email, firstName, lastName, isAdmin, persona, scopes)
}

func testAccCheckSendgridSSOTeammateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SSO teammate ID is set")
		}

		c := testAccProvider.Meta().(*sendgrid.Client)
		teammate, request := c.ReadSSOTeamMate(rs.Primary.ID)

		if request.Err != nil {
			return request.Err
		}

		if teammate.Email != rs.Primary.ID {
			return fmt.Errorf("SSO teammate not found")
		}

		return nil
	}
}
