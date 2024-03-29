package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/rpc"
	"sync"
)

type URLStore interface {
	get(key string) (string, bool)
	set(key, url string) (bool, error)
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
	http.Redirect(w, r, url, http.StatusFound)
}

var shortener Shortener

var (
	listenAddr     = flag.String("http", ":8080", "http listen address")
	binaryFile     = flag.String("bin_file", "store2.gob", "binary (gob) data store file name")
	jsonFile       = flag.String("json_file", "store.json", "json data store file name")
	fileFormatFlag = flag.String("storage_format", "json", "Possible values, to file format storage: json | bin ")
	rpcEnabled     = flag.Bool("rpc_enabled", false, "enable RPC server")
	masterAddr     = flag.String("master", "", "RPC master address")
)

func main() {
	//store := &URLStoreMaps{urls: make(map[string]string)}
	//store := NewURLStoreGob("store.gob")
	var fileFormat int
	var actualFileName string
	if *fileFormatFlag == "json" {
		fileFormat = FILE_FORMAT_JSON
		actualFileName = *jsonFile
	} else {
		fileFormat = FILE_FORMAT_GOB
		actualFileName = *binaryFile
	}

	var store URLStore
	if *masterAddr != "" {
		store = NewProxyStore(actualFileName, fileFormat, *masterAddr)
	} else {
		store = NewURLStorageGobChan(actualFileName, fileFormat)
	}
	shortener = Shortener{store: store}

	if *rpcEnabled { // flag has been set
		rpc.RegisterName("Store", shortener)
		rpc.HandleHTTP()
	}

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
