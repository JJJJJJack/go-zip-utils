package goziputils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

type ZipMap map[string]*zip.File

// NewZipMapFromFilename creates a ZipMap from the given zip filename.
func NewZipMapFromFilename(filename string) (ZipMap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	zipMap := make(ZipMap)
	for _, f := range r.File {
		zipMap[f.Name] = f
	}

	return zipMap, nil
}

// NewZipMapFromBytes creates a ZipMap from the given byte slice containing zip data.
func NewZipMapFromBytes(data []byte) (ZipMap, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	zipMap := make(ZipMap)
	for _, f := range r.File {
		zipMap[f.Name] = f
	}

	return zipMap, nil
}

// CopyFile copies the file f into zipWriter, preserving its
// original name and compression method.
func CopyFile(zipWriter *zip.Writer, f *zip.File) error {
	fileInZip, err := f.Open()
	if err != nil {
		return err
	}
	defer fileInZip.Close()

	newFile, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   f.Name,
		Method: f.FileHeader.Method,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, fileInZip)
	if err != nil {
		return err
	}

	return nil
}

// ReadZipFileContent reads the content of a file in a zip archive.
func ReadZipFileContent(f *zip.File) ([]byte, error) {
	fileInZip, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer fileInZip.Close()

	return io.ReadAll(fileInZip)
}

// RewriteFileIntoZipWriter reacreates the file f in zipWriter
// with the provided content while preserving the f.FileHeader.
func RewriteFileIntoZipWriter(zipWriter *zip.Writer, f *zip.File, content []byte) error {
	newHeader := f.FileHeader

	newHeader.UncompressedSize64 = uint64(len(content))
	newHeader.CompressedSize64 = 0

	newFile, err := zipWriter.CreateHeader(&newHeader)
	if err != nil {
		return err
	}

	_, err = newFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// WriteFile creates a new file in the zip archive with the given filename
func WriteFile(zipWriter *zip.Writer, filename string, content []byte) error {
	newFile, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:               filename,
		Method:             zip.Store,
		UncompressedSize64: uint64(len(content)),
	})
	if err != nil {
		return err
	}

	_, err = newFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}
