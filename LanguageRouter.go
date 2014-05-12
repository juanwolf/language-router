package main

import ("log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"strings"
)

// Path to the static files
const(
	ROOT_PATH = "/home/juanwolf/Documents/Devel/juanwolf.fr/"
	EN_DIR = "en/"
	ES_DIR = "es/"
	FR_DIR = "fr/"
)

var languageMap map[string]string

var LANG_DEFAULT_DIR = EN_DIR
var language_detected  =  LANG_DEFAULT_DIR


func languageDetection() {


}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	languagesRequest := header.Get("Accept-Language")
	fmt.Println("Accept-Language: ", languagesRequest)
	languages := strings.Split(languagesRequest, ",")
	fmt.Println(languages)
	language_detected = strings.Split(languages[0], "-")[0]
	fmt.Println("Language detected", language_detected)
	language_directory := language_detected + "/"
	http.Redirect(w, r, r.URL.Path + language_directory, http.StatusFound)
}

func languageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	http.ServeFile(w, r, ROOT_PATH + langAsked)
}

func stylesheetsHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stylesheet request detected")
	vars := mux.Vars(r)
	fileWanted := vars["fileWanted"]
	http.ServeFile(w, r, ROOT_PATH + "stylesheets/" + fileWanted)
}

func fontsHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stylesheet request detected")
	vars := mux.Vars(r)
	fileWanted := vars["fileWanted"]
	http.ServeFile(w, r, ROOT_PATH + "stylesheets/fonts/" + fileWanted)
}

func jsHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stylesheet request detected")
	vars := mux.Vars(r)
	fileWanted := vars["fileWanted"]
	http.ServeFile(w, r, ROOT_PATH + "js/" + fileWanted)
	fmt.Println("File wanted ", fileWanted, "path: ", ROOT_PATH + "js/" + fileWanted)
}

func jslibHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stylesheet request detected")
	vars := mux.Vars(r)
	fileWanted := vars["fileWanted"]
	http.ServeFile(w, r, ROOT_PATH + "js/lib/" + fileWanted)
	fmt.Println("File wanted ", fileWanted, "path: ", ROOT_PATH + "js/" + fileWanted)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/stylesheets/{fileWanted}", stylesheetsHandler)
	router.HandleFunc("/stylesheets/fonts/{fileWanted}", fontsHandler)
	router.HandleFunc("/js/{fileWanted}", jsHandler)
	router.HandleFunc("/js/lib/{fileWanted}", jslibHandler)
	router.HandleFunc("/{lang}/", languageHandler)

	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
