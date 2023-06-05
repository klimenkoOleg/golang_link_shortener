package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

const saveQueueLength = 1000
const (
	FILE_FORMAT_GOB  int = 0x1
	FILE_FORMAT_JSON int = 0x2
)

type URLStorageGobChan struct {
	URLStorageGob
	saveCh     chan record
	fileFormat int // use constants like FILE_FORMAT_GOB, FILE_FORMAT_JSON
}

func (s *URLStorageGobChan) set(key, url string) (bool, error) {
	// assume that we never can get error from Map implementation, so ignore error
	if ok, _ := s.URLStoreMaps.set(key, url); ok {
		s.saveCh <- record{key, url}
		return true, nil
	}
	panic("shouldn't get here")
}

func NewURLStorageGobChan(filename string) *URLStorageGobChan {
	s := &URLStorageGobChan{
		URLStorageGob: URLStorageGob{
			URLStoreMaps: URLStoreMaps{
				urls: make(map[string]string)}},
		saveCh: make(chan record, saveQueueLength)}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("URLStore:", err)
	}
	s.file = f
	if err := s.load(); err != nil {
		log.Println("Error loading URLStore:", err)
	}
	f.Close()
	s.file = nil
	go s.saveLoop(filename)
	return s
}

type FileEncoder struct {
	encoder interface{}
}

func (fe *FileEncoder) encode(r record) error {
	switch e := fe.encoder.(type) {
	case *gob.Encoder:
		return e.Encode(r)
	case *json.Encoder:
		return e.Encode(r)
	}
	return errors.New("unknown Encoder time for file writing")
}

func NewFileEncoder(fileFormat int, f *os.File) *FileEncoder {
	var encoder interface{}
	if fileFormat == FILE_FORMAT_GOB {
		encoder = gob.NewEncoder(f)
	} else {
		encoder = json.NewEncoder(f)
	}
	return &FileEncoder{encoder: encoder}
}

func (s *URLStorageGobChan) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStore: ", err)
	}
	defer f.Close()
	encoder := NewFileEncoder(s.fileFormat, f)
	for {
		r := <-s.saveCh // taking a record from the channel and encoding it
		fmt.Println("reading from ch: " + r.URL)
		if err := encoder.encode(r); err != nil {
			log.Println("URLStore:", err)
		}
	}
}
