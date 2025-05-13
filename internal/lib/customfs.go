package lib

import (
	"net/http"
	"path/filepath"
)

type NeuteredFileSystem struct {
	Fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.Fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.Fs.Open(index); err != nil {
			closeEr := f.Close()
			if closeEr != nil {
				return nil, closeEr
			}

			return nil, err
		}
	}

	return f, nil
}
