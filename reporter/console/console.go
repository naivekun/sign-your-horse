package console

import (
	"encoding/json"
	"sign-your-horse/common"
	"sign-your-horse/reporter"
)

type ConsoleReporter struct {
}

func (c *ConsoleReporter) Init(config json.RawMessage) error {
	return json.Unmarshal(config, c)
}

func (c *ConsoleReporter) Report(msg string) error {
	common.LogWithModule("Console Reporter", msg)
	return nil
}

func init() {
	reporter.RegisterReporter("console", &ConsoleReporter{})
}
