package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

// Firebase context and client used by Firestore functions throughout the program.
var (
	ctx    context.Context
	client *firestore.Client
)

// InitFirebase initializes the Firebase app and Firestore client.
// It returns an error if initialization fails.
func InitFirebase() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Check if GOOGLE_APPLICATION_CREDENTIALS environment variable is set
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsFile == "" {
		return fmt.Errorf("environment variable GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	// Initialize Firebase app
	ctx = context.Background()
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	// Initialize Firestore client
	client, err = app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize Firestore client: %v", err)
	}

	return nil
}

// CloseFirebase closes the Firestore client if it is not nil.
// This function should be called to properly release resources.
func CloseFirebase() error {
	if client != nil {
		if err := client.Close(); err != nil {
			return fmt.Errorf("failed to close Firestore client: %v", err)
		}
	}
	return nil
}

// GetCollectionIterator returns a DocumentIterator for the specified Firestore collection.
func GetCollectionIterator(collection string) *firestore.DocumentIterator {
	return client.Collection(collection).Documents(ctx)
}

// GetLimitedSortedDocuments returns a DocumentIterator for the specified Firestore collection, sorted by a field and limited to a specified number of documents.
func GetLimitedSortedDocuments(collection, sortBy string, limit int) *firestore.DocumentIterator {
	return client.Collection(collection).Limit(limit).OrderBy(sortBy, firestore.Asc).Documents(ctx)
}

// GetDocumentRef returns a DocumentRef for the specified Firestore collection and document ID.
func GetDocumentRef(collection, docID string) *firestore.DocumentRef {
	return client.Collection(collection).Doc(docID)
}

// GetDocument retrieves a DocumentSnapshot for the specified Firestore collection and document ID.
func GetDocument(collection, docID string) (*firestore.DocumentSnapshot, error) {
	return GetDocumentRef(collection, docID).Get(ctx)
}

// GetDocumentByRef retrieves a DocumentSnapshot for the specified DocumentRef.
func GetDocumentByRef(docRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	return docRef.Get(ctx)
}

// AddToCollection adds a new document to the specified Firestore collection and returns the resulting DocumentRef, WriteResult and error.
func AddToCollection(collection string, data any) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	return client.Collection(collection).Add(ctx, data)
}

// DeleteDocument deletes a document from the specified Firestore DocumentRef and returns the resulting WriteResult and error.
func DeleteDocument(docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	return docRef.Delete(ctx)
}

// SetDocument sets a document in the specified Firestore collection and returns an error if it fails.
func SetDocument(collection, id string, content any) error {
	_, err := client.Collection(collection).Doc(id).Set(ctx, content)
	if err != nil {
		return fmt.Errorf("error setting document: %v", err)
	}
	return nil
}

// UpdateDocument updates a document with the specified data and returns an error if it fails.
func UpdateDocument(docRef *firestore.DocumentRef, data []firestore.Update) error {
	_, err := docRef.Update(ctx, data)
	if err != nil {
		return fmt.Errorf("error updating document: %v", err)
	}
	return nil
}

// GetAllDocuments retrieves all documents from the specified Firestore collection.
func GetAllDocuments(collection string) ([]*firestore.DocumentSnapshot, error) {
	return client.Collection(collection).Documents(ctx).GetAll()
}

// IsFirebaseClientInitialized checks if the Firestore client is initialized.
func IsFirebaseClientInitialized() bool {
	return client != nil
}

// DocumentExists checks if a document exists in Firestore.
func DocumentExists(doc *firestore.DocumentSnapshot) bool {
	return doc.Exists()
}
