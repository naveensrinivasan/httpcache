// Package redis provides a redis interface for http caching.
package redis

import (
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/naveensrinivasan/httpcache"
)

// cache is an implementation of httpcache.Cache that caches responses in a
// redis server.
type cache struct {
	redis.Conn
}

// cacheKey modifies an httpcache key for use in redis. Specifically, it
// prefixes keys to avoid collision with other data stored in redis.
func cacheKey(key string) string {
	return "rediscache:" + key
}

// Get returns the response corresponding to key if present.
func (c cache) Get(key string) (resp []byte, ok bool) {
	item, err := redis.Bytes(c.Do("GET", cacheKey(key)))
	if err != nil {
		return nil, false
	}
	return item, true
}

// Set saves a response to the cache as key.
func (c cache) Set(key string, resp []byte) {
	_, e := c.Do("SET", cacheKey(key), resp)
	if e != nil {
		log.Fatal(e)
	}

}

// Delete removes the response with key from the cache.
func (c cache) Delete(key string) {
	c.Do("DEL", cacheKey(key))
}

// NewWithClient returns a new Cache with the given redis connection.
func NewWithClient(client redis.Conn) httpcache.Cache {
	return cache{client}
}

// New returns a new Transport using the redis cache implementation
func New(client redis.Conn, roundTripper http.RoundTripper) *httpcache.Transport {
	return &httpcache.Transport{Cache: cache{client}, MarkCachedResponses: true, Transport: roundTripper}
}

// NewTransport returns a new Transport with the
// provided Cache implementation and MarkCachedResponses set to true
func NewTransport(c httpcache.Cache) *httpcache.Transport {
	return &httpcache.Transport{Cache: c, MarkCachedResponses: true}
}
