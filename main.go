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
	"strconv"
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
	password = ""
	commands = ""
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
	if index := strings.Index(dir, "?"); index > 0 {
		dir = dir[0:index]
	}
	byDate := r.FormValue("bydate")
	if byDate == "true" {
		dir = filepath.Join(time.Now().Format(DATE_FORMAT), dir)
	}
	storeDir := filepath.Join(upload_dir, dir)
	b ,_ := strconv.ParseBool(r.FormValue("override"))
	fp, err := extractFile(r, storeDir, b)
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
		output := "the path in url is required"
		log.Println(output)
		writeError(w, output)
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

var downloadHandle = http.NotFoundHandler()

func download(w http.ResponseWriter, r *http.Request) {
	if password != "" && password != r.FormValue("password") {
		http.Error(w, "403 wrong password", http.StatusForbidden)
		return
	}
	downloadHandle.ServeHTTP(w, r)
}

func cmd(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	cmd := r.FormValue("cmd")
	if cmd == "" {
		output := "the parameter 'cmd' is required"
		log.Println(output)
		writeError(w, output)
		return
	}
	if strings.Index(commands, "," + cmd + ",") == -1 {
		output := fmt.Sprintf("server does not allow to execute '%s'", cmd)
		log.Println(output)
		writeError(w, output)
		return
	}
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
	var getVersion bool

	p := flag.String("p", port, "port to serve on")
	d := flag.String("d", upload_dir, "the directory of static file to host")
	flag.BoolVar(&showVersion, "version", false, "Print version information.")
	flag.BoolVar(&showVersion, "v", false, "Print version information.")
	flag.BoolVar(&getVersion, "getversion", false, "Get version.")
	flag.StringVar(&password, "password", "", "Set file server with simple password security mode.")
	flag.StringVar(&commands, "commands", "", "Which commands server can excuted. Add comma for multi. Do not allow any command by default.")
	flag.Parse()

	commands = strings.Replace("," + commands + ",", ",,", ",", -1)

	if showVersion {
		fmt.Printf("HFS %s\n", GetHumanVersion())
		return
	} else if getVersion {
		fmt.Printf(GetVersion())
		return
	}

	port = *p
	upload_dir = *d

	downloadHandle = http.StripPrefix("/download/", http.FileServer(http.Dir(upload_dir)))

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/remove/", remove)
	http.HandleFunc("/download/", download)
	http.HandleFunc("/execute", cmd)
	http.HandleFunc("/", welcome)
	log.Println(GitCommit)
	log.Println(GitDescribe)
	log.Printf("Started HTTP File Server on port: %s\n", port)
	log.Fatalln(http.ListenAndServe(":" + port, nil))
}
