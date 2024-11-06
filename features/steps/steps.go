package steps

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	maxAgeDirective       = `max-age`
	serverMaxAgeDirective = `s-maxage`
)

var (
	reMaxAge       = regexp.MustCompile(maxAgeDirective + `=(\d+)`)
	reServerMaxAge = regexp.MustCompile(serverMaxAgeDirective + `=(\d+)`)
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)
	c.babbageFeature.RegisterSteps(ctx)
	c.legacyCacheAPIFeature.RegisterSteps(ctx)
	c.releaseCalendarFeature.RegisterSteps(ctx)
	c.searchControllerFeature.RegisterSteps(ctx)

	ctx.Step(`^I should receive an empty response$`, c.iShouldReceiveAnEmptyResponse)
	ctx.Step(`^I should receive the same, unmodified response from Babbage$`, c.iShouldReceiveTheSameUnmodifiedResponseFromBabbage)
	ctx.Step(`^the Proxy receives a GET request for "([^"]*)"$`, c.apiFeature.IGet)
	ctx.Step(`^the Proxy receives a POST request for "([^"]*)"$`, c.apiFeature.IPostToWithBody)
	ctx.Step(`^the Proxy receives a PUT request for "([^"]*)"$`, c.apiFeature.IPut)
	ctx.Step(`^the Proxy receives a PATCH request for "([^"]*)"$`, c.apiFeature.IPatch)
	ctx.Step(`^the Proxy receives a DELETE request for "([^"]*)"$`, c.apiFeature.IDelete)
	ctx.Step(`^the (\S+) directives? should be calculated, rather than predefined$`, c.theDirectiveShouldBeCalculatedRatherThanPredefined)
	ctx.Step(`^the (\S+) directive should be (\d+)$`, c.theDirectiveShouldBe)
	ctx.Step(`^the Proxy has the publish expiry offset disabled$`, c.disablePublishExpiryOffset)
	ctx.Step(`^config includes ([A-Z0-9_]+) with a value of "([^"]*)"$`, c.configIncludes)
}

func (c *Component) disablePublishExpiryOffset() {
	c.Config.EnablePublishExpiryOffset = false
}

func (c *Component) configIncludes(configItem, configVal string) error {
	switch configItem {
	case "STALE_WHILE_REVALIDATE_SECONDS":
		seconds, err := strconv.Atoi(configVal)
		if err != nil {
			return err
		}
		c.Config.StaleWhileRevalidateSeconds = int64(seconds)
	case "ENABLE_MAX_AGE_COUNTDOWN":
		isEnabled, err := strconv.ParseBool(configVal)
		if err != nil {
			return err
		}
		c.Config.EnableMaxAgeCountdown = isEnabled
	case "ENABLE_SEARCH_CONTROLLER":
		isEnabled, err := strconv.ParseBool(configVal)
		if err != nil {
			return err
		}
		c.Config.EnableSearchController = isEnabled
	default:
		return fmt.Errorf("not a valid config item")
	}
	return nil
}

func (c *Component) iShouldReceiveAnEmptyResponse() error {
	emptyResponse := &godog.DocString{Content: ""}
	return c.apiFeature.IShouldReceiveTheFollowingResponse(emptyResponse)
}

func (c *Component) iShouldReceiveTheSameUnmodifiedResponseFromBabbage() error {
	// Ensure the body is the same
	babbageBody := &godog.DocString{Content: c.babbageFeature.Body}
	err := c.apiFeature.IShouldReceiveTheFollowingResponse(babbageBody)
	if err != nil {
		return err
	}

	// Ensure all the headers that the tester set in the mock Babbage response are present in the Proxy response
	for name, value := range c.babbageFeature.Headers {
		err = c.apiFeature.TheResponseHeaderShouldBe(name, value)
		if err != nil {
			return err
		}
	}

	// Ensure all the headers in the Proxy response are the same as the ones the tester set in the mock Babbage response
	for name, values := range c.apiFeature.HTTPResponse.Header {
		if shouldEvaluateHeader(name) {
			for _, value := range values {
				babbageHeaderValue := c.babbageFeature.Headers[name]
				errorMessage := fmt.Sprintf(`The Proxy response's %q header has a different value to the one sent by Babbage`, name)
				assert.Equal(c, babbageHeaderValue, value, errorMessage)
			}
		}
	}

	// Ensure the status code is the same
	babbageStatusCode := strconv.Itoa(c.babbageFeature.StatusCode)
	err = c.apiFeature.TheHTTPStatusCodeShouldBe(babbageStatusCode)
	if err != nil {
		return err
	}

	return c.StepError()
}

