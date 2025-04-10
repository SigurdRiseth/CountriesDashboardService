package testutils

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockFirestoreClient simulates Firestore client behavior for testing purposes.
// It provides mock implementations for Get and Delete methods based on injected functions.
type MockFirestoreClient struct {
	GetFunc    func(ctx context.Context, id string) (*firestore.DocumentSnapshot, error)
	DeleteFunc func(ctx context.Context, id string) error
}

// Get mocks Firestore.Collection().Doc(id).Get() behavior by calling the configured GetFunc.
func (m *MockFirestoreClient) Get(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	return m.GetFunc(ctx, id)
}

// Delete mocks Firestore.Collection().Doc(id).Delete() behavior by calling the configured DeleteFunc.
func (m *MockFirestoreClient) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

// MockCheckWebhooks is a stubbed implementation for webhook checks.
// It's used in tests to bypass actual webhook logic.
func MockCheckWebhooks(country, event, id string) {
	// Mock implementation of webhook checking
	return
}

// MockAddToCollection simulates adding a document to a Firestore collection.
// It returns a mock DocumentRef (with id = "mockID"), WriteResult, and no error.
func MockAddToCollection(collection string, data any) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	return &firestore.DocumentRef{ID: "mockID"}, &firestore.WriteResult{}, nil
}

// MockDeleteDocument simulates deletion of a document and always succeeds.
func MockDeleteDocument(docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	return &firestore.WriteResult{}, nil
}

// MockGetDocumentRef simulates retrieving a Firestore document reference.
func MockGetDocumentRef(collection, docID string) *firestore.DocumentRef {
	return &firestore.DocumentRef{ID: "mockID"}
}

// MockGetDocumentByRef simulates retrieving a document snapshot from a reference.
func MockGetDocumentByRef(docRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	return &firestore.DocumentSnapshot{}, nil
}

// MockFirebaseClientInitialized simulates that the Firebase client is always initialized.
func MockFirebaseClientInitialized() bool {
	return true
}

// MockDocumentExists simulates a check that always confirms a document exists.
func MockDocumentExists(docRef *firestore.DocumentSnapshot) bool {
	// Simulate checking if a document exists
	return true
}

// MockGetDocument simulates retrieving a Firestore document.
// Returns a document for ID "mockID", otherwise returns a NotFound error.
func MockGetDocument(collection, docID string) (*firestore.DocumentSnapshot, error) {
	if docID == "mockID" {
		return &firestore.DocumentSnapshot{}, nil
	}
	return nil, status.Error(codes.NotFound, "Document not found")
}
