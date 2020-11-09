package provider

type Provider interface {
	Init() error
	Run() error
}

var providerList map[string]Provider

func RegisterProvider(name string, provider Provider) {
	providerList[name] = provider
}
