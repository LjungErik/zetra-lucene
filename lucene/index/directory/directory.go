package directory

import (
	"os"
	"path/filepath"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
	"github.com/LjungErik/zetra-lucene/lucene/internal/stream"
)

type Directory interface {
	OpenOutputStream(string) (internal.DataOutputStream, error)
	OpenInputStream(string) (internal.DataInputStream, error)
	GetEntries() ([]os.DirEntry, error)
}

type FSDirectory struct {
	directory string
}

var _ Directory = (*FSDirectory)(nil)

func OpenFSDirectory(directroy string) *FSDirectory {
	return &FSDirectory{
		directory: directroy,
	}
}

func (d *FSDirectory) OpenOutputStream(filename string) (internal.DataOutputStream, error) {
	path := d.resolve(filename)
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return stream.NewOutputStream(f), nil
}

func (d *FSDirectory) OpenInputStream(filename string) (internal.DataInputStream, error) {
	path := d.resolve(filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return stream.NewInputStream(f), nil
}

func (d *FSDirectory) GetEntries() ([]os.DirEntry, error) {
	entries, err := os.ReadDir(d.directory)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (d *FSDirectory) resolve(filename string) string {
	return filepath.Join(d.directory, filename)
}
