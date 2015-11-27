package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path"

	"runa.ninja/vecty/dom"
	"runa.ninja/vecty/dom/example/shared"
)

// TODO: cleanup, separate package
// TODO: serve GOROOT
// TODO: multiple GOPATH dirs
func sourceMiddleware(next http.Handler) http.Handler {
	gopathSrc := path.Join(os.Getenv("GOPATH"), "src")
	gopathFs := http.FileServer(http.Dir(gopathSrc))

	gorootSrc := path.Join(os.Getenv("GOROOT"), "src")
	gorootFs := http.FileServer(http.Dir(gorootSrc))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file := path.Join(gopathSrc, r.URL.Path)
		fi, err := os.Stat(file)
		if err == nil && !fi.IsDir() {
			log.Println("SRC", file)
			gopathFs.ServeHTTP(w, r)
			return
		}

		file = path.Join(gorootSrc, r.URL.Path)
		fi, err = os.Stat(file)
		if err == nil && !fi.IsDir() {
			log.Println("SRC", file)
			gorootFs.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/assets/", logMiddleware(
		http.StripPrefix("/assets/", sourceMiddleware(http.FileServer(
			http.Dir("./assets"),
		))),
	))

	html := dom.NewDOM(new(shared.Page).Render())

	http.Handle("/", sourceMiddleware(logMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.NewBuffer(nil)
		if err := html.Encode(buf, nil); err != nil {
			panic(err)
		}
		log.Println(buf.String())

		if err := html.Encode(w, nil); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))))

	log.Println("listening on http://localhost:3434/")
	log.Fatal(http.ListenAndServe(":3434", nil))
}
