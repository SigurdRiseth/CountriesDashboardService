package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

// Firebase context and client used by Firestore functions throughout the program.
var (
	Ctx    context.Context
	Client *firestore.Client
)

// InitFirebase initializes the Firebase app and Firestore client.
// It sets up the context and client variables used throughout the program.
// If initialization fails, the function logs a fatal error and terminates the program.
func InitFirebase() {

	// Initialize Firebase app
	Ctx = context.Background()
	sa := option.WithCredentialsFile("./firebase-adminsdk.json")

	app, err := firebase.NewApp(Ctx, nil, sa)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Initialize Firestore client
	Client, err = app.Firestore(Ctx)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
}

// CloseFirebase closes the Firestore client if it is not nil.
// This function should be called to properly release resources
// associated with the Firestore client.
func CloseFirebase() {
	if Client != nil {
		_ = Client.Close()
	}
}
