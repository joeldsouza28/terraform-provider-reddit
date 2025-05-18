package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccPostProvider = Provider()

var testAccPostProviderFactories = map[string]func() (*schema.Provider, error){
	"reddit": func() (*schema.Provider, error) {
		return testAccPostProvider, nil
	},
}

func testAccPosPreCheck(t *testing.T) {
	requiredVars := []string{"REDDIT_CLIENT_ID", "REDDIT_CLIENT_SECRET", "REDDIT_USERNAME", "REDDIT_PASSWORD"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for acceptance tests", v)
		}
	}
}

func TestAccRedditPost_basic(t *testing.T) {
	postTitle := "Terraform Test Post"
	initialText := "Hello from Terraform!"
	updatedText := "Updated via Terraform!"
	subreddit := "test" // Or a subreddit where your bot has posting rights

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPosPreCheck(t) },
		ProviderFactories: testAccPostProviderFactories,
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
`, os.Getenv("REDDIT_CLIENT_ID"),
		os.Getenv("REDDIT_CLIENT_SECRET"),
		os.Getenv("REDDIT_USERNAME"),
		os.Getenv("REDDIT_PASSWORD"),
		subreddit, title, text)
}
