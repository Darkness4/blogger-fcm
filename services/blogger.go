package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Darkness4/blogger-fcm/models"
)

// HTTPClient interface
type HTTPClient interface {
	Get(string) (*http.Response, error)
}

// Blogger service
type Blogger struct {
	client  HTTPClient
	baseURL url.URL
	BlogID  string
	Key     string
}

// NewBlogger instanciates a Blogger
func NewBlogger(client HTTPClient) (*Blogger, error) {
	// Base URL
	url, err := url.Parse("https://www.googleapis.com/blogger/v3/blogs/")
	if err != nil {
		return nil, err
	}

	blogger := &Blogger{
		client:  client,
		baseURL: *url,
		BlogID:  "BLOG_ID",
		Key:     "BLOGGER_API_KEY",
	}

	// Check Env Variables
	blogID := os.Getenv("BLOG_ID")
	if blogID != "" {
		blogger.BlogID = blogID
	} else {
		text := `
The env variable BLOG_ID does not exist !
	
Please see https://developers.google.com/blogger/docs/3.0/using#RetrievingABlog
And add a env variable BLOG_ID.`
		log.Println(text)
	}

	key := os.Getenv("BLOGGER_API_KEY")
	if key != "" {
		blogger.Key = key
	} else {
		text := `
The env variable with BLOGGER_API_KEY does not exist !

Please see https://developers.google.com/blogger/docs/3.0/using#APIKey
And add a env variable BLOGGER_API_KEY.`
		log.Println(text)
	}

	return blogger, nil
}

// GetBlog fetch blog from Blogger
func (b *Blogger) GetBlog() (*models.Blog, error) {
	u, err := url.Parse(b.BlogID + "?key=" + b.Key)
	if err != nil {
		return nil, err
	}
	endpoint := b.baseURL.ResolveReference(u)
	log.Println("GET " + endpoint.String())
	resp, err := b.client.Get(endpoint.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP Error : " + resp.Status)
	}
	defer resp.Body.Close()

	blog := new(models.Blog)
	err = json.NewDecoder(resp.Body).Decode(blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

// GetLatestPost fetch latest post from Blogger
func (b *Blogger) GetLatestPost() (*models.Post, error) {
	u, err := url.Parse(b.BlogID + "/posts" + "?key=" + b.Key)
	if err != nil {
		return nil, err
	}
	endpoint := b.baseURL.ResolveReference(u)
	log.Println("GET " + endpoint.String())
	resp, err := b.client.Get(endpoint.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP Error : " + resp.Status)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := struct {
		Items []models.Post `json:"items"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data.Items[0], nil
}
