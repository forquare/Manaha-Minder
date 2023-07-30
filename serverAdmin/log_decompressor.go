package serverAdmin

import (
	"compress/gzip"
	"github.com/forquare/manaha-minder/config"
	logger "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	root string
)

func LogDecompressor() {
	logger.Debug("Starting Log Decompressor")
	config := config.GetConfig()
	root = config.MinecraftServer.LogDir
	fileSystem := os.DirFS(root)

	fs.WalkDir(fileSystem, ".", walk)
}

func walk(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() && filepath.Ext(d.Name()) == ".gz" {
		decompress(filepath.Join(root, p))
	}
	return nil
}

func decompress(f string) {
	gzippedFile, err := os.Open(f)
	if err != nil {
		logger.Error(err)
	}

	defer gzippedFile.Close()

	// Create a new gzip reader
	gzipReader, err := gzip.NewReader(gzippedFile)
	defer gzipReader.Close()

	// Create a new file to hold the uncompressed data
	dirname := filepath.Dir(f)
	basename := filepath.Base(f)
	decompressedName := strings.TrimSuffix(basename, filepath.Ext(basename))
	decompressedFile, err := os.Create(filepath.Join(dirname, decompressedName))
	if err != nil {
		logger.Error(err)
	}
	defer decompressedFile.Close()

	// Copy the contents of the gzip reader to the new file
	_, err = io.Copy(decompressedFile, gzipReader)
	if err != nil {
		logger.Error(err)
	}

	os.Remove(f)
}
