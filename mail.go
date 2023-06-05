package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

type URLStore interface {
	get(key string) string
	set(key, url string) (bool, error)
	count() int
}

type URLStoreMaps struct {
	urls map[string]string
	mu   sync.RWMutex
}

func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	url := r.FormValue("url")
	if url == "" {
		fmt.Fprint(w, AddForm)
		return
	}
	key, err := shortener.Put(url)
	if err != nil {
		fmt.Fprintf(w, "ERROR: %s", err)
	} else {
		fmt.Fprintf(w, "%s", key)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	fmt.Println("key: " + key)
	url := shortener.Get(key)
	if url == "" {
		http.NotFound(w, r)
		return
	}
	fmt.Println("url: " + url)

	http.Redirect(w, r, url, http.StatusFound)
}

var shortener Shortener

var (
	listenAddr = flag.String("http", ":8080", "http listen address")
	binaryFile = flag.String("bin_file", "store2.gob", "binary (gob) data store file name")
	jsonFile   = flag.String("json_file", "store.json", "json data store file name")
)

func main() {
	//store := &URLStoreMaps{urls: make(map[string]string)}
	//store := NewURLStoreGob("store.gob")
	store := NewURLStorageGobChan(*binaryFile)

	shortener = Shortener{store: store}

	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(*listenAddr, nil)
}

const AddForm = `
<html><body>
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
<\html><\body>`
