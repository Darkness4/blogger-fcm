package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Darkness4/blogger-fcm/models"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

type mockHTTPClient struct {
	GetFunc func(u string) (*http.Response, error)
}

var (
	// GetFunc fetches the mock client's `Get` func
	GetFunc func(u string) (*http.Response, error)
)

func (m *mockHTTPClient) Get(u string) (*http.Response, error) {
	return GetFunc(u)
}

func TestNewBlogger(t *testing.T) {
	testCases := []struct {
		blogID        string
		bloggerAPIKey string
	}{
		{"", ""},
		{"123456789", ""},
		{"", "123456789abcdef"},
		{"123456789", "123456789abcdef"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("{blogID: %s, BloggerAPIKey: %s}", tc.blogID, tc.bloggerAPIKey), func(t *testing.T) {
			os.Unsetenv("BLOG_ID")
			os.Unsetenv("BLOGGER_API_KEY")
			if tc.blogID != "" {
				err := os.Setenv("BLOG_ID", tc.blogID)
				if err != nil {
					t.Fatal("Couldn't set env variable BLOG_ID")
				}
			}
			if tc.bloggerAPIKey != "" {
				err := os.Setenv("BLOGGER_API_KEY", tc.bloggerAPIKey)
				if err != nil {
					t.Fatal("Couldn't set env variable BLOGGER_API_KEY")
				}
			}

			mock := &mockHTTPClient{}
			result, err := NewBlogger(mock)
			if err != nil {
				t.Fatal("Error raised with NewBlogger")
			}
			if tc.blogID == "" && result.BlogID != "BLOG_ID" {
				t.Fatalf("Got: %s. Expected: %s", result.BlogID, "BLOG_ID")
			}
			if tc.blogID != "" && result.BlogID != tc.blogID {
				t.Fatalf("Got: %s. Expected: %s", result.BlogID, tc.blogID)
			}
			if tc.bloggerAPIKey == "" && result.Key != "BLOGGER_API_KEY" {
				t.Fatalf("Got: %s. Expected: %s", result.Key, "BLOGGER_API_KEY")
			}
			if tc.bloggerAPIKey != "" && result.Key != tc.bloggerAPIKey {
				t.Fatalf("Got: %s. Expected: %s", result.Key, tc.bloggerAPIKey)
			}
		})
	}
}

func TestGetBlog(t *testing.T) {
	body := helperLoadBytes(t, "blog.json")
	testCases := []struct {
		title         string
		blogID        string
		bloggerAPIKey string
		resp          *http.Response
		getErr        error
	}{
		{
			http.StatusText(http.StatusOK),
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
				ContentLength: int64(len(body)),
				Header:        make(http.Header, 0),
			},
			nil,
		},
		{
			http.StatusText(http.StatusNotFound),
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{
				Status:        http.StatusText(http.StatusNotFound),
				StatusCode:    http.StatusNotFound,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString(http.StatusText(http.StatusNotFound))),
				ContentLength: int64(len(http.StatusText(http.StatusNotFound))),
				Header:        make(http.Header, 0),
			},
			nil,
		},
		{
			"Get should throw an error",
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{},
			errors.New("No Wi-Fi"),
		},
		{
			"ResolveReference should throw an error",
			"$:$ù%§µ%M£",
			"tBloggerAPIKey",
			&http.Response{},
			nil,
		},
		{
			"json.Decode should throw an error",
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString("Not a json")),
				ContentLength: int64(len("Not a json")),
				Header:        make(http.Header, 0),
			},
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.title), func(t *testing.T) {
			os.Unsetenv("BLOG_ID")
			os.Unsetenv("BLOGGER_API_KEY")
			if tc.blogID != "" {
				err := os.Setenv("BLOG_ID", tc.blogID)
				if err != nil {
					t.Fatal("Couldn't set env variable BLOG_ID")
				}
			}
			if tc.bloggerAPIKey != "" {
				err := os.Setenv("BLOGGER_API_KEY", tc.bloggerAPIKey)
				if err != nil {
					t.Fatal("Couldn't set env variable BLOGGER_API_KEY")
				}
			}
			mock := &mockHTTPClient{}
			GetFunc = func(u string) (*http.Response, error) {
				if tc.getErr != nil {
					return nil, tc.getErr
				}
				return tc.resp, nil
			}
			blogger, err := NewBlogger(mock)
			if err != nil {
				t.Fatal("Error raised with NewBlogger")
			}
			resultData, resultErr := blogger.GetBlog()

			if tc.title == http.StatusText(http.StatusOK) {
				if resultData == nil {
					t.Fatal("Got: resp == nil. Expected: resp != nil")
				} else {
					expectedBytes := helperLoadBytes(t, "blog.json")
					expected := &models.Blog{}
					err = json.Unmarshal(expectedBytes, &expected)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(expected, resultData) {
						t.Fatalf("Got: %s. Expected: %s", resultData, expected)
					}
				}
				if resultErr != nil {
					t.Fatalf("Got: err = %s. Expected: err == nil", resultErr)
				}
			} else {
				if resultData != nil {
					t.Fatal("Got: resp != nil. Expected: resp == nil")
				}
				if resultErr == nil {
					t.Fatal("Got: err == nil. Expected: err != nil")
				}
			}
		})
	}
}

