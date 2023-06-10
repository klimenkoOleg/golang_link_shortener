package main

import (
	"io"
	"log"
	"os"
)

const saveQueueLength = 1000
const (
	FILE_FORMAT_GOB  int = 0x1
	FILE_FORMAT_JSON int = 0x2
)

type FileEncoder struct {
	encoder interface{}
}

type FileDecoder struct {
	decoder interface{}
}

type URLMuliTypeStorageChan struct {
	URLStorageGob
	saveCh     chan record
	fileFormat int // use constants like FILE_FORMAT_GOB, FILE_FORMAT_JSON
}

func (s *URLMuliTypeStorageChan) set(key, url string) (bool, error) {
	// assume that we never can get error from Map implementation, so ignore error
	if ok, _ := s.URLStoreMaps.set(key, url); ok {
		if s.saveCh != nil {
			s.saveCh <- record{key, url}
		}
		return true, nil
	}
	panic("shouldn't get here, mixed storage")
}

func (s *URLMuliTypeStorageChan) save(key, url string) error {
	e := newFileEncoder(s.fileFormat, s.file)
	return e.encode(record{key, url})
}

func (s *URLMuliTypeStorageChan) load() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	d := newFileDecoder(s.fileFormat, s.file)
	var err error
	for err == nil {
		var r record
		if err = d.decode(&r); err == nil {
			s.set(r.Key, r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func NewURLStorageGobChan(filename string, fileFormat int) *URLMuliTypeStorageChan {
	s := &URLMuliTypeStorageChan{
		fileFormat: fileFormat,
		URLStorageGob: URLStorageGob{ // TODO user New factory function, otherwise it's too complex
			URLStoreMaps: URLStoreMaps{
				urls: make(map[string]string)}},
		saveCh: make(chan record, saveQueueLength)}

	if len(filename) == 0 {
		return s
	}

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

func (s *URLMuliTypeStorageChan) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening URLStore: ", err)
	}
	defer f.Close()
	encoder := newFileEncoder(s.fileFormat, f)
	for {
		r := <-s.saveCh // taking a record from the channel and encoding it
		if err := encoder.encode(r); err != nil {
			log.Println("URLStore:", err)
		}
	}
}
