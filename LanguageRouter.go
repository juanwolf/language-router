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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	// Static css files
	router.PathPrefix("/stylesheets/").Handler(http.StripPrefix("/stylesheets/",
		http.FileServer(http.Dir(ROOT_PATH + "stylesheets/"))))
	// Static js files
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir(ROOT_PATH + "js/"))))
	// Language management
	router.HandleFunc("/{lang}/", languageHandler)

	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
