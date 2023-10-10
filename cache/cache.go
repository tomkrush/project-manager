package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"time"
)

type FileCache struct {
	cacheDir string
	ttl      time.Duration
}

func NewFileCache(cacheDir string, ttl time.Duration) *FileCache {
	return &FileCache{
		cacheDir: cacheDir,
		ttl:      ttl,
	}
}

func (f *FileCache) generateKey(url string) string {
	hash := sha256.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))
}

func (f *FileCache) getCachePath(key string) string {
	return f.cacheDir + "/" + key
}

func (f *FileCache) Remember(key string, closure func() []byte) []byte {
	data, ok := f.Get(key)
	if ok {
		return data
	}

	data = closure()
	f.Set(key, data)

	return data
}

func (f *FileCache) Get(url string) ([]byte, bool) {
	key := f.generateKey(url)
	filePath := f.getCachePath(key)

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}
	if time.Since(info.ModTime()) > f.ttl {
		return nil, false
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, false
	}
	return data, true
}

func (f *FileCache) Set(url string, data []byte) {
	key := f.generateKey(url)

	// Ensure the cache directory exists
	if err := os.MkdirAll(f.cacheDir, 0700); err != nil {
		// Handle error, e.g., log it or return
		return
	}

	filePath := f.getCachePath(key)

	ioutil.WriteFile(filePath, data, 0644)
}
