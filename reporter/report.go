package reporter

import (
	"encoding/json"
	"sign-your-horse/common"
)

type Reporter interface {
	Init(json.RawMessage) error
	Report(string) error
}

var reporterMap map[string]Reporter

func GetReporter(name string) Reporter {
	if reporter, found := reporterMap[name]; found {
		return reporter
	}
	return nil
}

func GetAllReporter() ([]string, []Reporter) {
	var nameList []string
	var reporterList []Reporter
	for name, reporter := range reporterMap {
		nameList = append(nameList, name)
		reporterList = append(reporterList, reporter)
	}
	return nameList, reporterList
}

func SetReporterConfig(name string, config json.RawMessage) error {
	reporter := GetReporter(name)
	if reporter != nil {
		return common.Raise("no reporter named: " + name)
	}
	return reporter.Init(config)
}

func RegisterReporter(name string, reporter Reporter) {
	reporterMap[name] = reporter
}
