package main

import "sync/atomic"

type Shortener struct {
	store URLStore
	count int64
}

func (s *Shortener) Get(key string) string {
	url, _ := s.store.get(key)
	return url
}

/*func (s *Shortener) Put(url string) string {
	for {
		key := genKey(s.store.count()) // generate the short URL
		if s.store.set(key, url) {
			return key
		}
	}
	panic("generateShortUrl - should not reach this point")
}*/

func (s *Shortener) Put(url string) (string, error) {
	for {
		key := genKey(s.count)
		ok, err := s.store.set(key, url)
		if err != nil {
			return "", err
		}
		if ok {
			return key, nil
		}
		atomic.AddInt64(&s.count, 1)
	}
	panic("shouldn't get here, shortener")
}
