package steps

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cucumber/godog"
	"github.com/gorilla/mux"
)

type LegacyCacheAPIFeature struct {
	Server     *httptest.Server
	statusCode int
	db         map[string]string
}

func NewLegacyCacheAPIFeature() *LegacyCacheAPIFeature {
	f := LegacyCacheAPIFeature{
		statusCode: 0,
		db:         make(map[string]string),
	}

	router := mux.NewRouter()
	router.HandleFunc("/v1/cache-times/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		cacheTimeResource := f.db[id]

		// Set the status code if it hasn't been set already
		if f.statusCode == 0 {
			if cacheTimeResource == "" {
				f.statusCode = http.StatusNotFound
			} else {
				f.statusCode = http.StatusOK
			}
		}

		w.WriteHeader(f.statusCode)

		if _, err := w.Write([]byte(cacheTimeResource)); err != nil {
			panic(err)
		}
	}).Methods(http.MethodGet)

	f.Server = httptest.NewServer(router)

	return &f
}

func (f *LegacyCacheAPIFeature) Reset() {
	f.statusCode = 0
	f.db = make(map[string]string)
}

func (f *LegacyCacheAPIFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^the Legacy Cache API has an error$`, f.theLegacyCacheAPIHasAnError)
	ctx.Step(`^the Legacy Cache API does not have any data for the "([^"]*)" page$`, f.theLegacyCacheAPIDoesNotHaveAnyDataForThePage)
	ctx.Step(`^the "([^"]*)" page does not have a release time$`, f.thePageDoesNotHaveAReleaseTime)
	ctx.Step(`^the "([^"]*)" page will have a release in the near future$`, f.thePageWillHaveAReleaseInTheNearFuture)
	ctx.Step(`^the "([^"]*)" page will have a release in the distant future$`, f.thePageWillHaveAReleaseInTheDistantFuture)
	ctx.Step(`^the "([^"]*)" page was released long ago$`, f.thePageWasReleasedLongAgo)
	ctx.Step(`^the "([^"]*)" page was released recently$`, f.thePageWasReleasedRecently)
}

func (f *LegacyCacheAPIFeature) theLegacyCacheAPIHasAnError() error {
	f.statusCode = http.StatusInternalServerError

	return nil
}

func (f *LegacyCacheAPIFeature) theLegacyCacheAPIDoesNotHaveAnyDataForThePage(path string) error {
	delete(f.db, getID(path))

	return nil
}

func (f *LegacyCacheAPIFeature) thePageDoesNotHaveAReleaseTime(path string) error {
	f.upsertCacheTimeResource(path, time.Time{})

	return nil
}

func (f *LegacyCacheAPIFeature) thePageWillHaveAReleaseInTheNearFuture(path string) error {
	releaseTime := time.Now().Add(5 * time.Second)
	f.upsertCacheTimeResource(path, releaseTime)

	return nil
}

func (f *LegacyCacheAPIFeature) thePageWillHaveAReleaseInTheDistantFuture(path string) error {
	releaseTime := time.Now().Add(999999 * time.Hour)
	f.upsertCacheTimeResource(path, releaseTime)

	return nil
}

func (f *LegacyCacheAPIFeature) thePageWasReleasedLongAgo(path string) error {
	releaseTime, err := time.Parse(time.DateOnly, "1980-01-01")
	if err != nil {
		return err
	}

	f.upsertCacheTimeResource(path, releaseTime)

	return nil
}

func (f *LegacyCacheAPIFeature) thePageWasReleasedRecently(path string) error {
	releaseTime := time.Now().Add(-1 * time.Second)
	f.upsertCacheTimeResource(path, releaseTime)

	return nil
}

func getID(path string) string {
	pathHash := md5.Sum([]byte(path))
	return hex.EncodeToString(pathHash[:])
}

func generateCacheTimeResource(path string, releaseTime time.Time) string {
	id := getID(path)

	if releaseTime.IsZero() {
		return fmt.Sprintf(`{"_id": %q, "path": %q}`, id, path)
	}

	return fmt.Sprintf(`{"_id": %q, "path": %q, "release_time": %q}`, id, path, releaseTime.Format(time.RFC3339))
}

func (f *LegacyCacheAPIFeature) upsertCacheTimeResource(path string, releaseTime time.Time) {
	f.db[getID(path)] = generateCacheTimeResource(path, releaseTime)
}
