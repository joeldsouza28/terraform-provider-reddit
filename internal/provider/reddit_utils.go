package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type redditClient struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

func (r *redditClient) GetToken() (string, error) {
	return GetAccessToken(r.ClientID, r.ClientSecret, r.Username, r.Password)
}

func ExtractPostIDFromHTMLJSON(rawJSON []byte) (string, error) {
	var resp map[string]interface{}
	if err := json.Unmarshal(rawJSON, &resp); err != nil {
		return "", err
	}

	jqueryArr, ok := resp["jquery"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no jquery field found")
	}

	for _, item := range jqueryArr {
		arr, ok := item.([]interface{})
		if !ok || len(arr) < 4 {
			continue
		}

		// The 3rd element is the method name (string)
		method, ok := arr[2].(string)
		if !ok || method != "call" {
			continue
		}

		// The 4th element is the args array
		args, ok := arr[3].([]interface{})
		if !ok || len(args) == 0 {
			continue
		}

		// Look for URL containing "/comments/"
		for _, arg := range args {
			strArg, ok := arg.(string)
			if !ok {
				continue
			}
			if strings.Contains(strArg, "/comments/") {
				// Extract the post ID from the URL
				// URL format: https://www.reddit.com/r/{subreddit}/comments/{postID}/...
				parts := strings.Split(strArg, "/")
				for i, part := range parts {
					if part == "comments" && i+1 < len(parts) {
						return parts[i+1], nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("post URL not found in jquery")
}

func GetAccessToken(clientID, clientSecret, username, password string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequest("POST", "https://oauth.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", "go-reddit-script/0.1 by "+username)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("failed to get access token: %s", body)
}

func SubmitPost(accessToken, subreddit, title, text, flair string, nsfw bool) (string, error) {
	data := url.Values{}
	data.Set("sr", subreddit)
	data.Set("kind", "self")
	data.Set("title", title)
	data.Set("text", text)
	log.Println(flair)
	if flair != "" {
		data.Set("flair_text", flair)
		flair_id, err := GetFlairID(accessToken, subreddit, flair)
		if err != nil {
			return "", err
		}
		data.Set("flair_id", flair_id)
	}
	data.Set("api_type", "json")
	data.Set("nsfw", strconv.FormatBool(nsfw))

	req, err := http.NewRequest("POST", "https://oauth.reddit.com/api/submit", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("User-Agent", "go-reddit-script/0.1 by Sensitive-Cake-1569")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Content-Type:", resp.Header.Get("Content-Type"))

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	var response struct {
		Json struct {
			Data struct {
				URL  string `json:"url"`
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"data"`
			Errors [][]interface{} `json:"errors"` // if any
		} `json:"json"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(response.Json.Errors) > 0 {
		return "", fmt.Errorf("Reddit API error: %+v", response.Json.Errors)
	}

	// Parse relevant fields
	return response.Json.Data.ID, nil
}

func DeletePost(accessToken string, postFullname string) error {
	data := url.Values{}
	data.Set("id", postFullname)

	req, err := http.NewRequest("POST", "https://oauth.reddit.com/api/del", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("User-Agent", "go-reddit-script/0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Delete response:", string(body))

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", body)
	}

	return nil
}

func UpdatePostText(accessToken, postID, newText string) error {
	// Ensure fullname format: "t3_" + postID
	if !strings.HasPrefix(postID, "t3_") {
		postID = "t3_" + postID
	}

	data := url.Values{}
	data.Set("thing_id", postID)
	data.Set("text", newText)

	req, err := http.NewRequest("POST", "https://oauth.reddit.com/api/editusertext", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("User-Agent", "go-reddit-script/0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update post: %s", string(body))
	}

	return nil
}

func GetFlairID(accessToken, subreddit, flairText string) (string, error) {
	url := fmt.Sprintf("https://oauth.reddit.com/r/%s/api/link_flair_v2", subreddit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("User-Agent", "go-reddit-script/0.1 by yourusername")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch flairs: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var flairs []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(body, &flairs); err != nil {
		return "", err
	}

	for _, flair := range flairs {
		if flair.Text == flairText {
			log.Printf("Flair Id: %s", flair.ID)
			return flair.ID, nil
		}
	}

	return "", fmt.Errorf("flair with text '%s' not found", flairText)
}

func AddComment(accessToken, parentFullname, text string) (string, error) {
	endpoint := "https://oauth.reddit.com/api/comment"

	data := url.Values{}
	data.Set("api_type", "json")
	data.Set("thing_id", parentFullname) // e.g., t3_abc123 (post) or t1_xyz789 (comment)
	data.Set("text", text)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("User-Agent", "go-reddit-script/0.1 by yourusername")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Check for API errors

	if jsonData, ok := result["json"].(map[string]interface{}); ok {
		if data, ok := jsonData["data"].(map[string]interface{}); ok {
			if things, ok := data["things"].([]interface{}); ok && len(things) > 0 {
				if thing, ok := things[0].(map[string]interface{}); ok {
					if data, ok := thing["data"].(map[string]interface{}); ok {
						if id, ok := data["id"].(string); ok {
							return "t1_" + id, nil
						}
					}
				}
			}
		}
	}

	return "", nil

}

type RedditPost struct {
	Title     string
	Text      string
	Subreddit string
}

// FetchPostByID fetches a Reddit post by its ID using Reddit's JSON API
func FetchPostByID(token, postID string) (*RedditPost, error) {
	// Ensure postID is in the correct format
	postID, err := strconv.Unquote(postID)
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(postID, "t3_") {
		postID = "t3_" + postID
	}
	url := fmt.Sprintf("https://oauth.reddit.com/api/info.json?id=%s", postID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("User-Agent", "terraform-provider-reddit/0.1 by your-username")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Data struct {
			Children []struct {
				Data struct {
					Title     string `json:"title"`
					SelfText  string `json:"selftext"`
					Subreddit string `json:"subreddit"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Data.Children) == 0 {
		return nil, fmt.Errorf("no post found with ID %s", postID)
	}

	postData := data.Data.Children[0].Data

	return &RedditPost{
		Title:     postData.Title,
		Text:      postData.SelfText,
		Subreddit: postData.Subreddit,
	}, nil
}
