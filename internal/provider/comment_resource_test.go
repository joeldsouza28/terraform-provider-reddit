package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccRedditCommentConfigBasic() string {
	return `
provider "reddit" {
  client_id     = "YOUR_CLIENT_ID"
  client_secret = "YOUR_CLIENT_SECRET"
  username      = "YOUR_USERNAME"
  password      = "YOUR_PASSWORD"
}

resource "reddit_comment" "test" {
  post_id = "t3_abcdef"
  text    = "This is a test comment"
}
`
}

var testAccProvider = Provider()

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"reddit": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}

func testAccPreCheck(t *testing.T) {
	// Optional: check env vars or credentials here
}

func TestAccRedditComment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRedditCommentConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("reddit_comment.test", "post_id", "t3_abcdef"), // Use real post ID
					resource.TestCheckResourceAttr("reddit_comment.test", "text", "This is a test comment"),
				),
			},
		},
	})
}
