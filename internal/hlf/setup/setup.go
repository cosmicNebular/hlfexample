package setup

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
	"log"
)

type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	OrdererID       string
	ChannelID       string
	ChainCodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
}

func (s *FabricSetup) Initialize() error {

	if s.initialized {
		return errors.New("sdk already initialized")
	}

	sdk, err := fabsdk.New(config.FromFile(s.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	s.sdk = sdk
	log.Println("SDK created")

	rmcc := s.sdk.Context(fabsdk.WithUser(s.OrgAdmin), fabsdk.WithOrg(s.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to load Admin identity")
	}
	c, err := resmgmt.New(rmcc)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	s.admin = c
	log.Println("Ressource management client created")

	mspc, err := mspclient.New(sdk.Context(), mspclient.WithOrg(s.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}
	adminIdentity, err := mspc.GetSigningIdentity(s.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "failed to get admin signing identity")
	}
	req := resmgmt.SaveChannelRequest{ChannelID: s.ChannelID, ChannelConfigPath: s.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := s.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(s.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	log.Println("Channel created")

	if err = s.admin.JoinChannel(s.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(s.OrdererID)); err != nil {
		return errors.WithMessage(err, "failed to make admin join channel")
	}
	log.Println("Channel joined")

	log.Println("Initialization Successful")
	s.initialized = true
	return nil
}

func (s *FabricSetup) InstallAndInstantiateCC() error {

	ccPkg, err := packager.NewCCPackage(s.ChaincodePath, s.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	log.Println("ccPkg created")

	installCCReq := resmgmt.InstallCCRequest{Name: s.ChainCodeID, Path: s.ChaincodePath, Version: "0", Package: ccPkg}
	_, err = s.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}
	log.Println("Chaincode installed")

	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.hf.hlfexample.io"})

	resp, err := s.admin.InstantiateCC(s.ChannelID, resmgmt.InstantiateCCRequest{Name: s.ChainCodeID, Path: s.ChaincodeGoPath, Version: "0", Args: nil, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	log.Println("Chaincode instantiated")

	clientContext := s.sdk.ChannelContext(s.ChannelID, fabsdk.WithUser(s.UserName))
	s.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	log.Println("Channel client created")

	s.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	log.Println("Event client created")

	log.Println("Chaincode Installation & Instantiation Successful")
	return nil
}

func (s *FabricSetup) CloseSDK() {
	s.sdk.Close()
}
