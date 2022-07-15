package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	port       = "80"
	UPLOAD_DIR = "upload"
)

func init() {
	nowpath, _ := os.Getwd()

	UPLOAD_DIR = filepath.Join(nowpath, UPLOAD_DIR)

	if _, err := os.Stat(UPLOAD_DIR); err != nil {
		if os.IsNotExist(err) {

			os.MkdirAll(UPLOAD_DIR, 0777)
		}
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		return
	}
	port = args[0]
}

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/upload", uploadPage)

	log.Println(": " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err.Error())
	}
}



func uploadPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, "<!DOCTYPE html PUBLIC \"-
			"<html xmlns=\"http:
			"<head>"+
			"<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />"+
			"<title></title>"+
			"</head>"+
			"<body>"+
			"<form id=\"form1\"  enctype=\"multipart/form-data\" method=\"post\" action=\"/upload\">"+
			":"+
			"<input name=\"image\" type=\"file\"  /><br/>"+
			"<input type=\"submit\" name=\"button\" id=\"button\" value=\"\" />"+
			"</form>"+
			"</body>"+
			"</html>")
		return
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileName := h.Filename
		defer f.Close()

		t, err := os.Create(filepath.Join(UPLOAD_DIR, fileName))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()

		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := filepath.Join(UPLOAD_DIR, imageId)
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)

}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir(UPLOAD_DIR)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var listHtml string
	for _, fileInfo := range fileInfoArr {
		imgid := fileInfo.Name()
		listHtml += "<li><a href=\"/view?id=" + imgid + "\">" + imgid + "</a></li>"
	}
	io.WriteString(w, "<!DOCTYPE html PUBLIC \"-
		"<html xmlns=\"http:
		"<head>"+
		"<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />"+
		"<title></title>"+
		"</head>"+
		"<body>"+
		"<a href='/upload'></a></br>"+
		"<ol>"+listHtml+"</ol>"+
		"</body>"+
		"</html>")

}
