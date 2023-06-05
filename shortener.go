package main

type Shortener struct {
	store URLStore
}

func (s *Shortener) Get(key string) string {
	return s.store.get(key)
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
		key := genKey(s.store.count())
		ok, err := s.store.set(key, url)
		if err != nil {
			return "", err
		}
		if ok {
			return key, nil
		}
	}
	panic("shouldn't get here")
}
