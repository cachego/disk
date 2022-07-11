package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cachego/cache"
)

// NewInDiskStrCache returns a new in-disk string cache
// with the given cachePath. eg: ".cache", "~/.cache", "/tmp/cache"
func NewInDiskStrCache(cachePath string) cache.Cache {
	return &inDiskStrCache{
		CachePath: cachePath,
		KeyPath:   fmt.Sprintf("%s/keys", cachePath),
	}
}

type inDiskStrCache struct {
	CachePath string // cache file path
	KeyPath   string // cache file path
}

type inDiskStrCacheItem struct {
	Val string
	Exp int64
}

func (i *inDiskStrCacheItem) JsonEncode() ([]byte, error) {
	return json.Marshal(i)
}

func (i *inDiskStrCacheItem) JsonDecode(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), i)
}

type cacheKey struct {
	KMap map[string]bool
}

func (c *cacheKey) JsonEncode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *cacheKey) JsonDecode(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), c)
}

func (c *inDiskStrCache) Get(key string) (val interface{}, err error) {
	v, err := GetFileData(c.CachePath, key)
	if err != nil || v == nil {
		return nil, nil
	}
	item := inDiskStrCacheItem{}
	err = item.JsonDecode(string(v))
	if err != nil {
		return
	}
	if item.Exp != 0 && item.Exp < time.Now().Unix() {
		c.Del(key)
		val = nil
	} else {
		val = item.Val
	}
	return
}

func (c *inDiskStrCache) Del(key string) (err error) {
	return DeleteFile(c.CachePath, key)
}

func (c *inDiskStrCache) Set(key string, val interface{}, ttl time.Duration) error {
	if v, ok := val.(string); ok {
		exp := int64(0)
		if ttl != 0 {
			exp = time.Now().Add(ttl).Unix()
		}
		cacheItem := inDiskStrCacheItem{
			Val: v,
			Exp: exp,
		}
		err := c.saveKey(key)
		if err != nil {
			return err
		}
		data, err := cacheItem.JsonEncode()
		if err != nil {
			return err
		}
		return CoverFile(c.CachePath, key, data)
	} else {
		return errors.New("inDiskStrCache just support string values")
	}
}

func (c *inDiskStrCache) IsHit(key string) (isHit bool, err error) {
	val, err := c.Get(key)
	if err != nil {
		return
	}
	return val != nil, nil
}

func (c *inDiskStrCache) Clear() (err error) {
	keyMaps, err := c.getAllKeyMaps()
	if err != nil || keyMaps.KMap == nil {
		return nil
	}
	for key, _ := range keyMaps.KMap {
		if _, err = c.Get(key); err != nil {
			return
		}
	}
	return c.deleteKeyMap()
}

func (c *inDiskStrCache) saveKey(key string) error {
	v, _ := GetFileData(c.KeyPath, "keyMap")
	keyMaps := cacheKey{}
	if v == nil {
		keyMaps.KMap = make(map[string]bool)
	} else {
		err := keyMaps.JsonDecode(string(v))
		if err != nil {
			return err
		}
		if keyMaps.KMap == nil {
			keyMaps.KMap = make(map[string]bool)
		}
	}
	keyMaps.KMap[key] = true
	data, err := keyMaps.JsonEncode()
	if err != nil {
		return err
	}
	return CoverFile(c.KeyPath, "keyMap", data)
}

func (c *inDiskStrCache) getAllKeyMaps() (keys cacheKey, err error) {
	v, _ := GetFileData(c.KeyPath, "keyMap")
	if v == nil {
		keys.KMap = make(map[string]bool)
	} else {
		err = keys.JsonDecode(string(v))
		if err != nil {
			return
		}
	}
	return
}

func (c *inDiskStrCache) deleteKeyMap() error {
	return DeleteFile(c.KeyPath, "keyMap")
}
