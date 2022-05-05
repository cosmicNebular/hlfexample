package setup

import (
	"awesomeProject/internal/hlf/chaincode"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"log"
	"time"
)

func (s *FabricSetup) Create(r *chaincode.Record) (fab.TransactionID, error) {
	var a [][]byte
	a = append(a, []byte(r.Passport))
	a = append(a, []byte(r.Name))
	a = append(a, []byte(r.FamilyName))
	a = append(a, []byte(r.City))
	a = append(a, []byte(r.Address))
	a = append(a, []byte(r.Phone))
	a = append(a, []byte(r.FamilyStatus))

	eventID := "eventInvoke"

	rce, ch, err := s.client.RegisterChaincodeEvent(s.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}

	response, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "create", Args: a})
	if err != nil {
		return "", err
	}

	select {
	case ccEvent := <-ch:
		log.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	s.client.UnregisterChaincodeEvent(rce)
	return response.TransactionID, nil
}

func (s *FabricSetup) Read(p string) (*chaincode.Record, error) {
	response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "create", Args: [][]byte{[]byte(p)}})
	if err != nil {
		return nil, err
	}

	r := new(chaincode.Record)
	err = json.Unmarshal(response.Payload, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *FabricSetup) Update(p string, f string, v string) (fab.TransactionID, error) {
	var a [][]byte
	a = append(a, []byte(p))
	a = append(a, []byte(f))
	a = append(a, []byte(v))

	eventID := "eventInvoke"

	rce, ch, err := s.client.RegisterChaincodeEvent(s.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}

	response, err := s.client.Execute(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "create", Args: a})
	if err != nil {
		return "", err
	}

	select {
	case ccEvent := <-ch:
		log.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	s.client.UnregisterChaincodeEvent(rce)
	return response.TransactionID, nil
}

func (s *FabricSetup) History(p string) (string, error) {
	response, err := s.client.Query(channel.Request{ChaincodeID: s.ChainCodeID, Fcn: "create", Args: [][]byte{[]byte(p)}})
	if err != nil {
		return "", err
	}
	return string(response.Payload), nil
}
