package main

import (
	"context"
	"encoding/json"
	"log"
	"me/blogger-fcm/services"
	"net/http"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

var totalItems int = 0
var ticker *time.Ticker = nil

func noonTask(b *services.Blogger, fcm *messaging.Client) {
	if ticker == nil {
		period := 6 * time.Hour
		ticker = time.NewTicker(period)
		log.Println("Task will run every ", period)
	}
	for {
		log.Println(time.Now())
		log.Println("Refreshing data...")
		blog, err := b.GetBlog()
		if err != nil {
			log.Fatalln("Cannot get blog : ", err)
		}
		if totalItems != blog.Posts.TotalItems {
			totalItems = blog.Posts.TotalItems
			post, err := b.GetLatestPost()
			if err != nil {
				log.Fatalln("Cannot get latest post : ", err)
			}
			services.SendLatestPost(context.Background(), fcm, post)
			log.Println("New post !")
			log.Println("Number of posts : ", totalItems)
		} else {
			log.Println("No new post available.")
		}
		<-ticker.C
	}
}

func main() {
	// Instanciate HTTPClient
	client := new(http.Client)

	log.Println("Initializing services...")
	ctx := context.Background()
	blogger, err := services.NewBlogger(client)
	if err != nil {
		log.Fatalln("Error initializing Blogger client: ", err)
	}
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln("Error initializing Firebase App: ", err)
	}
	fcmService, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalln("Error getting Messaging client: ", err)
	}
	log.Println("Initialized !")

	log.Println("Fetching for a first time...")
	blog, err := blogger.GetBlog()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Blog name is : ", blog.Name)
	totalItems = blog.Posts.TotalItems
	log.Println("Number of posts : ", totalItems)
	post, err := blogger.GetLatestPost()
	if err != nil {
		log.Fatalln("Cannot get latest post : ", err)
	}
	b, err := json.Marshal(post)
	if err != nil {
		log.Fatalln("Cannot parse to json : ", err)
	}
	log.Println("Latest post : ", string(b))

	// Program task
	call := func() {
		noonTask(blogger, fcmService)
	}
	time.AfterFunc(duration(), call)
	wg.Add(1)
	// Do other things here
	wg.Wait()
}

func duration() time.Duration {
	// Start task at ...
	now := time.Now()
	scheduled := now.Add(30 * time.Second)
	start := time.Date(scheduled.Year(), scheduled.Month(), scheduled.Day(), scheduled.Hour(), scheduled.Minute(), scheduled.Second(), 0, scheduled.Location())
	log.Println("Task will run at ", start)
	if now.After(start) {
		start = start.Add(24 * time.Hour)
	}
	duration := start.Sub(now)
	log.Println("Which means in ", duration)
	return duration
}

var wg sync.WaitGroup
