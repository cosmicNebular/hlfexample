package chaincode

import (
	"bytes"
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"reflect"
	"strconv"
	"time"
)

type RecordChaincode struct {
}

type Record struct {
	Passport     string `json:"passport"`
	Name         string `json:"name"`
	FamilyName   string `json:"family_name"`
	City         string `json:"city"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	FamilyStatus string `json:"family_status"`
}

func (r *RecordChaincode) Init(s shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (r *RecordChaincode) Invoke(s shim.ChaincodeStubInterface) peer.Response {
	f, a := s.GetFunctionAndParameters()
	switch f {
	case "create":
		return r.create(s, a)
	case "read":
		return r.read(s, a)
	case "update":
		return r.update(s, a)
	case "history":
		return r.history(s, a)
	}
	return shim.Error("Received unknown function invocation")
}

// arguments are in the struct fields order
func (r *RecordChaincode) create(s shim.ChaincodeStubInterface, a []string) peer.Response {
	if len(a) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	st, err := s.GetState(a[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if st != nil {
		return shim.Error("Person already exists, passport number: " + a[0])
	}

	rec := Record{a[0], a[1], a[2], a[3], a[4], a[5], a[6]}
	recj, err := json.Marshal(rec)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = s.PutState(a[0], recj)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = s.SetEvent("eventInvoke", nil)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 0 - passport number
func (r *RecordChaincode) read(s shim.ChaincodeStubInterface, a []string) peer.Response {
	if len(a) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	st, err := s.GetState(a[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if st == nil {
		return shim.Error("Person does not exist, passport number: " + a[0])
	}
	return shim.Success(st)
}

// 0 - passport number, 1 - field to update, 2 - new value of field
func (r *RecordChaincode) update(s shim.ChaincodeStubInterface, a []string) peer.Response {
	if len(a) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	if a[1] == "passport" {
		return shim.Error("Incorrect field to change. You can't change passport number")
	}

	st, err := s.GetState(a[0])
	if err != nil {
		return shim.Error(err.Error())
	} else if st == nil {
		return shim.Error("Person does not exist, passport number: " + a[0])
	}

	rec := Record{}
	err = json.Unmarshal(st, &rec)
	if err != nil {
		return shim.Error(err.Error())
	}

	reflect.ValueOf(&rec).Elem().FieldByName(a[1]).SetString(a[2])

	recj, err := json.Marshal(rec)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = s.PutState(a[0], recj)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = s.SetEvent("eventInvoke", nil)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (r *RecordChaincode) history(s shim.ChaincodeStubInterface, a []string) peer.Response {
	if len(a) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	iter, err := s.GetHistoryForKey(a[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	defer iter.Close()

	var buff bytes.Buffer
	buff.WriteString("[")
	first := false
	for iter.HasNext() {
		item, err := iter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if first == true {
			buff.WriteString(",")
		}
		buff.WriteString("{\"TxId\":")
		buff.WriteString("\"")
		buff.WriteString(item.TxId)
		buff.WriteString("\"")

		buff.WriteString(", \"Value\":")
		if item.IsDelete {
			buff.WriteString("null")
		} else {
			buff.WriteString(string(item.Value))
		}

		buff.WriteString(", \"Timestamp\":")
		buff.WriteString("\"")
		buff.WriteString(time.Unix(item.Timestamp.Seconds, int64(item.Timestamp.Nanos)).String())
		buff.WriteString("\"")

		buff.WriteString(", \"IsDelete\":")
		buff.WriteString("\"")
		buff.WriteString(strconv.FormatBool(item.IsDelete))
		buff.WriteString("\"")

		buff.WriteString("}")
		first = true
	}
	buff.WriteString("]")

	return shim.Success(buff.Bytes())
}
