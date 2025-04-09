package firebase

import (
	"CountriesDashboardService/consts"
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

// Add this line below your imports in firebase.go
const NotificationsCollection = "notifications"

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

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(Ctx, nil, opt)
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

// GetCollectionIterator returns a DocumentIterator for the specified Firestore collection.
func GetCollectionIterator(collection string) *firestore.DocumentIterator {
	return Client.Collection(collection).Documents(Ctx)
}

// GetDocumentRef returns a DocumentRef for the specified Firestore collection and document ID.
func GetDocumentRef(collection, docID string) *firestore.DocumentRef {
	return Client.Collection(collection).Doc(docID)
}

// GetDocument retrieves a DocumentSnapshot for the specified Firestore collection and document ID.
func GetDocument(collection, docID string) (*firestore.DocumentSnapshot, error) {
	return GetDocumentRef(collection, docID).Get(Ctx)
}

// GetDocumentByRef retrieves a DocumentSnapshot for the specified DocumentRef.
func GetDocumentByRef(docRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	return docRef.Get(Ctx)
}

// AddToCollection adds a new document to the specified Firestore collection and returns the resulting DocumentRef, WriteResult and error.
func AddToCollection(collection string, data any) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	return Client.Collection(collection).Add(Ctx, data)
}

// DeleteDocument deletes a document from the specified Firestore DocumentRef and returns the resulting WriteResult and error.
func DeleteDocument(docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	return docRef.Delete(Ctx)
}

// FirebaseClientInitialized checks if the Firestore client is initialized.
// It returns true if the client is initialized, false otherwise.
func FirebaseClientInitialized() bool {
	if Client == nil {
		log.Println(consts.FTNotInitialized)
		return false
	}
	return true
}

func DocumentExists(doc *firestore.DocumentSnapshot) bool {
	return doc.Exists()
}

// NewClient initializes and returns a Firestore client using the given context and project ID.
// This is used mainly in tests where we want to manually manage the client.
func NewClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	return firestore.NewClient(ctx, projectID)
}
