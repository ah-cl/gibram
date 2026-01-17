// Package backup provides backup and recovery for GibRAM
package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Archiver handles creating and extracting backup archives
type Archiver struct {
	baseDir string
}

// NewArchiver creates a new archiver
func NewArchiver(baseDir string) *Archiver {
	return &Archiver{baseDir: baseDir}
}

// Archive creates a tar.gz archive of the data directory
func (a *Archiver) Archive(outputPath string) error {
	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer f.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(f)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Walk directory and add files
	return filepath.Walk(a.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Use relative path
		relPath, err := filepath.Rel(a.baseDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			return err
		}

		return nil
	})
}

// Extract extracts a tar.gz archive to the data directory
func (a *Archiver) Extract(archivePath string) error {
	// Open archive
	f, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("open archive: %w", err)
	}
	defer f.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar next: %w", err)
		}

		targetPath := filepath.Join(a.baseDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}

			// Create file
			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			// Copy content
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// ArchiveInfo holds information about an archive
type ArchiveInfo struct {
	Path      string
	Size      int64
	ModTime   time.Time
	FileCount int
}

// ListArchives lists all archives in a directory
func ListArchives(dir string) ([]*ArchiveInfo, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.tar.gz"))
	if err != nil {
		return nil, err
	}

	infos := make([]*ArchiveInfo, 0, len(files))
	for _, path := range files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		infos = append(infos, &ArchiveInfo{
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	return infos, nil
}

// VerifyArchive verifies archive integrity
func VerifyArchive(archivePath string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("invalid gzip: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	// Try to read all entries
	for {
		_, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("invalid tar: %w", err)
		}
	}

	return nil
}
