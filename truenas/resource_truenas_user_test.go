package truenas

import (
	"context"
	"crypto/tls"
	"fmt"
	api "github.com/dariusbakunas/truenas-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"net/http"
	"strconv"
	"testing"
)

func TestAccResourceTruenasUser_basic(t *testing.T) {
	var user api.User
	suffix := acctest.RandStringFromCharSet(3, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("%s-%s", testResourcePrefix, suffix)
	resourceName := "truenas_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceTruenasUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckResourceTruenasUserConfig(name, fmt.Sprintf("%s@example.com", name), fmt.Sprintf("%s %s", name, name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "create_group", "false"),
					resource.TestCheckResourceAttr(resourceName, "email", fmt.Sprintf("%s@example.com", name)),
					resource.TestCheckResourceAttr(resourceName, "full_name", fmt.Sprintf("%s %s", name, name)),
					resource.TestCheckResourceAttr(resourceName, "gid", "1"),
					resource.TestCheckResourceAttr(resourceName, "group_ids", "[]"),
					//resource.TestCheckResourceAttr(resourceName, "quota_bytes", "2147483648"),
					//resource.TestCheckResourceAttr(resourceName, "quota_critical", "90"),
					//resource.TestCheckResourceAttr(resourceName, "quota_warning", "70"),
					//resource.TestCheckResourceAttr(resourceName, "ref_quota_bytes", "1073741824"),
					//resource.TestCheckResourceAttr(resourceName, "ref_quota_critical", "90"),
					//resource.TestCheckResourceAttr(resourceName, "ref_quota_warning", "70"),
					//resource.TestCheckResourceAttr(resourceName, "deduplication", "off"),
					//resource.TestCheckResourceAttr(resourceName, "exec", "on"),
					//resource.TestCheckResourceAttr(resourceName, "snap_dir", "hidden"),
					//resource.TestCheckResourceAttr(resourceName, "readonly", "off"),
					//resource.TestCheckResourceAttr(resourceName, "record_size", "256K"),
					//resource.TestCheckResourceAttr(resourceName, "case_sensitivity", "mixed"),
					testAccCheckTruenasUserResourceExists(resourceName, &user),
				),
			},
		},
	})
}

func testAccCheckResourceTruenasUserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.APIClient)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// loop through the resources in state, verifying each widget
	// is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "truenas_user" {
			continue
		}

		// Try to find the user
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		_, http, err := client.UserApi.GetUser(context.Background(), int32(id)).Execute()

		if err == nil {
			return fmt.Errorf("dataset (%s) still exists", rs.Primary.ID)
		}

		// check if error is in fact 404 (not found)
		if http.StatusCode != 404 {
			return fmt.Errorf("Error occured while checking for absence of dataset (%s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckResourceTruenasUserConfig(name string, email string, fullName string) string {
	return fmt.Sprintf(`
	resource "truenas_user" "test" {
		name = "%s"
		create_group = false
		email = "%s"
		full_name = "%s"
		gid = 1
		group_ids = []
		home_directory = "/nonexistent"
		home_mode = "0755"
		locked = false
		microsoft_account = true
		password = "changeme"
		password_disabled = false
		shell = "/bin/sh"
		smb = true
		ssh_public_key = null
		sudo = true
		sudo_commands = ["/bin/ls"]
		sudo_no_password = true
		uid = 1000
	}
	`, name, email, fullName)
}

func testAccCheckTruenasUserResourceExists(n string, user *api.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("user resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no user ID is set")
		}

		client := testAccProvider.Meta().(*api.APIClient)
		id, err := strconv.Atoi(rs.Primary.ID)
		resp, _, err := client.UserApi.GetUser(context.Background(), int32(id)).Execute()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(resp.Id)) != rs.Primary.ID {
			return fmt.Errorf("user not found")
		}

		user = resp
		return nil
	}
}
