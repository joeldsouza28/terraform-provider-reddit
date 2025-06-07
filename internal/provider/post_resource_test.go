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

func TestRedditPost_Basic(t *testing.T) {
	postTitle := "Terraform Test Post"
	initialText := "Hello from Terraform!"
	updatedText := "Updated via Terraform!"
	subreddit := "test" // Replace with a subreddit where your bot has permission

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testPostPreCheck(t) },
		ProtoV6ProviderFactories: testPostProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRedditPostConfig(postTitle, initialText, subreddit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("reddit_post.test", "title", postTitle),
					resource.TestCheckResourceAttr("reddit_post.test", "text", initialText),
					resource.TestCheckResourceAttr("reddit_post.test", "subreddit", subreddit),
				),
			},
			{
				Config: testAccRedditPostConfig(postTitle, updatedText, subreddit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("reddit_post.test", "text", updatedText),
				),
			},
		},
	})
}

func testPostPreCheck(t *testing.T) {
	required := []string{"REDDIT_CLIENT_ID", "REDDIT_CLIENT_SECRET", "REDDIT_USERNAME", "REDDIT_PASSWORD"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			t.Fatalf("Environment variable %s must be set for acceptance tests", v)
		}
	}
}

func testPostProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"reddit": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

func testAccRedditPostConfig(title, text, subreddit string) string {
	return fmt.Sprintf(`
provider "reddit" {
  client_id     = "%s"
  client_secret = "%s"
  username      = "%s"
  password      = "%s"
}

resource "reddit_post" "test" {
  subreddit = "%s"
  title     = "%s"
  text      = "%s"
}
`,
		os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_CLIENT_SECRET"),
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		subreddit, title, text)
}
