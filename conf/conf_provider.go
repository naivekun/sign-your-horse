package conf

import (
	"encoding/json"
	"sign-your-horse/common"
	"sign-your-horse/provider"
)

type ProviderConfigBlock struct {
	Name   string          `json:"name"`
	Module string          `json:"module"`
	Config json.RawMessage `json:"config"`
}

func DumpProviderConfigBlock() []ProviderConfigBlock {
	providerNameList, providerInstanceList := provider.GetAllProvider()
	var ret []ProviderConfigBlock
	for i, moduleName := range providerNameList {
		confJson, _ := json.Marshal(providerInstanceList[i])
		ret = append(ret, ProviderConfigBlock{
			Name:   moduleName + "_default",
			Module: moduleName,
			Config: confJson,
		})
	}
	return ret
}

func UpdateProviderConfig(config *Config) {
	for _, p := range config.Provider {
		err := provider.CreateProviderWithConfig(p.Module, p.Name, p.Config)
		common.Must(err)
	}
}
