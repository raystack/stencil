package stencil

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goburrow/cache"
)

// RefreshStrategy clients can configure which refresh strategy to use to download latest schema.
// Default is LongPollingRefresh strategy
type RefreshStrategy int

const (
	// LongPollingRefresh this refresh strategy tries to update schema on every specified interval.
	// It doesn't check for schema changes explicitly.
	LongPollingRefresh RefreshStrategy = iota
	// VersionBasedRefresh this refresh strategy utilizes versions API provided by Stencil Server.
	// If new version is available then only schema cache would be updated.
	VersionBasedRefresh
)

func (r RefreshStrategy) getLoader(opts Options) cache.LoaderFunc {
	switch r {
	case VersionBasedRefresh:
		return versionBasedRefresh(opts)
	default:
		return longPollingRefresh(opts)
	}
}

func loadFromURL(url string, opts Options) (*Resolver, error) {
	data, err := downloader(url, opts.HTTPOptions)
	if err != nil {
		return nil, err
	}
	return NewResolver(data)
}

func longPollingRefresh(opts Options) cache.LoaderFunc {
	return func(k cache.Key) (cache.Value, error) {
		url := k.(string)
		return loadFromURL(url, opts)
	}
}

type versionsModel struct {
	Versions []int `json:"versions"`
}

func versionBasedRefresh(opts Options) cache.LoaderFunc {
	lastVersion := 0
	return func(k cache.Key) (cache.Value, error) {
		url := k.(string)
		versionsURL := fmt.Sprintf("%s/versions", strings.TrimRight(url, "/"))
		data, err := downloader(versionsURL, opts.HTTPOptions)
		if err != nil {
			return nil, err
		}
		versionsResp := &versionsModel{}
		if err := json.Unmarshal(data, versionsResp); err != nil {
			return nil, err
		}
		versions := versionsResp.Versions
		maxVersion := getMaxVersion(versions)
		if maxVersion > lastVersion {
			data, err := loadFromURL(fmt.Sprintf("%s/%d", versionsURL, maxVersion), opts)
			if err != nil {
				return nil, err
			}
			lastVersion = maxVersion
			return data, err
		}
		return nil, fmt.Errorf("schema already upto date")
	}
}

func getMaxVersion(array []int) int {
	max := array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}
