package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := info.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > fileSize {
		limit = fileSize
	}

	reader := io.LimitReader(file, limit)

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(reader)
	defer bar.Finish()

	file.Seek(offset, io.SeekStart)

	createdFile, err := create(toPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(createdFile, barReader)
	if err != nil {
		return err
	}

	return nil
}
