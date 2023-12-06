package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		outFile, err := os.CreateTemp("./", "testFile")
		require.NoError(t, err)

		defer func() {
			outFile.Close()
			err = os.Remove(outFile.Name())
			if err != nil {
				fmt.Println("Error removing temp file:", err)
			}
		}()

		var testLimit int64 = 100
		err = Copy("./testdata/input.txt", outFile.Name(), 0, testLimit)
		fileInfo, err := outFile.Stat()
		fileSize := fileInfo.Size()

		require.NoError(t, err)
		require.Equal(t, testLimit, fileSize, "not all bytes were copied")
	})

	t.Run("limit too big", func(t *testing.T) {
		outFile, err := os.CreateTemp("./", "testFile")
		require.NoError(t, err)

		defer func() {
			outFile.Close()
			err = os.Remove(outFile.Name())
			if err != nil {
				fmt.Println("Error removing temp file:", err)
			}
		}()

		var testLimit int64 = 100
		err = Copy("./testdata/out_offset0_limit10.txt", outFile.Name(), 0, testLimit)
		fileInfo, err := outFile.Stat()
		fileSize := fileInfo.Size()

		require.NoError(t, err)
		require.Less(t, fileSize, testLimit)
	})

	t.Run("error with too big offset ", func(t *testing.T) {
		outFile, err := os.CreateTemp("./", "testFile")
		require.NoError(t, err)

		defer func() {
			outFile.Close()
			err = os.Remove(outFile.Name())
			if err != nil {
				fmt.Println("Error removing temp file:", err)
			}
		}()

		var testLimit int64 = 100
		var testOffset int64 = 20
		err = Copy("./testdata/out_offset0_limit10.txt", outFile.Name(), testOffset, testLimit)

		require.EqualError(t, ErrOffsetExceedsFileSize, err.Error())
	})
}
