package main

import (
	"bytes"
	"io"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	fromFilePath := "testdata/input.txt"

	t.Run("copy all", func(t *testing.T) {
		// tempFile, err := os.CreateTemp("", "out.txt")
		// if err != nil {
		// 	t.Fatalf("Failed to create temp file: %v", err)
		// }

		// Arrange
		outFilePath := "testdata/out1.txt"
		checkFilePath := fromFilePath

		// Act
		copyErr := Copy(fromFilePath, outFilePath, 0, 0)
		defer os.Remove(outFilePath)
		isSame, compareErr := deepCompare(checkFilePath, outFilePath)

		// Assert
		require.NoError(t, copyErr)
		require.NoError(t, compareErr)
		assert.True(t, isSame, "files should be the same")
	})

	t.Run("copy limit without offset", func(t *testing.T) {
		// Arrange
		outFilePath := "testdata/out2.txt"
		checkFilePath := "testdata/out_offset0_limit1000.txt"
		offset := int64(0)
		limit := int64(1000)

		// Act
		copyErr := Copy(fromFilePath, outFilePath, offset, limit)
		defer os.Remove(outFilePath)
		isSame, compareErr := deepCompare(checkFilePath, outFilePath)
		outFileNumberOfBytes, bytesErr := getFileNumberOfBytes(outFilePath)

		// Assert
		require.NoError(t, copyErr)
		require.NoError(t, compareErr)
		require.NoError(t, bytesErr)
		assert.True(t, isSame, "files should be the same")
		assert.Equal(t, outFileNumberOfBytes, limit, "coped number of bytes should be equal to limit")
	})

	t.Run("copy limit with offset", func(t *testing.T) {
		// Arrange
		outFilePath := "testdata/out3.txt"
		checkFilePath := "testdata/out_offset100_limit1000.txt"
		offset := int64(100)
		limit := int64(1000)

		// Act
		copyErr := Copy(fromFilePath, outFilePath, offset, limit)
		defer os.Remove(outFilePath)
		isSame, compareErr := deepCompare(checkFilePath, outFilePath)
		outFileNumberOfBytes, bytesErr := getFileNumberOfBytes(outFilePath)

		// Assert
		require.NoError(t, copyErr)
		require.NoError(t, compareErr)
		require.NoError(t, bytesErr)
		assert.True(t, isSame, "files should be the same")
		assert.Equal(t, outFileNumberOfBytes, limit, "coped number of bytes should be equal to limit")
	})

	t.Run("copy all with offset", func(t *testing.T) {
		// Arrange
		outFilePath := "testdata/out4.txt"
		checkFilePath := "testdata/out_offset100_limit0.txt"
		offset := int64(100)
		limit := int64(0)

		// Act
		copyErr := Copy(fromFilePath, outFilePath, offset, limit)
		defer os.Remove(outFilePath)
		isSame, compareErr := deepCompare(checkFilePath, outFilePath)

		// Assert
		require.NoError(t, copyErr)
		require.NoError(t, compareErr)
		assert.True(t, isSame, "files should be the same")
	})

	t.Run("copy limit with offset if limit > file size", func(t *testing.T) {
		// Arrange
		outFilePath := "testdata/out5.txt"
		checkFilePath := "testdata/out_offset6000_limit1000.txt"
		offset := int64(6000)
		limit := int64(1000)

		// Act
		copyErr := Copy(fromFilePath, outFilePath, offset, limit)
		defer os.Remove(outFilePath)
		isSame, compareErr := deepCompare(checkFilePath, outFilePath)
		outFileNumberOfBytes, bytesErr := getFileNumberOfBytes(outFilePath)

		// Assert
		require.NoError(t, copyErr)
		require.NoError(t, compareErr)
		require.NoError(t, bytesErr)
		assert.True(t, isSame, "files should be the same")
		assert.Less(t, outFileNumberOfBytes, limit, "coped number of bytes should be less than limit")
	})

	t.Run("copy limit with offset > file size", func(t *testing.T) {
		// Arrange
		outFilePath := "testdata/out6.txt"
		offset := int64(math.MaxInt64)
		limit := int64(0)

		// Act
		copyErr := Copy(fromFilePath, outFilePath, offset, limit)
		defer os.Remove(outFilePath)

		// Assert
		assert.Error(t, copyErr)
		assert.ErrorIs(t, copyErr, ErrOffsetExceedsFileSize)
	})
}

const chunkSize = 1024

func deepCompare(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil // Both files reached EOF, they are the same
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil // One file reached EOF before the other, they are not the same
			}

			return false, err1 // Both files returned an error, they are not the same
		}

		if !bytes.Equal(b1, b2) {
			return false, nil // Chunks are not equal, files are not the same
		}
	}
}

func getFileNumberOfBytes(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
