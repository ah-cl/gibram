// Package types provides tests for core data structures
package types

import (
	"sync"
	"testing"
)

// =============================================================================
// IDGenerator Tests
// =============================================================================

func TestIDGenerator_New(t *testing.T) {
	gen := NewIDGenerator()
	if gen == nil {
		t.Fatal("NewIDGenerator() returned nil")
	}
}

func TestIDGenerator_NextDocumentID(t *testing.T) {
	gen := NewIDGenerator()

	id1 := gen.NextDocumentID()
	id2 := gen.NextDocumentID()

	if id1 != 1 {
		t.Errorf("First document ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("Second document ID = %d, want 2", id2)
	}
}

func TestIDGenerator_NextTextUnitID(t *testing.T) {
	gen := NewIDGenerator()

	id1 := gen.NextTextUnitID()
	id2 := gen.NextTextUnitID()

	if id1 != 1 {
		t.Errorf("First text unit ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("Second text unit ID = %d, want 2", id2)
	}
}

func TestIDGenerator_NextEntityID(t *testing.T) {
	gen := NewIDGenerator()

	id1 := gen.NextEntityID()
	id2 := gen.NextEntityID()

	if id1 != 1 {
		t.Errorf("First entity ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("Second entity ID = %d, want 2", id2)
	}
}

func TestIDGenerator_NextRelationshipID(t *testing.T) {
	gen := NewIDGenerator()

	id1 := gen.NextRelationshipID()
	id2 := gen.NextRelationshipID()

	if id1 != 1 {
		t.Errorf("First relationship ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("Second relationship ID = %d, want 2", id2)
	}
}

func TestIDGenerator_NextCommunityID(t *testing.T) {
	gen := NewIDGenerator()

	id1 := gen.NextCommunityID()
	id2 := gen.NextCommunityID()

	if id1 != 1 {
		t.Errorf("First community ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("Second community ID = %d, want 2", id2)
	}
}

func TestIDGenerator_NextQueryID(t *testing.T) {
	gen := NewIDGenerator()

	id1 := gen.NextQueryID()
	id2 := gen.NextQueryID()

	if id1 != 1 {
		t.Errorf("First query ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("Second query ID = %d, want 2", id2)
	}
}

func TestIDGenerator_SetCounters(t *testing.T) {
	gen := NewIDGenerator()

	gen.SetCounters(10, 20, 30, 40, 50, 60)

	if gen.NextDocumentID() != 11 {
		t.Error("Document counter not restored correctly")
	}
	if gen.NextTextUnitID() != 21 {
		t.Error("TextUnit counter not restored correctly")
	}
	if gen.NextEntityID() != 31 {
		t.Error("Entity counter not restored correctly")
	}
	if gen.NextRelationshipID() != 41 {
		t.Error("Relationship counter not restored correctly")
	}
	if gen.NextCommunityID() != 51 {
		t.Error("Community counter not restored correctly")
	}
	if gen.NextQueryID() != 61 {
		t.Error("Query counter not restored correctly")
	}
}

func TestIDGenerator_GetCounters(t *testing.T) {
	gen := NewIDGenerator()

	gen.NextDocumentID()
	gen.NextDocumentID()
	gen.NextTextUnitID()
	gen.NextEntityID()
	gen.NextRelationshipID()
	gen.NextCommunityID()
	gen.NextQueryID()

	doc, tu, ent, rel, comm, query := gen.GetCounters()

	if doc != 2 {
		t.Errorf("Document counter = %d, want 2", doc)
	}
	if tu != 1 {
		t.Errorf("TextUnit counter = %d, want 1", tu)
	}
	if ent != 1 {
		t.Errorf("Entity counter = %d, want 1", ent)
	}
	if rel != 1 {
		t.Errorf("Relationship counter = %d, want 1", rel)
	}
	if comm != 1 {
		t.Errorf("Community counter = %d, want 1", comm)
	}
	if query != 1 {
		t.Errorf("Query counter = %d, want 1", query)
	}
}

func TestIDGenerator_Concurrent(t *testing.T) {
	gen := NewIDGenerator()
	var wg sync.WaitGroup
	const n = 100

	// Concurrent ID generation
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gen.NextDocumentID()
			gen.NextTextUnitID()
			gen.NextEntityID()
			gen.NextRelationshipID()
			gen.NextCommunityID()
			gen.NextQueryID()
		}()
	}

	wg.Wait()

	doc, tu, ent, rel, comm, _ := gen.GetCounters()

	if doc != n {
		t.Errorf("Document counter = %d, want %d", doc, n)
	}
	if tu != n {
		t.Errorf("TextUnit counter = %d, want %d", tu, n)
	}
	if ent != n {
		t.Errorf("Entity counter = %d, want %d", ent, n)
	}
	if rel != n {
		t.Errorf("Relationship counter = %d, want %d", rel, n)
	}
	if comm != n {
		t.Errorf("Community counter = %d, want %d", comm, n)
	}
}
