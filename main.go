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
	DEFAULT_PORT = "8011"
	FORM_FILE = "file"
	STORAGE_PATH = "files"
	DATE_FORMAT = "20060102"
)

var (
	port = DEFAULT_PORT
	upload_dir = STORAGE_PATH
)

func logRequest(r *http.Request) {
	log.Printf("Request URI: %s", r.RequestURI)
}

func writeSuccess(w http.ResponseWriter, output string) {
	write(w, "[Success] " + output)
}

func writeError(w http.ResponseWriter, output string) {
	write(w, "[Error] " + output)
}

func write(w http.ResponseWriter, output string) {
	io.WriteString(w, output)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	write(w, "Welcome to HTTP File Server")
}

func upload(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	dir := strings.TrimPrefix(r.RequestURI, "/upload")
	dir = strings.TrimPrefix(dir, "/")
	byDate := r.FormValue("bydate")
	if byDate == "true" {
		dir = filepath.Join(time.Now().Format(DATE_FORMAT), dir)
	}
	storeDir := filepath.Join(upload_dir, dir)
	fp, err := extractFile(r, storeDir)
	if err != nil {
		writeError(w, err.Error())
	} else {
		writeSuccess(w, "Uploaded " + filepath.Base(fp) + "")
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	path := strings.TrimPrefix(r.RequestURI, "/remove")
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		log.Println("Path is required")
		writeError(w, "Path is required")
		return
	}
	err := os.RemoveAll(filepath.Join(upload_dir, path))
	if err != nil {
		log.Println(err)
		writeError(w, err.Error())
	} else {
		output := "Removed " + path
		log.Println(output)
		writeSuccess(w, output)
	}
}

func cmd(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	cmd := r.FormValue("cmd")
	args := r.Form["args"]
	pCmd := exec.Command(cmd, args...)
	var out bytes.Buffer
	pCmd.Stdout = &out
	err := pCmd.Run()
	if err != nil {
		writeError(w, err.Error())
	} else {
		writeSuccess(w, out.String())
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

	log.Printf("Started HTTP File Server on port: %s\n", port)
	log.Fatalln(http.ListenAndServe(":" + port, nil))
}
