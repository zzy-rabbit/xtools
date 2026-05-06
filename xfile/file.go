package xfile

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func IsExist(ctx context.Context, path string) bool {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		return true
	}
	return true
}

func GetFileSize(ctx context.Context, file string) (int64, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func Zip(ctx context.Context, zipPath string, paths []string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	for _, p := range paths {
		err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// 关键：保留权限
			header.SetMode(info.Mode())

			// 保留目录结构（相对路径）
			relPath, err := filepath.Rel(filepath.Dir(p), path)
			if err != nil {
				return err
			}
			header.Name = relPath

			// 目录需要以 / 结尾
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}

			// 目录不需要写内容
			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func Unzip(ctx context.Context, zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	unzipFile := func(f *zip.File, destDir string) error {
		filePath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", filePath)
		}

		if f.FileInfo().IsDir() {
			return os.MkdirAll(filePath, f.Mode())
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		srcFile, err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		return os.Chmod(filePath, f.Mode())
	}

	for _, f := range r.File {
		err = unzipFile(f, destDir)
		if err != nil {
			return err
		}
	}
	return nil
}
