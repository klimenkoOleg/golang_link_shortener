package main

func (s *URLStoreMaps) get(key string) (val string, isPresent bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, isPresent = s.urls[key]
	return
}

func (s *URLStoreMaps) set(key, url string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, present := s.urls[key]
	if present {
		return false, nil
	}
	s.urls[key] = url
	return true, nil
}

func (s *URLStoreMaps) count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func NewURLStoreMaps() *URLStoreMaps {
	return &URLStoreMaps{urls: make(map[string]string)}
}
