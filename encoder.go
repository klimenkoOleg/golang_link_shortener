package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"os"
)

func (fe *FileEncoder) encode(r record) error {
	switch e := fe.encoder.(type) {
	case *gob.Encoder:
		return e.Encode(r)
	case *json.Encoder:
		return e.Encode(r)
	}
	return errors.New("unknown Encoder time for file writing")
}

func (fe *FileDecoder) decode(e any) error {
	switch e := fe.decoder.(type) {
	case *gob.Decoder:
		return e.Decode(e)
	case *json.Decoder:
		return e.Decode(e)
	}
	return errors.New("unknown Encoder time for file writing")
}

func newFileEncoder(fileFormat int, f *os.File) *FileEncoder {
	var encoder interface{}
	if fileFormat == FILE_FORMAT_GOB {
		encoder = gob.NewEncoder(f)
	} else {
		encoder = json.NewEncoder(f)
	}
	return &FileEncoder{encoder: encoder}
}

func newFileDecoder(fileFormat int, f *os.File) *FileDecoder {
	var decoder interface{}
	if fileFormat == FILE_FORMAT_GOB {
		decoder = gob.NewDecoder(f)
	} else {
		decoder = json.NewDecoder(f)
	}
	return &FileDecoder{decoder: decoder}
}
