package amo

import (
	"encoding/json"
	"io/ioutil"
)

type Storage interface {
	Get() *TokenData
	Set(d *TokenData) error
}

type RuntimeStorage struct {
	tokenData *TokenData
}

func (s *RuntimeStorage) Get() *TokenData {
	return s.tokenData
}

func (s *RuntimeStorage) Set(d *TokenData) error {
	s.tokenData = d
	return nil
}

type FileStorage struct {
	Path      string
	tokenData *TokenData
	loaded    bool
}

func (s *FileStorage) load() {
	s.loaded = true
	path := "amo.storage"
	if s.Path != "" {
		path = s.Path
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &s.tokenData)
}

func (s *FileStorage) flush() error {
	path := "amo.storage"
	if s.Path != "" {
		path = s.Path
	}
	b, err := json.Marshal(s.tokenData)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0644)
}

func (s *FileStorage) Get() *TokenData {
	if !s.loaded {
		s.load()
	}
	return s.tokenData
}

func (s *FileStorage) Set(d *TokenData) error {
	s.tokenData = d
	return s.flush()
}
