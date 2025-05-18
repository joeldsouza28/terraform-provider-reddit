package provider

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePost() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreatePost,
		Read:   resourceReadPost,
		Delete: resourceDeletePost,
		Update: resourceUpdatePost,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // title can't be updated on Reddit
			},
			"text": {
				Type:     schema.TypeString,
				Optional: true, // this allows update!
			},
			"subreddit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // changing subreddit creates a new post
			},
			"post_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flair": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"nsfw": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCreatePost(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]string)
	client_id := config["client_id"]
	client_secret := config["client_secret"]
	username := config["username"]
	password := config["password"]
	token, err := GetAccessToken(client_id, client_secret, username, password)
	subreddit := d.Get("subreddit").(string)
	title := d.Get("title").(string)
	text := d.Get("text").(string)
	flair := d.Get("flair").(string)
	nsfw := d.Get("nsfw").(bool)

	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err
	}

	postID, err := SubmitPost(token, subreddit, title, text, flair, nsfw)
	if err != nil {
		log.Printf("[ERROR] Something went wrong while creating post %s", err)
		return err
	}

	d.SetId(postID)
	d.Set("post_id", postID)
	log.Printf("[INFO] Reddit post created with ID: %s", postID)

	return nil
}

func resourceReadPost(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceDeletePost(d *schema.ResourceData, meta interface{}) error {

	config := meta.(map[string]string)
	client_id := config["client_id"]
	client_secret := config["client_secret"]
	username := config["username"]
	password := config["password"]
	postID := d.Id()
	fullname := postID
	if !strings.HasPrefix(postID, "t3_") {
		fullname = "t3_" + postID
	}
	token, err := GetAccessToken(client_id, client_secret, username, password)
	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err
	}

	if err := DeletePost(token, fullname); err != nil {
		log.Printf("[ERROR] Something went wrong while deleting post %s", err)
		return err
	}

	return nil
}

func resourceUpdatePost(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]string)
	client_id := config["client_id"]
	client_secret := config["client_secret"]
	username := config["username"]
	password := config["password"]
	postID := d.Id()
	newText := d.Get("text").(string)
	token, err := GetAccessToken(client_id, client_secret, username, password)

	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err
	}

	if err := UpdatePostText(token, postID, newText); err != nil {
		log.Printf("[ERROR] Something went wrong while updating post %s", err)
		return err
	}

	return nil
}
