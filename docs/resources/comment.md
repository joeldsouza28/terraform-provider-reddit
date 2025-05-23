---
page_title: "reddit_post Resource - terraform-provider-reddit"
subcategory: ""
description: |-
  
---

# Resource `reddit_comment`

## Example usage


```terraform 

resource "reddit_post" "my_comment" {
   post_id = var.post_id
   comment     = var.comment
}

variable "post_id" {
    type = string
}

variable "comment" {
    type = string
}



```

## Schema

### Required
- **post_id** The target post id under which you want to add a comment
- **comment** The comment text of your reddit comment
