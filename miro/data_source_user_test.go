package miro

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.miro_user.user2", "team_id", ""),
					resource.TestCheckResourceAttr(
						"data.miro_user.user2", "id", ""),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig() string {
	return fmt.Sprintf(`	  
	resource "miro_user" "user2" {
		email        = ""
		team_id		 = ""
	  }
	data "miro_user" "user1" {
		email        = ""
		team_id		 = ""
	}
	`)
}