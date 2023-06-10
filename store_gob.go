package main

import (
	"encoding/gob"
	"io"
	"log"
	"os"
)

type URLStorageGob struct {
	URLStoreMaps
	file *os.File
}

type record struct {
	Key, URL string
}

func (s *URLStorageGob) set(key, url string) (bool, error) {
	// assume that we never can get error from Map implementation, so ignore error
	if ok, _ := s.URLStoreMaps.set(key, url); ok {
		if err := s.save(key, url); err != nil {
			// TODO probably better error processing for PROD but this depends on overall error processing strategy.
			log.Println("Error saving to URLStore:", err)
			return false, err
		}
		return true, nil
	}
	e := gob.NewEncoder(s.file)
	return true, e.Encode(record{key, url})
}

func (s *URLStorageGob) save(key, url string) error {
	e := gob.NewEncoder(s.file)
	return e.Encode(record{key, url})
}

func (s *URLStorageGob) load() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	// TODO make json OR gob encoder, similar to save to chan. OR it's causing lost of prev data on startup for JSON format.
	d := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.set(r.Key, r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func NewURLStoreGob(filename string) *URLStorageGob {
	s := &URLStorageGob{URLStoreMaps: URLStoreMaps{urls: make(map[string]string)}}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStore:", err)
	}
	s.file = f
	if err := s.load(); err != nil {
		log.Println("Error loading data in URLStore:", err)
	}
	return s
}
