package testutils

import "cloud.google.com/go/firestore"

// MockAddToCollection simulates adding a document to a Firestore collection.
// It returns a mock DocumentRef (with id = "mockID"), WriteResult, and no error.
func MockAddToCollection(collection string, data any) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	return &firestore.DocumentRef{ID: "mockID"}, &firestore.WriteResult{}, nil
}

func MockDeleteDocument(docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	return &firestore.WriteResult{}, nil
}

func MockGetDocumentRef(collection, docID string) *firestore.DocumentRef {
	return &firestore.DocumentRef{ID: "mockID"}
}

func MockGetDocumentByRef(docRef *firestore.DocumentRef) (*firestore.DocumentSnapshot, error) {
	return &firestore.DocumentSnapshot{}, nil
}

func MockFirebaseClientInitialized() bool {
	return true
}

func MockDocumentExists(docRef *firestore.DocumentSnapshot) bool {
	// Simulate checking if a document exists
	return true
}
