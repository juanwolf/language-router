package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"os"
	"path/filepath"
)

// Path to the static files
const (
	ROOT_PATH        = "/home/juanwolf/Documents/Devel/juanwolf.fr/"
	LANG_DEFAULT 	 = "en"
)

var languageMap map[string]string

func serverLanguageAvailable() {
	languageMap = make(map[string]string)
	filepath.Walk(ROOT_PATH, (filepath.WalkFunc)(func(path string, info os.FileInfo, err error) error {
		if (info.IsDir()) {
			if (info.Name()[0] == '.' || info.Name() == "js") {
				fmt.Println("path skipped: " + path)
				return filepath.SkipDir
			}
			if len(info.Name()) <= 2 ||
				(len(info.Name()) <= 5 && strings.Contains(info.Name(), "-"))   {
				languageMap[info.Name()] = path
				fmt.Println("Adding a new language: ", info.Name())
			}
		}
		return nil
	}))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	languagesRequest := header.Get("Accept-Language")
	fmt.Println("Accept-Language: ", languagesRequest)
	languages := strings.Split(languagesRequest, ",")
	fmt.Println(languages)
	language_detected := strings.Split(languages[0], "-")[0]
	fmt.Println("Language detected", language_detected)
	language_directory := language_detected + "/"
	http.Redirect(w, r, r.URL.Path+language_directory, http.StatusFound)
}

func languageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	http.ServeFile(w, r, ROOT_PATH+langAsked)
}

func main() {
	serverLanguageAvailable()
	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler)
	// Static css files
	router.PathPrefix("/stylesheets/").Handler(http.StripPrefix("/stylesheets/",
		http.FileServer(http.Dir(ROOT_PATH+"stylesheets/"))))
	// Static js files
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir(ROOT_PATH+"js/"))))
	// Language management
	router.HandleFunc("/{lang}/", languageHandler)

	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
