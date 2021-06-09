package miro

import(
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("miro_user.user1", "email", "{email"),
					resource.TestCheckResourceAttr("miro_user.user1", "role", "member"),
				),
			},
		},
	})
}

func testAccCheckItemBasic() string {
	return fmt.Sprintf(`
		resource "miro_user" "user1" {
			email	= "{email}"
			role	= "member"
		}
	`)
}

func TestAccItem_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"miro_user.user1", "email", "{email}"),	
					resource.TestCheckResourceAttr(
						"miro_user.user1", "role", "member"),
				),
			},
			{
				Config: testAccCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"miro_user.user1", "email", "{email}"),	
					resource.TestCheckResourceAttr(
						"miro_user.user1", "role", "member"),
				),
			},
		},
	})
}

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
		resource "miro_user" "user1" {
			email	= "{email}"
			role	= "member"	
		}
	`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
		resource "miro_user" "user1" {
			email        = "{email}"
			role		 = "member"
		}
	`)
}
