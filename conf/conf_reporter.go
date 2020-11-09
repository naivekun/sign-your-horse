package conf

import (
	"encoding/json"
	"sign-your-horse/common"
	"sign-your-horse/reporter"
)

type ReporterConfigBlock struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

func DumpReporterConfigBlock() []ReporterConfigBlock {
	reporterNameList, reporterInstanceList := reporter.GetAllReporter()
	var ret []ReporterConfigBlock
	for i, name := range reporterNameList {
		confJson, _ := json.Marshal(reporterInstanceList[i])
		ret = append(ret, ReporterConfigBlock{
			Name:   name,
			Config: confJson,
		})
	}
	return ret
}

func UpdateReporterConfig(config *Config) {
	for _, c := range config.Reporter {
		err := reporter.SetReporterConfig(c.Name, c.Config)
		common.Must(err)
	}
}
