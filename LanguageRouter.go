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
	HOST             = "juanwolf.fr"
	PORT    		 = ":8080"
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
	for _, language := range languages {
		language_without_quality := strings.Split(language, ";")[0]
		language_detected := strings.Split(language_without_quality, "-")[0]
		fmt.Println("Language detected", language_detected)
		if languageMap[language_detected] != "" {
			language_directory := language_detected + "/"
			fmt.Println("url asked: " + r.URL.Path)
			http.Redirect(w, r, r.URL.Path+language_directory, http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, r.URL.Path + LANG_DEFAULT + "/", http.StatusFound)
}

func languageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	http.ServeFile(w, r, languageMap[langAsked])
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	fmt.Println("Finding the resume in " + languageMap[langAsked])
	http.ServeFile(w, r, languageMap[langAsked] + "/" + "resume.html")
}

func main() {
	serverLanguageAvailable()
	router := mux.NewRouter()
	router.Host(HOST).Schemes("http")
	subrouter := router.Host("resume." + HOST).Subrouter()
	subrouter.HandleFunc("/", rootHandler)
	subrouter.HandleFunc("/{lang}/", resumeHandler)
	router.HandleFunc("/", rootHandler)
	// Static css files
	router.PathPrefix("/stylesheets/").Handler(http.StripPrefix("/stylesheets/",
		http.FileServer(http.Dir(ROOT_PATH + "stylesheets/"))))
	// Static js files
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir(ROOT_PATH + "js/"))))
	// Language management
	router.HandleFunc("/{lang}/", languageHandler)
	// Subrouter for resumes

	http.Handle("/", router)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
