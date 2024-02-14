package response

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type CacheTime struct {
	ReleaseTime *time.Time `json:"release_time"`
}

func getReleaseTime(path, legacyCacheAPIURL string) (time.Time, int, error) {
	pathHash := md5.Sum([]byte(path))
	cacheTimeID := hex.EncodeToString(pathHash[:])
	cacheTimeResourceURL := legacyCacheAPIURL + "/v1/cache-times/" + cacheTimeID

	cacheTimeResource, statusCode, err := fetchCacheTimeResource(cacheTimeResourceURL)
	if err != nil {
		return time.Time{}, 0, err
	}

	var releaseTime time.Time
	if cacheTimeResource.ReleaseTime != nil {
		releaseTime = *cacheTimeResource.ReleaseTime
	} else {
		releaseTime = time.Time{}
	}

	return releaseTime, statusCode, nil
}

func fetchCacheTimeResource(cacheTimeResourceURL string) (CacheTime, int, error) {
	req, err := http.NewRequest(http.MethodGet, cacheTimeResourceURL, http.NoBody)
	if err != nil {
		return CacheTime{}, 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CacheTime{}, 0, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return CacheTime{}, resp.StatusCode, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CacheTime{}, 0, err
	}

	var cacheTimeResource CacheTime
	err = json.Unmarshal(body, &cacheTimeResource)
	if err != nil {
		return CacheTime{}, 0, err
	}

	return cacheTimeResource, resp.StatusCode, nil
}
