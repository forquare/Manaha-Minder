package serverAdmin

import (
	"compress/gzip"
	logger "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"manaha_minder/config"
	"os"
	"path/filepath"
)

func LogDecompressor() {
	logger.Debug("Starting Log Decompressor")
	config := config.GetConfig()
	root := config.MinecraftServer.LogDir
	fileSystem := os.DirFS(root)

	fs.WalkDir(fileSystem, ".", walk)
}

func walk(p string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() && filepath.Ext(d.Name()) == ".gz" {
		decompress(filepath.Join(p, d.Name()))
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
	// TODO Fix this
	uncompressedFile, err := os.Create("example.txt")
	if err != nil {
		logger.Error(err)
	}
	defer uncompressedFile.Close()

	// Copy the contents of the gzip reader to the new file
	_, err = io.Copy(uncompressedFile, gzipReader)
	if err != nil {
		logger.Error(err)
	}

	os.Remove(f)
}
