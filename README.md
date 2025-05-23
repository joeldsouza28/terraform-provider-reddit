
# <img src="https://raw.githubusercontent.com/joeldsouza28/terraform-provider-reddit/61657acd37cc130a849773ffe361f311a4be4c8b/reddit.png" width=30 height=30/> Reddit Terraform Provider

![Terraform](https://img.shields.io/badge/Terraform-Provider-blue?logo=terraform)
![Go](https://img.shields.io/badge/Built%20with-Go-00ADD8?logo=go)

The `joeldsouza28/reddit` Terraform provider allows you to write posts to public subreddits directly from your Terraform configuration.

---

## 🚀 Features

- ✅ Write post on subreddits
- ✅ Write comments to those post
- ✅ Get data from a post using post id
---

## 📦 Installation

Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    reddit = {
      source  = "joeldsouza28/reddit"
      version = "~> 0.29.0"
    }
  }
}
```

Then run:

```bash
terraform init
```

---

## 🔧 Usage

### Example: Writing post and fetching post content from Subreddit

```hcl

resource "reddit_post" "example_post" {
  subreddit = "kubernetes"
  title     = "Hi there everyone"
  text      = "Text some"
}


data "reddit_post" "example_data_post" {
  post_id = reddit_post.example_post.id
}


output "title" {
  value = data.example_data_post.title
}

output "text" {
  value = data.example_data_post.text
}

output "subreddit" {
  value = data.example_data_post.subreddit
}
```

---

## 🛠️ Inputs

| Name       | Type     | Description                         | Required |
|------------|----------|-----------------------------------|----------|
| `subreddit`| `string` | The name of the subreddit         | ✅ Yes   |
| `title`    | `string` | The title of the post             | ✅ Yes   |
| `text`     | `string` | The text of the post              | :negative_squared_cross_mark: No   |


---

## 📤 Outputs

The `reddit_post` data source returns data of the posts with the following attributes:

### `posts` (list of objects):

| Attribute    | Type     | Description                      |
|--------------|----------|---------------------------------|
| `title`      | `string` | Title of the reddit post         |
| `text`       | `string` | Text of the reddit post                  |
| `subreddit`  | `string` | Subreddit of the reddit post     |

---


## 📄 License

MIT License © 2025 Joel D’Souza

---

## 🙌 Contributing

PRs welcome! Feel free to open issues or improvements.