func TestGetLatestPost(t *testing.T) {
	body := helperLoadBytes(t, "posts.json")
	testCases := []struct {
		title         string
		blogID        string
		bloggerAPIKey string
		resp          *http.Response
		getErr        error
	}{
		{
			http.StatusText(http.StatusOK),
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
				ContentLength: int64(len(body)),
				Header:        make(http.Header, 0),
			},
			nil,
		},
		{
			http.StatusText(http.StatusNotFound),
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{
				Status:        http.StatusText(http.StatusNotFound),
				StatusCode:    http.StatusNotFound,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString(http.StatusText(http.StatusNotFound))),
				ContentLength: int64(len(http.StatusText(http.StatusNotFound))),
				Header:        make(http.Header, 0),
			},
			nil,
		},
		{
			"Get should throw an error",
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{},
			errors.New("No Wi-Fi"),
		},
		{
			"ResolveReference should throw an error",
			"$:$ù%§µ%M£",
			"tBloggerAPIKey",
			&http.Response{},
			nil,
		},
		{
			"json.Decode should throw an error",
			"tBlogID",
			"tBloggerAPIKey",
			&http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString("Not a json")),
				ContentLength: int64(len("Not a json")),
				Header:        make(http.Header, 0),
			},
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprint(tc.title), func(t *testing.T) {
			os.Unsetenv("BLOG_ID")
			os.Unsetenv("BLOGGER_API_KEY")
			if tc.blogID != "" {
				err := os.Setenv("BLOG_ID", tc.blogID)
				if err != nil {
					t.Fatal("Couldn't set env variable BLOG_ID")
				}
			}
			if tc.bloggerAPIKey != "" {
				err := os.Setenv("BLOGGER_API_KEY", tc.bloggerAPIKey)
				if err != nil {
					t.Fatal("Couldn't set env variable BLOGGER_API_KEY")
				}
			}
			mock := &mockHTTPClient{}
			GetFunc = func(u string) (*http.Response, error) {
				if tc.getErr != nil {
					return nil, tc.getErr
				}
				return tc.resp, nil
			}
			blogger, err := NewBlogger(mock)
			if err != nil {
				t.Fatal("Error raised with NewBlogger")
			}
			resultData, resultErr := blogger.GetLatestPost()

			if tc.title == http.StatusText(http.StatusOK) {
				if resultData == nil {
					t.Fatal("Got: resp == nil. Expected: resp != nil")
				} else {
					expectedBytes := helperLoadBytes(t, "posts.json")
					expected := struct {
						Items []models.Post `json:"items"`
					}{}
					err = json.Unmarshal(expectedBytes, &expected)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(expected.Items[0], *resultData) {
						t.Fatalf("Got: %s. Expected: %s", *resultData, expected.Items[0])
					}
				}
				if resultErr != nil {
					t.Fatalf("Got: err = %s. Expected: err == nil", resultErr)
				}
			} else {
				if resultData != nil {
					t.Fatal("Got: resp != nil. Expected: resp == nil")
				}
				if resultErr == nil {
					t.Fatal("Got: err == nil. Expected: err != nil")
				}
			}
		})
	}
}
