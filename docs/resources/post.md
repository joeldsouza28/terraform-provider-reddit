---
page_title: "reddit_post Resource - terraform-provider-reddit"
subcategory: ""
description: |-
  
---
# Resource `reddit_post`

## Example usage


```terraform 

resource "reddit_post" "my_post" {
   subreddit = var.subreddit
   title     = var.title
   text      = var.text
}

variable "subreddit" {
    type = string
}

variable "title" {
    type = string
}

variable "text" {
    type = string
}


```

## Schema

### Required
- **subreddit** The target subreddit where you want to post your content
- **title** The title of your reddit post

### Optional
- **text** The text of your reddit post
- **flair** Optional key. May be required for certain subreddits. If not provided for certain subreddits and error may be thrown 
- **nsfw** Optional flag with value true or false to mark your post as nsfw
