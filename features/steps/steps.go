package steps

import (
	"github.com/cucumber/godog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)
	c.babbageFeature.RegisterSteps(ctx)

	ctx.Step(`^the Proxy receives a GET request for "([^"]*)"$`, c.apiFeature.IGet)
}
