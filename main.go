// Copyright 2016 carsonsx. All rights reserved.

package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"strings"
)

const (
	DEFAULT_PORT    = "8011"
	FORM_FILE       = "file"
	STORAGE_PATH    = "files"
	DATE_FORMAT     = "20060102"
)

var (
	port       = DEFAULT_PORT
	upload_dir = STORAGE_PATH
)

func welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to HTTP File Server!\n")
}

func upload(w http.ResponseWriter, r *http.Request) {
	dir := strings.TrimPrefix(r.RequestURI, "/upload")
	dir = strings.TrimPrefix(dir, "/")
	byDate := r.FormValue("bydate")
	if byDate == "true" {
		dir = filepath.Join(time.Now().Format(DATE_FORMAT), dir)
	}
	storeDir := filepath.Join(upload_dir, dir)
	fp, err := extractFile(r, storeDir)
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, "[SUCCESS] " + filepath.Join(dir, filepath.Base(fp)))
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.RequestURI, "/remove")
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		io.WriteString(w, "File or dir path is required")
		return
	}
	err := os.Remove(filepath.Join(upload_dir, path))
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, "[SUCCESS] Removed")
	}
}

func cmd(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	args := r.Form["args"]
	pCmd := exec.Command(cmd, args...)
	var out bytes.Buffer
	pCmd.Stdout = &out
	err := pCmd.Run()
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, "[SUCCESS] " + out.String())
	}
}

func main() {
	p := flag.String("p", port, "port to serve on")
	d := flag.String("d", upload_dir, "the directory of static file to host")
	flag.Parse()

	port = *p
	upload_dir = *d
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/remove/", remove)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(upload_dir))))
	http.HandleFunc("/execute", cmd)
	http.HandleFunc("/", welcome)

	log.Printf("Started HTTP File Server in port: %s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
