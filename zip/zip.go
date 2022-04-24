package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jrp0h/backpack/utils"
)

func walkDir(path string, writer *zip.Writer, first bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		if first {
			header.Name, err = filepath.Rel(filepath.Dir(path), entry.Name())
			if err != nil {
				return err
			}
		} else {
			header.Name = strings.Join(strings.Split(path, "/")[1:], "/")
			header.Name = filepath.Join(header.Name, entry.Name())
		}

		if info.IsDir() {
			header.Name += "/"
		}

		header.SetMode(info.Mode())

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if err = walkDir(filepath.Join(path, entry.Name()), writer, false); err != nil {
				return err
			}
			continue
		}

		f, err := os.Open(filepath.Join(path, entry.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = io.Copy(headerWriter, f); err != nil {
			return err
		}

	}

	return nil
}

func Zip(input, output string) (outErr error) {

	if utils.PathExists(output) {
		return fmt.Errorf("output path %s already exists", output)
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	return walkDir(input, writer, true)
}

// Taken from https://golangcode.com/unzip-files-in-go/
func Unzip(input, output string) (outErr error) {

	r, err := zip.OpenReader(input)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(output, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(output)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			if err = os.MkdirAll(fpath, f.Mode()); err != nil {
				return err
			}

			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), f.Mode()); err != nil {
			return err
		}

		if utils.PathExists(fpath) {
			return fmt.Errorf("%s already exists", fpath)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
