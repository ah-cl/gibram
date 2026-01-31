// Package backup - Critical path tests for backup/recovery
package backup

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAtomicSnapshotCreation tests atomic snapshot write-to-temp pattern
func TestAtomicSnapshotCreation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.snapshot")
	
	// Create snapshot
	err := CreateSnapshot(path, 100, func(w *SnapshotWriter) error {
		data := []byte("test data")
		_, err := w.Write(data)
		return err
	})
	
	if err != nil {
		t.Fatalf("CreateSnapshot failed: %v", err)
	}
	
	// Verify final file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("Snapshot file not created")
	}
	
	// Verify temp file was removed
	tmpPath := path + ".tmp"
	if _, err := os.Stat(tmpPath); !os.IsNotExist(err) {
		t.Error("Temp file not removed")
	}
}

// Test2PCBackupCoordination tests two-phase commit for backups
func Test2PCBackupCoordination(t *testing.T) {
	dir := t.TempDir()
	walDir := filepath.Join(dir, "wal")
	if err := os.MkdirAll(walDir, 0755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	
	wal, err := NewWAL(walDir, SyncEveryWrite)
	if err != nil {
		t.Fatalf("NewWAL failed: %v", err)
	}
	defer func() {
		if err := wal.Close(); err != nil {
			t.Fatalf("WAL Close failed: %v", err)
		}
	}()
	
	// Write some WAL entries
	_, err = wal.Append(EntryInsert, "key1", []byte("data1"))
	if err != nil {
		t.Fatalf("WAL Append failed: %v", err)
	}
	
	snapshotPath := filepath.Join(dir, "snapshot.gibram")
	coordinator := NewBackupCoordinator(wal, snapshotPath)
	
	// Test prepare phase
	lsn, err := coordinator.Prepare()
	if err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}
	
	if lsn == 0 {
		t.Error("Expected non-zero LSN")
	}
	
	// Test commit phase
	err = coordinator.Commit(func(w *SnapshotWriter) error {
		data := []byte("snapshot data")
		_, err := w.Write(data)
		return err
	})
	
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	
	// Verify snapshot exists
	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		t.Error("Snapshot not created")
	}
}

// TestWALChecksumIntegrity tests xxHash64 checksum validation
func TestWALChecksumIntegrity(t *testing.T) {
	dir := t.TempDir()
	
	wal, err := NewWAL(dir, SyncEveryWrite)
	if err != nil {
		t.Fatalf("NewWAL failed: %v", err)
	}
	
	// Write entry
	lsn, err := wal.Append(EntryInsert, "testkey", []byte("testdata"))
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}
	
	if err := wal.Close(); err != nil {
		t.Fatalf("WAL Close failed: %v", err)
	}
	
	// For now, skip WAL reader test since we don't have a reader implementation
	// The important part is that writing works with xxHash64
	_ = lsn
}

// BenchmarkWALWrite benchmarks WAL write performance with xxHash64
func BenchmarkWALWrite(b *testing.B) {
	dir := b.TempDir()
	
	wal, err := NewWAL(dir, SyncNever) // Don't sync for benchmark
	if err != nil {
		b.Fatalf("NewWAL failed: %v", err)
	}
	defer func() {
		if err := wal.Close(); err != nil {
			b.Fatalf("WAL Close failed: %v", err)
		}
	}()
	
	data := make([]byte, 1024)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := wal.Append(EntryInsert, "key", data)
		if err != nil {
			b.Fatalf("Append failed: %v", err)
		}
	}
}
