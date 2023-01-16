package zipper

import (
	"archive/zip"
	"fmt"
	"github.com/gozelle/vfs"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type WriteableFileSystem interface {
	http.FileSystem
	Write(dir, file string) error
}

type Option func(c *Config)

type Config struct {
	sourceReader io.Reader
	sourceDir    string
	sourceFile   string
	sourceFs     http.FileSystem
	targetWriter io.Writer
	targetFs     WriteableFileSystem
	targetFile   string
	targetForce  bool
}

func (c Config) zipValid() error {
	
	if c.sourceReader == nil &&
		c.sourceDir == "" &&
		c.sourceFile == "" &&
		c.sourceFs == nil {
		return fmt.Errorf("no souce config")
	}
	
	if c.targetFs != nil {
		return fmt.Errorf("zip not suprt file system target")
	}
	
	if c.targetWriter == nil &&
		c.targetFile == "" {
		return fmt.Errorf("no target config")
	}
	
	return nil
}

func WithSourceFile(file string) Option {
	return func(c *Config) {
		c.sourceFile = file
	}
}

func WithSourceDir(dir string) Option {
	return func(c *Config) {
		c.sourceDir = dir
	}
}

func WithSourceReader(reader io.Reader) Option {
	return func(c *Config) {
		c.sourceReader = reader
	}
}

func WithSourceFileSystem(fs http.FileSystem) Option {
	return func(c *Config) {
		c.sourceFs = fs
	}
}

func WithTargetWriter(writer io.Writer) Option {
	return func(c *Config) {
		c.targetWriter = writer
	}
}

func WithTargetFile(file string) Option {
	return func(c *Config) {
		c.targetFile = file
	}
}

func WithTargetFileSystem(wfs WriteableFileSystem) Option {
	return func(c *Config) {
		c.targetFs = wfs
	}
}

func WithTargetForce(force bool) Option {
	return func(c *Config) {
		c.targetForce = force
	}
}

func Zip(options ...Option) (err error) {
	if len(options) == 0 {
		err = fmt.Errorf("no options")
		return
	}
	
	c := &Config{}
	for _, v := range options {
		v(c)
	}
	err = c.zipValid()
	if err != nil {
		return
	}
	
	var out io.Writer
	if c.targetWriter != nil {
		out = c.targetWriter
	}
	
	if c.sourceFs != nil {
		err = zipFileSystem(c.sourceFs, out)
		if err != nil {
			return
		}
	}
	
	return
}

func zipFileSystem(fs http.FileSystem, out io.Writer) (err error) {
	
	zw := zip.NewWriter(out)
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
