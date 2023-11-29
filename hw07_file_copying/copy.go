package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrFileNotFound          = errors.New("file is not found")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	inFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotFound
		}
		return fmt.Errorf("error while opening inFile: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error while creating out file: %w", err)
	}
	defer outFile.Close()

	fileInfo, err := inFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()

	if limit > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit > fileSize {
		limit = fileSize
	}

	var count int64
	if fileSize-offset < limit {
		count = fileSize - offset
	} else {
		count = limit
	}

	bar := pb.Full.Start64(count)
	barReader := bar.NewProxyReader(inFile)

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(outFile)

	_, err = io.CopyN(multiWriter, barReader, limit)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	bar.Finish()

	return nil
}
