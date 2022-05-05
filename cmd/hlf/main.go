package main

import (
	"awesomeProject/internal/hlf/setup"
	web2 "awesomeProject/internal/hlf/web"
	"log"
	"os"
)

func main() {
	fSetup := setup.FabricSetup{
		OrdererID:       "orderer.hf.hlfexample.io",
		ChannelID:       "hlfexample",
		ChannelConfig:   "/Users/klykovandrey/my_folder/pet_projects/hlfexample/network/artifacts/hlfexample.channel.tx",
		ChainCodeID:     "example-service",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/killerter/hlfexample/internal/chaincode",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",
		UserName:        "User1",
	}

	err := fSetup.Initialize()
	if err != nil {
		log.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		log.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	web2.Serve(web2.NewController(fSetup))
}
