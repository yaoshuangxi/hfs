package util

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ExtractFile(r *http.Request, input, storePath string, override bool) (fp string, err error) {
	r.ParseMultipartForm(100000000)
	file, fh, err := r.FormFile(input)

	if err != nil {
		log.Printf("get form file failed: %v\n", err)
		return "", err
	}
	if file == nil || fh == nil {
		return "", errors.New("cannot get file data")
	}
	defer file.Close()

	fn := fh.Filename
	if fn == "" {
		return "", errors.New("filename cannot be empty!!!")
	}
	Mkdir(storePath)
	fp = filepath.Join(storePath, fn)

	idx := 0
	newfp := fp
	if _, err = os.Stat(newfp); err == nil {
		if (override) {
			os.Remove(fp)
		} else {
			for err == nil {
				idx++;
				ext := filepath.Ext(newfp)
				newfp = filepath.Join(storePath, strings.TrimRight(fn, ext) + "_" + strconv.Itoa(idx) + ext)
				_, err = os.Stat(newfp)
			}
		}
	}

	if idx > 0 {
		err = os.Rename(fp, newfp)
		if err != nil {
			return "", err
		}
		log.Printf("renamed %s to %s\n", fp, newfp)
	}

	f, err := os.Create(fp)
	if err != nil {
		log.Printf("crate file failed: %v\n", err)
		return "", err
	}
	_, err = io.CopyBuffer(f, file, nil)
	log.Printf("Saved [%s] to [%s]\n", fn, fp)
	if err != nil {
		log.Printf("copy buffer failed: %v\n", err)
		return "", err
	}

	return fp, nil
}
