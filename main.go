package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	http.HandleFunc("/", index)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

	cookie := checkAndGetCookie(w, r)
	if r.Method == http.MethodPost {
		// get file uploaded
		mf, fh, err := r.FormFile("nf")
		if err != nil {
			log.Fatal(err)
		}
		defer mf.Close()

		// get file extension
		ex := strings.Split(fh.Filename, ".")[1]

		h := sha1.New()
		io.Copy(h, mf)
		newFn := fmt.Sprintf("%x", h.Sum(nil)) + "." + ex

		// wd, err := os.Getwd()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		path := filepath.Join(".", "public", "pics", newFn)
		nf, err := os.Create(path)
		if err != nil {
			log.Fatalln(err)
		}
		defer nf.Close()

		// since we have read from mf already, seek back to 0
		mf.Seek(0, 0)
		io.Copy(nf, mf)

		cookie = *appendCookieValue(&cookie, newFn)
		http.SetCookie(w, &cookie)
	}
	tpl.ExecuteTemplate(w, "index.html", strings.Split(cookie.Value, "|")[1:])
}

func checkAndGetCookie(w http.ResponseWriter, r *http.Request) http.Cookie {

	cookie, err := r.Cookie("session")

	if err != nil {
		//generate new uuid and cookie if not exists
		newId := uuid.NewV4()
		cookie = &http.Cookie{
			Name:  "session",
			Value: newId.String(),
		}
	}

	// cookie, err = setPictureNames(cookie)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	http.SetCookie(w, cookie)

	return *cookie
}

func appendCookieValue(cookie *http.Cookie, val string) *http.Cookie {
	cookieVal := cookie.Value
	if !strings.Contains(cookieVal, val) {
		// append to cookie
		cookieVal += "|" + val
	}

	cookie.Value = cookieVal
	return cookie

}

// func setPictureNames(cookie *http.Cookie) (*http.Cookie, error) {
// 	var err error
// 	pics := []string{"tea.jpeg", "coop.jpeg", "both.jpeg"}

// 	// check if cookie value includes picture names.
// 	cookieValue := cookie.Value
// 	valueSlice := strings.Split(cookie.Value, "|")

// 	var result bool
// 	for _, pic := range pics {
// 		for _, v := range valueSlice {
// 			if v == pic {
// 				result = true
// 				break
// 			}
// 		}
// 		// if the pic name is not included, then append to cookie value
// 		if !result {
// 			cookieValue = cookieValue + "|" + pic
// 		}
// 	}

// 	// set new cookie
// 	cookie = &http.Cookie{
// 		Name:  "session",
// 		Value: cookieValue,
// 	}

// 	return cookie, err
// }
