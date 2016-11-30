// Copyright 2016 carsonsx. All rights reserved.

package main

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

func mkdir(dir string) error {
	if fi, err := os.Stat(dir); os.IsNotExist(err) {
		if fi != nil && !fi.IsDir() {
			dir = filepath.Dir(dir)
			return mkdir(dir)
		}
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func extractFile(r *http.Request, storePath string, override bool) (fp string, err error) {
	file, fh, err := r.FormFile(FORM_FILE)
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
	mkdir(storePath)
	fp = filepath.Join(storePath, fn)
	_, err = os.Stat(fp)
	if (override) {
		os.Remove(fp)
	} else {
		for idx := 1; err == nil; idx++ {
			ext := filepath.Ext(fn)
			fp = filepath.Join(storePath, strings.TrimRight(fn, ext) + "_" + strconv.Itoa(idx) + ext)
			_, err = os.Stat(fp)
		}
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
