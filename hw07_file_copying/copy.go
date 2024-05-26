package main

import (
	"errors"
	"io"
	"os"

	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	readFileMode := os.FileMode(0o444)
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, readFileMode)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	// Проверяем, что offset не превышает размер файла
	fileInfo, err := fromFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Перемещаемся к указанному offset
	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	// Определяем, сколько байт нужно скопировать
	var bytesToCopy int64
	if limit == 0 || limit > fileInfo.Size()-offset {
		bytesToCopy = fileInfo.Size() - offset
	} else {
		bytesToCopy = limit
	}

	// Создаем прогресс-бар
	p := mpb.New(mpb.WithWidth(60))
	bar := p.AddBar(bytesToCopy,
		mpb.PrependDecorators(
			decor.Name("Copying: "),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.CountersKiloByte("% .2f / % .2f"),
		),
	)

	// Копируем данные
	buf := make([]byte, 1024)
	for bytesToCopy > 0 {
		n, err := fromFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if int64(n) > bytesToCopy {
			n = int(bytesToCopy)
		}
		if _, err := toFile.Write(buf[:n]); err != nil {
			return err
		}
		bytesToCopy -= int64(n)

		bar.IncrBy(n)
	}

	p.Wait()
	return nil
}
