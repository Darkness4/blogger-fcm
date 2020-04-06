package services

import (
	"context"
	"me/blogger-fcm/models"

	"firebase.google.com/go/messaging"
)

// SendLatestPost to FCM
func SendLatestPost(ctx context.Context, fcm *messaging.Client, post *models.Post) (string, error) {
	notification := &messaging.Notification{
		Body:  post.Title,
		Title: "Nouvelle article !",
	}
	message := &messaging.Message{
		Notification: notification,
		Data: map[string]string{
			"kind":      post.Kind,
			"id":        post.ID,
			"published": post.Published.String(),
			"updated":   post.Updated.String(),
			"url":       post.URL,
			"selfLink":  post.SelfLink,
			"title":     post.Title,
			"content":   post.Content,
		},
		Topic: "actualite",
	}
	return fcm.Send(ctx, message)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println("Successfully sent message:", response)
}
