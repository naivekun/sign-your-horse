package conf

import (
	"encoding/json"
	"io/ioutil"
	"sign-your-horse/cloudscan"
	"sign-your-horse/cloudscan_client"
	"sign-your-horse/common"
)

type Config struct {
	CloudScanAPIServer json.RawMessage       `json:"cloudscan"`
	CloudScanClient    json.RawMessage       `json:"cloudscan_client"`
	Provider           []ProviderConfigBlock `json:"provider"`
	Reporter           []ReporterConfigBlock `json:"reporter"`
}

func ReadConfig(filename string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func CreateNewConfig(filename string) error {
	if common.FileExists(filename) {
		return common.Raise("file exists: " + filename)
	}

	config := &Config{}
	config.CloudScanAPIServer = cloudscan.DumpCloudScanAPIServerConfig()
	config.CloudScanClient = cloudscan_client.DumpCloudScanMessageClientConfig()
	config.Provider = DumpProviderConfigBlock()
	config.Reporter = DumpReporterConfigBlock()
	configBytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, configBytes, 0644)
}
