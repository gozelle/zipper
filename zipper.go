package zipper

import (
	"archive/zip"
	"github.com/gozelle/vfs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Zip(dst io.Writer, fs http.FileSystem) (err error) {
	
	zw := zip.NewWriter(dst)
	defer func() {
		_ = zw.Close()
	}()
	
	return vfs.Walk(fs, "/", func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))
		if fi.IsDir() {
			fh.Name += "/"
		}
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}
		if !fh.Mode().IsRegular() {
			return
		}
		fr, err := os.Open(path)
		defer func() {
			_ = fr.Close()
		}()
		if err != nil {
			return
		}
		_, err = io.Copy(w, fr)
		if err != nil {
			return
		}
		
		return
	})
}

func ZipToFile(dst string, fs http.FileSystem) (err error) {
	f, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	return Zip(f, fs)
}
