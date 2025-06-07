package provider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/joeldsouza28/terraform-provider-reddit/internal/provider"
)

func TestAccRedditComment_basic(t *testing.T) {
	const (
		testPostID  = "1l5fc29"                // Replace with real or mocked post ID
		commentText = "This is a test comment" // Update if needed
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testCommentPreCheck(t) },
		ProtoV6ProviderFactories: testCommentProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRedditCommentConfig(testPostID, commentText),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("reddit_comment.test", "post_id", testPostID),
					resource.TestCheckResourceAttr("reddit_comment.test", "comment", commentText),
				),
			},
		},
	})
}

// Ensure required env vars are set before running tests
func testCommentPreCheck(t *testing.T) {
	required := []string{
		"REDDIT_CLIENT_ID",
		"REDDIT_CLIENT_SECRET",
		"REDDIT_USERNAME",
		"REDDIT_PASSWORD",
	}
	for _, v := range required {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for acceptance tests", v)
		}
	}
}

// Return a framework-compatible ProviderFactory
func testCommentProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"reddit": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

// Generates Terraform config for a reddit_comment test
func testAccRedditCommentConfig(postID, commentText string) string {
	return fmt.Sprintf(`
provider "reddit" {
  client_id     = "%s"
  client_secret = "%s"
  username      = "%s"
  password      = "%s"
}

resource "reddit_comment" "test" {
  post_id = "%s"
  comment = "%s"
}
`,
		os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_CLIENT_SECRET"),
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		postID,
		commentText,
	)
}

// Actual test
