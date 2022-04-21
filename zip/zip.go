package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jrp0h/backuper/utils"
)

// Taken from https://gosamples.dev/zip-file/
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

  	return filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 3. Create a local file header
        header, err := zip.FileInfoHeader(info)
        if err != nil {
            return err
        }

        // set compression
        header.Method = zip.Deflate

        // 4. Set relative path of a file as the header name
        header.Name, err = filepath.Rel(filepath.Dir(input), path)
        if err != nil {
            return err
        }
        if info.IsDir() {
            header.Name += "/"
        }

        // 5. Create writer for the file header and save content of the file
        headerWriter, err := writer.CreateHeader(header)
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()

        _, err = io.Copy(headerWriter, f)
        return err
    })
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
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        // Make File
        if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
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