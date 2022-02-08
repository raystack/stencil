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
	logger := wrapLogger(opts.Logger)
	logger.Info(fmt.Sprintf("fetching schema from %s", url))
	data, err := downloader(url, opts.HTTPOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to fetch schema from %s", url))
		return nil, err
	}
	logger.Info(fmt.Sprintf("successfully fetched schema from %s", url))
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
	logger := wrapLogger(opts.Logger)
	return func(k cache.Key) (cache.Value, error) {
		url := k.(string)
		versionsURL := fmt.Sprintf("%s/versions", strings.TrimRight(url, "/"))
		data, err := downloader(versionsURL, opts.HTTPOptions)
		if err != nil {
			logger.Error(fmt.Sprintf("unable to download versions info, %s", err))
			return nil, err
		}
		versionsResp := &versionsModel{}
		if err := json.Unmarshal(data, versionsResp); err != nil {
			return nil, err
		}
		versions := versionsResp.Versions
		if len(versions) == 0 {
			logger.Error("no versions available for this schema")
			return nil, fmt.Errorf("no versions available")
		}
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
