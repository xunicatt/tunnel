package zippy

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzipFile(rootDir string, file *zip.File, dest string) (err error) {
	if rootDir == file.Name {
		os.MkdirAll(dest, os.ModePerm)
		return
	}

	dir := strings.Replace(file.Name, rootDir, "", 1)
	if file.FileInfo().IsDir() {
		err = os.MkdirAll(filepath.Join(dest, dir), os.ModePerm)
		return
	}

	f, err := file.Open()
	if err != nil {
		return
	}

	defer f.Close()

	writeFile, err := os.OpenFile(filepath.Join(dest, dir), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return
	}

	defer writeFile.Close()

	_, err = io.Copy(writeFile, f)
	return
}

func Unzip(src, dest string) (err error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return
	}

	defer r.Close()
	mainDir := r.File[0].Name
	for _, file := range r.File {
		unzipFile(mainDir, file, dest)
	}

	return
}
