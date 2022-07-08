package cache

import (
	"fmt"
	"time"

	"github.com/cachego/cache"
)

// NewInDiskCache returns a new in-disk string cache
// with the given cachePath. eg: ".cache", "~/.cache", "/tmp/cache"
func NewInDiskCache(cachePath string) cache.Cache {
	return &inDiskCache{
		CachePath: cachePath,
		KeyPath:   fmt.Sprintf("%s/keys", cachePath),
	}
}

type inDiskCache struct {
	CachePath string // cache file path
	KeyPath   string // cache file path
}

type inDiskCacheItem struct {
	Val interface{}
	Exp int64
}

type cacheKey struct {
	keyMap map[string]bool
}

func (c *inDiskCache) Get(key string) (val interface{}, err error) {
	v, err := GetFileData(c.CachePath, key)
	if err != nil {
		return
	}
	if v == nil {
		return
	}
	data := v.(map[string]interface{})
	if int64(data["Exp"].(float64)) < time.Now().Unix() {
		c.Del(key)
		val = nil
	} else {
		val = data["Val"]
	}
	return
}

func (c *inDiskCache) Del(key string) (err error) {
	return DeleteFile(c.CachePath, key)
}

func (c *inDiskCache) Set(key string, val interface{}, ttl time.Duration) error {
	exp := int64(0)
	if ttl != 0 {
		exp = time.Now().Unix() + int64(ttl.Seconds())
	}
	cacheItem := inDiskCacheItem{
		Val: val,
		Exp: exp,
	}
	c.saveKey(key)
	return CoverFile(c.CachePath, key, cacheItem)
}

func (c *inDiskCache) IsHit(key string) (isHit bool, err error) {
	val, err := c.Get(key)
	if err != nil {
		return
	}
	return val != nil, nil
}

func (c *inDiskCache) Clear() (err error) {
	keyMaps, err := c.getAllKeyMaps()
	if err != nil || keyMaps.keyMap == nil {
		return
	}
	for key, _ := range keyMaps.keyMap {
		if _, err = c.Get(key); err != nil {
			return
		}
	}
	return
}

func (c *inDiskCache) saveKey(key string) error {
	v, err := GetFileData(c.KeyPath, "keyMap")
	if err != nil {
		return err
	}
	var keyMaps cacheKey
	if v == nil {
		keyMaps = cacheKey{keyMap: map[string]bool{key: true}}
	} else {
		data := v.(map[string]interface{})
		keyMaps = cacheKey{keyMap: data["keyMap"].(map[string]bool)}
		keyMaps.keyMap[key] = true
	}
	CoverFile(c.KeyPath, "keyMap", keyMaps)
	return nil
}

func (c *inDiskCache) getAllKeyMaps() (keys cacheKey, err error) {
	v, err := GetFileData(c.KeyPath, "keyMap")
	if err != nil {
		return
	}
	if v == nil {
		keys = cacheKey{keyMap: map[string]bool{}}
	} else {
		data := v.(map[string]interface{})
		keys = cacheKey{keyMap: data["keyMap"].(map[string]bool)}
	}
	return
}
