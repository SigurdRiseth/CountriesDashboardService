package testutils

import (
	"cloud.google.com/go/firestore"
	"context"
)

// MockFirestoreClient simulates Firestore client behavior for testing.
type MockFirestoreClient struct {
	GetFunc    func(ctx context.Context, id string) (*firestore.DocumentSnapshot, error)
	DeleteFunc func(ctx context.Context, id string) error
}

// Simulate Firestore.Collection().Doc(id).Get()
func (m *MockFirestoreClient) Get(ctx context.Context, id string) (*firestore.DocumentSnapshot, error) {
	return m.GetFunc(ctx, id)
}

// Simulate Firestore.Collection().Doc(id).Delete()
func (m *MockFirestoreClient) Delete(ctx context.Context, id string) error {
	return m.DeleteFunc(ctx, id)
}

func MockCheckWebhooks(country, event, id string) {
	// Mock implementation of webhook checking
	return
}
