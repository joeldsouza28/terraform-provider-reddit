package provider

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRedditPost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRedditPostRead,
		Schema: map[string]*schema.Schema{
			"post_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subreddit": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRedditPostRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]string)
	token, err := GetAccessToken(
		config["client_id"],
		config["client_secret"],
		config["username"],
		config["password"],
	)

	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err

	}

	postID := d.Get("post_id").(string)
	post, err := FetchPostByID(token, postID)
	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting post %s", err)
		return err

	}

	d.SetId(postID)
	d.Set("title", post.Title)
	d.Set("text", post.Text)
	d.Set("subreddit", post.Subreddit)
	return nil
}
