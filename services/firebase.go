package services

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var App *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile("firebase-service-account.json")
	var err error
	App, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Firebase init error: %v", err)
	}
}
