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
	"encoding/json"
	"fmt"
)

const (
	DEFAULT_PORT = "8011"
	FORM_FILE = "file"
	STORAGE_PATH = "files"
	DATE_FORMAT = "20060102"

	STATUS_SUCCESS = 0
	STATUS_ERROR = 1
)

var (
	port = DEFAULT_PORT
	upload_dir = STORAGE_PATH
)

type ResponseData struct {
	Status  int `json:"status"`
	Message string `json:"message"`
}

func logRequest(r *http.Request) {
	log.Printf("Request URI: %s", r.RequestURI)
}

func writeSuccess(w http.ResponseWriter, message string) {
	write(w, STATUS_SUCCESS, message)
}

func writeError(w http.ResponseWriter, err string) {
	write(w, STATUS_ERROR, err)
}

func write(w http.ResponseWriter, status int, message string) {
	data := ResponseData{status, message}
	b, err := json.Marshal(data)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(b)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Welcome to HTTP File Server")
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
		writeSuccess(w, "uploaded " + filepath.Base(fp) + "")
	}
}

func remove(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	path := strings.TrimPrefix(r.RequestURI, "/remove")
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		log.Println("path is required")
		writeError(w, "path is required")
		return
	}
	err := os.RemoveAll(filepath.Join(upload_dir, path))
	if err != nil {
		log.Println(err)
		writeError(w, err.Error())
	} else {
		output := "removed " + path
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

	var showVersion bool

	p := flag.String("p", port, "port to serve on")
	d := flag.String("d", upload_dir, "the directory of static file to host")
	flag.BoolVar(&showVersion, "version", false, "Print version information.")
	flag.BoolVar(&showVersion, "v", false, "Print version information.")
	flag.Parse()

	if showVersion {
		fmt.Println(GetHumanVersion())
		return
	}

	port = *p
	upload_dir = *d
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/remove/", remove)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(upload_dir))))
	http.HandleFunc("/execute", cmd)
	http.HandleFunc("/", welcome)
	log.Println(GitCommit)
	log.Println(GitDescribe)
	log.Printf("Started HTTP File Server on port: %s\n", port)
	log.Fatalln(http.ListenAndServe(":" + port, nil))
}
