package conf

import (
	"encoding/json"
)

type ReporterConfigBlock struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

func DumpReporterConfig() []ReporterConfigBlock {

}
