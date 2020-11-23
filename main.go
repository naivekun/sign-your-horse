package main

import (
	"flag"
	"fmt"
	"log"
	"sign-your-horse/cloudscan"
	"sign-your-horse/common"
	"sign-your-horse/conf"
	"sign-your-horse/provider"
	_ "sign-your-horse/provider/chaoxing"
	_ "sign-your-horse/provider/chaoxing_cloud"
	_ "sign-your-horse/provider/teachermate_cloud"
	"sign-your-horse/reporter"
	_ "sign-your-horse/reporter/console"
	_ "sign-your-horse/reporter/wechat"
)

var configFileName string

func main() {
	fmt.Println(`
┌─┐┬┌─┐┌┐┌  ┬ ┬┌─┐┬ ┬┬─┐  ┬ ┬┌─┐┬ ┬┬─┐┌─┐┌─┐
└─┐││ ┬│││  └┬┘│ ││ │├┬┘  ├─┤│ ││ │├┬┘└─┐├┤ 
└─┘┴└─┘┘└┘   ┴ └─┘└─┘┴└─  ┴ ┴└─┘└─┘┴└─└─┘└─┘
Sign-in as a Service               @naivekun`)

	if !common.FileExists(configFileName) {
		log.Println("create default config to " + configFileName)
		common.Must(conf.CreateNewConfig(configFileName))
		return
	}
	config, err := conf.ReadConfig(configFileName)
	if err != nil {
		log.Fatalln("load config error: " + err.Error())
	}
	conf.UpdateProviderConfig(config)
	conf.UpdateReporterConfig(config)

	cloudScanServer, err := cloudscan.Init(config.CloudScanAPIServer)
	common.Must(err)
	go cloudScanServer.Run()

	_, providerList := provider.GetAllProviderInstance()
	for _, provider := range providerList {
		go provider.Run(reporter.CallReporter)
	}

	for {
		MessageData := <-cloudScanServer.APIMessageOutputChan
		for _, provider := range providerList {
			go provider.Push(MessageData)
		}
	}
}

func init() {
	flag.StringVar(&configFileName, "config", "config.json", "specify config file")
	flag.Parse()
}
