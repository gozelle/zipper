package tests

import (
	"github.com/gozelle/vfs"
	"github.com/gozelle/zipper"
	"testing"
)

func TestZip(t *testing.T) {
	err := zipper.Zip(
		zipper.WithSourceFileSystem(Templates),
		zipper.WithTargetFile("test.zip"),
		zipper.WithTargetForce(true),
	)
	if err != nil {
		panic(err)
	}
}

func TestZip2(t *testing.T) {
	
	fs := vfs.NewFS()
	
	fs.Add("/", "1.txt", []byte("1"))
	fs.Add("/", "2.txt", []byte("2"))
	fs.Add("/", "3.txt", []byte("2"))
	
	err := zipper.Zip(
		zipper.WithSourceFileSystem(fs),
		zipper.WithTargetFile("test2.zip"),
		zipper.WithTargetForce(true),
	)
	if err != nil {
		panic(err)
	}
}
