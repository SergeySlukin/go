package main

import (
	"hash"
	"crypto/md5"
	"io"
	"encoding/hex"
	"sync"
)

var (
	quoteCount, dupCount = 0, 0
)

type HashStorage struct {
	mu sync.Mutex
	hashMap map[string]bool
	hash hash.Hash
}


func NewHashStorage() *HashStorage {
	h := &HashStorage{
		hashMap: make(map[string]bool),
	}
	return h
}

func (h *HashStorage) SetAlgorithm(s string) {
	switch s {
	case "md5":
		h.hash = md5.New()
	}
}

func (h *HashStorage) Add(s string) ([]byte, bool)  {
	h.hash.Reset()
	io.WriteString(h.hash, s)
	localHash := h.hash.Sum(nil)
	hashString := hex.EncodeToString(localHash)
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.hashMap[hashString] {
		h.hashMap[hashString] = true
		dupCount = 0
		return localHash, true
	}
	return nil, false
}
