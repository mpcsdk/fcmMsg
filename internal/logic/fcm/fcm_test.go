package logic

import (
	"context"
	"log"
	"testing"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func Test_Init(t *testing.T) {
	// ctx := context.Background()
	opt := option.WithCredentialsFile("./token.json")
	config := &firebase.Config{ProjectID: "tantalum-f449b"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	response, err := client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: "title",
			Body:  "body",
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ClickAction: "data",
			},
		},
		Token: "fYdd0R_iWUQ2jlq6k7ZEjN:APA91bGvGtnYJkv0Xy6AYzIQ9kH5FRUeLSYOzkrO9BgvBiUPNE3-rvveHrPxl3TiVjsa63EfbMCY4K1ThJ9gFI0qEQYSn38AZrPC1PwyMW_lLnqHJRkHh7DtMoXziapYXz_QSKxAIOgA",
	})
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	log.Fatal(response)
}
