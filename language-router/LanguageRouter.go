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
	PORT    		 = ":9001"
	ROOT_PATH        = "/home/juanwolf/Documents/Devel/juanwolf.fr/"
	LANG_DEFAULT 	 = "en"
	COOKIE_NAME  	 = "lang"
	COOKIE_LANG_ID   = "lang"
	NOT_FOUND_PAGE	 = "404.html"
)

// Language available
var languageMap map[string]bool

/*
 * Detect languages available on the server (all directories at ROOT_PATH with
 * the ISO-639 (2 letter only))
 */
func serverLanguageAvailable() {
	languageMap = make(map[string]bool)
	filepath.Walk(ROOT_PATH, (filepath.WalkFunc) (func(path string,
			info os.FileInfo, err error) error {

		if (info.IsDir()) {
			if (info.Name()[0] == '.' || info.Name() == "js") {
				fmt.Println("path skipped: " + path)
				return filepath.SkipDir
			}
			if len(info.Name()) <= 2 ||
				(len(info.Name()) <= 5 && strings.Contains(info.Name(), "-")) {
				languageMap[info.Name()] = true
				fmt.Println("Adding a new language: ", info.Name())
			}
		}
		return nil
	}))
}

/*
 * Read a cookie with the lang attribute.
 */
func readCookie(r *http.Request) string {
	cookie,err := r.Cookie(COOKIE_NAME);
	if (err != nil) {
		return ""
	}
	language := ""
	cookieVal := strings.Split(cookie.String(), ";");
	for i := 0; i < len(cookieVal); i++ {
		if strings.Contains(cookieVal[i], COOKIE_LANG_ID) {
			langArray := strings.Split(cookieVal[i], "=")
			language = langArray[1]
		}
	}
	fmt.Println("Language cookie detected=" + language)
	return language
}

/**
 * Read the HTTP Header and choose the preferred language if exists.
 */
func detectLanguageFromHTTPHeader(r *http.Request) string {
	header := r.Header
	languagesRequest := header.Get("Accept-Language")
	fmt.Println("Accept-Language: ", languagesRequest)
	languages := strings.Split(languagesRequest, ",")
	fmt.Println(languages)
	for _, language := range languages {
		language_without_quality := strings.Split(language, ";")[0]
		language_detected := strings.Split(language_without_quality, "-")[0]
		if languageMap[language_detected] == true {
			return language_detected
		}
	}
	return LANG_DEFAULT
}
/*
 * Detect the best language for the user (cookie first, Accept-Language
 * otherwise).
 */
func detectLanguage(r *http.Request) string {
	cookieResult := readCookie(r)
	if cookieResult != "" {
		return cookieResult
	} else {
		return detectLanguageFromHTTPHeader(r)
	}

}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r,
					ROOT_PATH + "/" + detectLanguage(r) + "/" + NOT_FOUND_PAGE)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
		language := detectLanguage(r)
		http.Redirect(w, r, r.URL.Path + language + "/", http.StatusFound)
}

func languageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	if languageMap[langAsked] {
		http.ServeFile(w, r, ROOT_PATH + "/" + langAsked + "/index.html")
	} else {
		notFoundHandler(w, r)
	}
}

func resumeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	fmt.Println("Finding the resume in " + langAsked)
	if languageMap[langAsked] {
		http.ServeFile(w, r, ROOT_PATH + "/" + langAsked + "/resume.html")
	} else {
		notFoundHandler(w, r)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	fmt.Println("Finding the chat in " + langAsked)
	if languageMap[langAsked] {
		http.ServeFile(w, r, ROOT_PATH + "/" + langAsked + "/chat.html")
	} else {
		notFoundHandler(w, r)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	langAsked := vars["lang"]
	fmt.Println("Finding the about in " + langAsked)
	if languageMap[langAsked] {
		http.ServeFile(w, r, ROOT_PATH + "/" + langAsked + "/about.html")
	} else {
		notFoundHandler(w, r)
	}
}
func main() {
	serverLanguageAvailable()
	// Router settings
	router := mux.NewRouter()
	router.Host(HOST).Schemes("http")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// Subrouter resume section
	subrouter_resume := router.Host("resume." + HOST).Subrouter()
	subrouter_resume.HandleFunc("/", rootHandler)
	subrouter_resume.HandleFunc("/{lang}/", resumeHandler)

	//Subrouter about section
	subrouter_about := router.Host("about." + HOST).Subrouter()
	subrouter_about.HandleFunc("/", rootHandler)
	subrouter_about.HandleFunc("/{lang}/", aboutHandler)

	// Subrouter chat section
	subrouter_chat := router.Host("chat." + HOST).Subrouter()
	subrouter_chat.HandleFunc("/", rootHandler)
	subrouter_chat.HandleFunc("/{lang}/", chatHandler)

	// Router section
	router.HandleFunc("/", rootHandler)
	// Static css files
	router.PathPrefix("/stylesheets/").Handler(http.StripPrefix("/stylesheets/",
		http.FileServer(http.Dir(ROOT_PATH + "stylesheets/"))))
	// Static js files
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir(ROOT_PATH + "js/"))))
	// Static image files
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/",
		http.FileServer(http.Dir(ROOT_PATH + "img/"))))
	// Language management
	router.HandleFunc("/{lang}/", languageHandler)

	http.Handle("/", router)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

