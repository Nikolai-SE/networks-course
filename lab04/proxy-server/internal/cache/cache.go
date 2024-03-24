package cache

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"time"

	"github.com/peterbourgon/diskv/v3"
)

var (
	options = diskv.Options{
		BasePath:     "cache-storage",
		CacheSizeMax: 32 * 1024 * 1024,
	}
)

type Cache struct {
	storage *diskv.Diskv
}

func NewCache() Cache {
	return Cache{diskv.New(options)}
}

type Data struct {
	Key           string    `json:"key"`
	Value         []byte    `json:"value"`
	ExpiredTime   time.Time `json:"expired"`
	ModifiedSince string    `json:"modified_since"`
	Etag          string    `json:"etag"`
}

func (c Cache) Add(key string, value []byte, duration time.Duration, modifiedSince, etag string) error {
	data := Data{
		Key:           key,
		Value:         value,
		ExpiredTime:   time.Now().Add(duration),
		ModifiedSince: modifiedSince,
		Etag:          etag,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.storage.Write(convKey(key), jsonData)
}

func convKey(key string) string {
	h := fnv.New32a()
	h.Write([]byte(key))
	return fmt.Sprint(h.Sum32())
}

func (c Cache) Get(key string) *Data {
	keyHash := convKey(key)

	jsonData, err := c.storage.Read(keyHash)
	if err != nil {
		return nil
	}

	var data Data
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil
	}

	if time.Now().After(data.ExpiredTime) {
		c.removeImpl(keyHash)
		return nil
	}
	return &data
}

func (c Cache) removeImpl(keyHash string) {
	if err := c.storage.Erase(keyHash); err != nil {
		log.Printf("Key '%s' was not found in cache.\n", keyHash)
	}
}
func (c Cache) Remove(key string) {
	c.removeImpl(convKey(key))
}
