package awsclient

import (
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
)

// TODO: this should match exporters setting
var TAG_CACHE_DEFAULT_TIMEOUT = 60 * time.Second

type TagCache struct {
	cacheMutex *sync.Mutex
	entries    map[string]cacheEntry
	ttl        time.Duration
}

var tagCache TagCache
var singleTon sync.Once

func GetTagCache() *TagCache {
	singleTon.Do(func() {
		tagCache = TagCache{
			cacheMutex: &sync.Mutex{},
			entries:    map[string]cacheEntry{},
			ttl:        TAG_CACHE_DEFAULT_TIMEOUT,
		}
	})
	return &tagCache
}

// AddMetric adds a metric to the cache
func (mc *TagCache) addTag(arn string, tags []Tag, ttl *time.Duration) {
	mc.cacheMutex.Lock()
	mc.entries[arn] = cacheEntry{
		creation: time.Now(),
		tags:     tags,
		ttl:      ttl,
	}
	mc.cacheMutex.Unlock()
}

func (mc *TagCache) AddEc2Tags(arn string, tags []*ec2.Tag, ttl *time.Duration) {
	tagList := make([]Tag, 0)
	for _, t := range tags {
		tagList = append(tagList, Tag{Key: *t.Key, Value: *t.Value})
	}
	mc.addTag(arn, tagList, ttl)
}

func (mc *TagCache) AddRdsTags(arn string, tags []*rds.Tag, ttl *time.Duration) {
	tagList := make([]Tag, 0)
	for _, t := range tags {
		tagList = append(tagList, Tag{Key: *t.Key, Value: *t.Value})
	}
	mc.addTag(arn, tagList, ttl)
}

// GetAllMetrics Iterates over all cached metrics and discards expired ones.
func (mc *TagCache) GetTagsPerARN() map[string][]Tag {
	mc.cacheMutex.Lock()
	returnMap := make(map[string][]Tag, 0)

	for k, v := range mc.entries {
		if (v.ttl != nil && time.Since(v.creation).Seconds() > v.ttl.Seconds()) ||
			time.Since(v.creation).Seconds() > mc.ttl.Seconds() {
			delete(mc.entries, k)
		} else {
			returnMap[k] = v.tags
		}
	}
	mc.cacheMutex.Unlock()
	return returnMap
}

type cacheEntry struct {
	creation time.Time
	tags     []Tag
	ttl      *time.Duration
}

type Tag struct {
	Key   string
	Value string
}
