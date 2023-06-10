package main

import "net/rpc"

type ProxyStore struct {
	client *rpc.Client
	store  *URLMuliTypeStorageChan
}

func (s *ProxyStore) get(key string) (string, bool) {
	var url string
	if val, isPresent := s.store.get(key); isPresent {
		//*url = val
		return val, true
	}
	err := s.client.Call("Store.Get", &key, &url)
	if err != nil {
		panic(err)
	}
	return url, false
}

func (s *ProxyStore) set(key, url string) (bool, error) {
	err := s.client.Call("Store.Put", &url, &key)
	if err != nil {
		return false, err
	}
	isPresent, err := s.store.set(key, url)
	return isPresent, err
}

func NewProxyStore(storageFileName string, fileFormat int, serverUrl string) *ProxyStore {
	client, err := rpc.DialHTTP("tcp", serverUrl)
	if err != nil {
		panic(err)
	}
	store := NewURLStorageGobChan(storageFileName, fileFormat)
	return &ProxyStore{client, store}
}
