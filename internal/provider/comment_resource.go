package provider

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComment() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateComment,
		Delete: resourceDeleteComment,
		Update: resourceUpdateComment,
		Read:   resourceReadComment,
		Schema: map[string]*schema.Schema{
			"post_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Optional: false,
			},
			"comment": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"comment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateComment(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]string)
	client_id := config["client_id"]
	client_secret := config["client_secret"]
	username := config["username"]
	password := config["password"]
	text := d.Get("comment").(string)
	post_id := d.Get("post_id").(string)
	token, err := GetAccessToken(client_id, client_secret, username, password)
	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err
	}
	fullname := post_id
	if !strings.HasPrefix(post_id, "t3_") {
		fullname = "t3_" + post_id
	}
	comment_id, err := AddComment(token, fullname, text)
	if err != nil {
		log.Printf("[ERROR] Something went wrong while adding comment %s", err)
	}
	d.SetId(comment_id)
	d.Set("comment_id", comment_id)
	return nil
}

func resourceDeleteComment(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]string)
	client_id := config["client_id"]
	client_secret := config["client_secret"]
	username := config["username"]
	password := config["password"]
	comment_id := d.Get("comment_id").(string)
	token, err := GetAccessToken(client_id, client_secret, username, password)

	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err
	}

	fullname := comment_id
	if !strings.HasPrefix(comment_id, "t1_") {
		fullname = "t1_" + comment_id
	}

	if err := DeletePost(token, fullname); err != nil {
		log.Printf("[ERROR] Something went wrong while deleting comment %s", err)
	}
	return nil
}
func resourceReadComment(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceUpdateComment(d *schema.ResourceData, meta interface{}) error {
	config := meta.(map[string]string)
	client_id := config["client_id"]
	client_secret := config["client_secret"]
	username := config["username"]
	password := config["password"]
	token, err := GetAccessToken(client_id, client_secret, username, password)
	newText := d.Get("text").(string)

	comment_id := d.Get("comment_id").(string)

	if err != nil {
		log.Printf("[ERROR] Something went wrong while getting access token %s", err)
		return err
	}

	if err := UpdatePostText(token, comment_id, newText); err != nil {
		log.Printf("[ERROR] Something went wrong while updating comment %s", err)
		return err
	}

	return nil
}
