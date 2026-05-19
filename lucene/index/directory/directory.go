package directory

import (
	"os"
	"path/filepath"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type Directory interface {
	OpenOutputStream(string) (*internal.OutputStream, error)
	OpenInputStream(string) (*internal.InputStream, error)
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

func (d *FSDirectory) OpenOutputStream(filename string) (*internal.OutputStream, error) {
	path := d.resolve(filename)
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return internal.NewOutputStream(f), nil
}

func (d *FSDirectory) OpenInputStream(filename string) (*internal.InputStream, error) {
	path := d.resolve(filename)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return internal.NewInputStream(f), nil
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
