package main

import ("log"
	"net/http"
)

// The Root path of your static files
var ROOT_PATH = "/home/juanwolf/Documents/Devel/juanwolf.fr/"

// Setup here what language you want to use.
var LANG_DEFAULT_DIR = EN_DIR
var EN_DIR = "en/"
var FR_DIR = "fr/"
var ES_DIR = "es/"

var language_detected string

func main() {
	language_detected = LANG_DEFAULT_DIR
	http.Handle("/" , http.StripPrefix("/",
			http.FileServer(http.Dir(ROOT_PATH))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
