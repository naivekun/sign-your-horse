package chaoxing

import (
	"sign-your-horse/provider"
)

type ChaoxingProvider struct {
	Cookie string
}

func (c *ChaoxingProvider) Init() error {

	return nil
}

func (c *ChaoxingProvider) Run() error {

	return nil
}

func init() {
	provider.RegisterProvider("chaoxing", &ChaoxingProvider{})
}
