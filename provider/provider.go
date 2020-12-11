package provider

import (
	"encoding/json"
	"sign-your-horse/common"
)

type Provider interface {
	//Init(alias, configJson) error
	Init(string, json.RawMessage) error
	//Run(messageCallback)
	Run(func(string, string) error)
	//Push(message)
	Push(string)
}

var providerMap = make(map[string]Provider)
var providerInstanceMap = make(map[string]Provider)

func RegisterProvider(name string, provider Provider) {
	providerMap[name] = provider
}

func GetProvider(name string) Provider {
	if provider, found := providerMap[name]; found {
		return provider
	}
	return nil
}

func GetAllProvider() ([]string, []Provider) {
	var nameList []string
	var providerList []Provider
	for name, provider := range providerMap {
		nameList = append(nameList, name)
		providerList = append(providerList, provider)
	}
	return nameList, providerList
}

func GetAllProviderInstance() ([]string, []Provider) {
	var nameList []string
	var providerList []Provider
	for name, provider := range providerInstanceMap {
		nameList = append(nameList, name)
		providerList = append(providerList, provider)
	}
	return nameList, providerList
}

func CreateProviderWithConfig(name string, alias string, config json.RawMessage) error {
	provider := GetProvider(name)
	if provider == nil {
		return common.Raise("no provider named: " + name)
	}
	providerInstance := common.CloneEmpty(provider).(Provider)
	providerInstanceMap[name+"_"+alias] = providerInstance
	return providerInstance.Init(alias, config)
}