func (c *Component) getMaxAgeAndServerMaxAge() (maxAge, serverMaxAge int, err error) {
	cacheControl := c.apiFeature.HTTPResponse.Header.Get("Cache-Control")

	maxAgeMatch := reMaxAge.FindStringSubmatch(cacheControl)
	serverMaxAgeMatch := reServerMaxAge.FindStringSubmatch(cacheControl)

	maxAgeFound := assert.GreaterOrEqual(c, len(maxAgeMatch), 1)
	if !maxAgeFound {
		return 0, 0, errors.New("the max-age directive was not found or is invalid")
	}
	serverMaxAgeFound := assert.GreaterOrEqual(c, len(serverMaxAgeMatch), 1)
	if !serverMaxAgeFound {
		return 0, 0, errors.New("the s-maxage directive was not found or is invalid")
	}

	if maxAge, err = strconv.Atoi(maxAgeMatch[1]); err != nil {
		return
	}
	if serverMaxAge, err = strconv.Atoi(serverMaxAgeMatch[1]); err != nil {
		return
	}
	return
}

func (c *Component) checkMaxAgeAndServerMaxAge(checkMaxAgeCalculated, checkServerMaxAgeCalculated bool) error {
	maxAge, serverMaxAge, err := c.getMaxAgeAndServerMaxAge()
	if err != nil {
		return err
	}
	defaultCacheTime := int(c.Config.CacheTimeDefault.Seconds())
	preConfiguredCacheTimes := []int{
		defaultCacheTime,
		int(c.Config.CacheTimeErrored.Seconds()),
		int(c.Config.CacheTimeLong.Seconds()),
		int(c.Config.CacheTimeShort.Seconds()),
	}

	type testDirectives struct {
		directive            string
		age                  int
		checkAgeIsCalculated bool
	}

	for _, td := range []testDirectives{
		{
			directive:            maxAgeDirective,
			age:                  maxAge,
			checkAgeIsCalculated: checkMaxAgeCalculated,
		},
		{
			directive:            serverMaxAgeDirective,
			age:                  serverMaxAge,
			checkAgeIsCalculated: checkServerMaxAgeCalculated,
		}} {
		isAgeAPreconfiguredValue := assert.Contains(c, preConfiguredCacheTimes, td.age)
		if !isAgeAPreconfiguredValue && !td.checkAgeIsCalculated && td.age != 0 {
			return fmt.Errorf("%s is not calculated, its value (%d) is the same as one of the pre-configured cache times %+v", td.directive, td.age, preConfiguredCacheTimes)
		}

		isAgeLowerThanDefaultCacheTime := assert.Less(c, td.age, defaultCacheTime)
		if td.checkAgeIsCalculated && !isAgeLowerThanDefaultCacheTime {
			return fmt.Errorf("%s (%d) is not lower than the default cache time (%d)", td.directive, td.age, defaultCacheTime)
		} else if !td.checkAgeIsCalculated && isAgeLowerThanDefaultCacheTime && td.age > 0 {
			return fmt.Errorf("%s (%d) is lower than the default cache time (%d)", td.directive, td.age, defaultCacheTime)
		}
	}

	return nil
}

func (c *Component) theDirectiveShouldBeCalculatedRatherThanPredefined(directiveNames string) error {
	checkMaxAgeIsCalculated, checkServerMaxAgeIsCalculated := false, false
	for _, directiveName := range strings.Split(directiveNames, ",") {
		switch directiveName {
		case maxAgeDirective:
			checkMaxAgeIsCalculated = true
		case serverMaxAgeDirective:
			checkServerMaxAgeIsCalculated = true
		default:
			return fmt.Errorf("did not recognise directive %q", directiveName)
		}
	}
	return c.checkMaxAgeAndServerMaxAge(checkMaxAgeIsCalculated, checkServerMaxAgeIsCalculated)
}

func (c *Component) theDirectiveShouldBe(directiveName string, expectedValue int) error {
	maxAge, serverMaxAge, err := c.getMaxAgeAndServerMaxAge()
	if err != nil {
		return err
	}
	var obtainedValue int
	if directiveName == maxAgeDirective {
		obtainedValue = maxAge
	} else if directiveName == serverMaxAgeDirective {
		obtainedValue = serverMaxAge
	} else {
		return fmt.Errorf("did not recognise directive %q", directiveName)
	}
	if obtainedValue != expectedValue {
		return fmt.Errorf("%s (%d) does not match expected (%d)", directiveName, obtainedValue, expectedValue)
	}
	return nil
}

// shouldEvaluateHeader helps determine which headers should be skipped when comparing the Babbage and the Proxy response
//
// When writing a feature, the tester can specify which headers the mock Babbage server will return. These headers are
// saved in BabbageFeature.Headers. However, when trying to determine if the Proxy response's headers and the Babbage
// response's headers are identical, we can't just compare BabbageFeature.Headers against APIFeature.HTTPResponse.Header
// because the mock Babbage server will automatically add extra headers that may have not been defined by the tester,
// such as "Content-Length".
func shouldEvaluateHeader(headerName string) bool {
	switch headerName {
	case "Content-Length", "Content-Type", "Date":
		return false
	default:
		return true
	}
}
