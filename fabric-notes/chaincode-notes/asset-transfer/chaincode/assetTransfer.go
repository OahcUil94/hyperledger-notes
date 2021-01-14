package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"strconv"
	"time"
)

const idPrefix = "pid:skillcamps:"
const index = "size~color"

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

type AssetHistoryItem struct {
	TxId                 string               `json:"tx_id"`
	Value                *Asset               `json:"value"`
	Timestamp            *timestamp.Timestamp `json:"timestamp"`
	IsDelete             bool                 `json:"is_delete"`
}

// peer chaincode invoke -o localhost:7050 -C mychannel -n basic --peerAddresses localhost:7051 --peerAddresses localhost:9051 -c '{"function":"InitLedger","Args":["a", "b", "c", "d"]}'

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, args1 string, args2 string, args3 string, args4 string) error {
	{
		fmt.Println("获取交易id")
		fmt.Println(ctx.GetStub().GetTxID()) // 65f547a42214618fe7994084f03a9f336ad9729e69a4aba11547fb89924b5922
		fmt.Println("获取传递的参数列表")
		argList := ctx.GetStub().GetArgs() // 包含了函数名
		for k, v := range argList {
			fmt.Println("byte参数:", k, "值: ", string(v))
		}

		argList2 := ctx.GetStub().GetStringArgs()
		for k, v := range argList2 {
			fmt.Println("string类型参数:", k, "值: ", v)
		}

		funcname, paramsList := ctx.GetStub().GetFunctionAndParameters() // 将函数名和参数名区分开了
		fmt.Println(funcname)
		for k, v := range paramsList {
			fmt.Println("stringFunctionAndParameters类型参数:", k, "值: ", v)
		}

		buf, err := ctx.GetStub().GetArgsSlice()
		if err != nil {
			fmt.Println(err.Error(), ": GetArgsSlice错误")
		}

		fmt.Println("GetArgsSlice结果, ", string(buf)) // InitLedgerabcd
	}

	{
		/*
		Org1MSP�-----BEGIN CERTIFICATE-----
		MIICKTCCAdCgAwIBAgIRAMZpLg1Kz02nuxXsIus6QO8wCgYIKoZIzj0EAwIwczEL
		MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
		cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
		Lm9yZzEuZXhhbXBsZS5jb20wHhcNMjEwMTE0MDUyMzAwWhcNMzEwMTEyMDUyMzAw
		WjBrMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
		U2FuIEZyYW5jaXNjbzEOMAwGA1UECxMFYWRtaW4xHzAdBgNVBAMMFkFkbWluQG9y
		ZzEuZXhhbXBsZS5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAARyeCXisseB
		uTj/cXckoGKCudkOKVDifAYjDDDRK6n5JHAgwc8p91njfLQipfsU+jrFy045am2O
		AcnxOaQ4ll6Ko00wSzAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADArBgNV
		HSMEJDAigCB0LG7eiC7fIUyutXj4tY+9VO1TSi3Diwq3dDTiwtl0ZzAKBggqhkjO
		PQQDAgNHADBEAiBP1P0E9SSfDjqlshBtyk1LA1DnaVbkhl+UYJ3NZIjAMQIgXcJ1
		N3pl34IA4RuPRv1XUbzxHLrP5ajitLzLG5PUoeM=
		-----END CERTIFICATE-----
		 */
		buf, err := ctx.GetStub().GetCreator()
		if err != nil {
			fmt.Println("获取creator信息失败: ", err.Error())
		}
		fmt.Println("creator信息: ", string(buf))

		if err := s.getCreatorName(buf); err != nil {
			fmt.Println(err)
		}

		// 时间戳信息:  2021-01-14T08:27:28.141775063Z
		ts, err := ctx.GetStub().GetTxTimestamp()
		if err != nil {
			fmt.Println("获取时间戳:", err.Error())
		}
		fmt.Println("时间戳信息: ", ptypes.TimestampString(ts))
	}

	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, asset := range assets { // TransactionContextInterface
		asset.ID = idPrefix + asset.ID
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		// 添加数据
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}

		sizeColorIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{strconv.Itoa(asset.Size), asset.Color})
		if err != nil {
			return err
		}
		//  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the asset.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		err = ctx.GetStub().PutState(sizeColorIndexKey, value)
		if err != nil {
			return err
		}
	}

	fmt.Println("初始化账本数据")
	fmt.Println("获取channelid")
	// 获取链码部署的id
	fmt.Println(ctx.GetStub().GetChannelID())
	return nil
}

// peer chaincode query -C mychannel -n basic -c '{"Args":["GetCompositeKey"]}'
func (s *SmartContract) GetCompositeKey(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()
	// namespace:"basic" key:"\000size~color\00015\000black\000" value:"\000"
	it, err := stub.GetStateByPartialCompositeKey(index, []string{"15"})
	if err != nil {
		return err
	}

	defer it.Close()

	for it.HasNext() {
		res, err := it.Next()
		if err != nil {
			return err
		}

		// size~color  [15 white]
		obtype, parts, err := stub.SplitCompositeKey(res.Key)
		if err != nil {
			return err
		}

		fmt.Println(obtype)
		fmt.Println(parts)
		fmt.Println(res.String())
	}

	return nil
}

func (s *SmartContract) GetHistories(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	res, err := s.AssetExists(ctx, key)
	if err != nil {
		return "", err
	}

	if !res {
		return "", fmt.Errorf("没有该数据")
	}

	it, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return "", err
	}
	defer it.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for it.HasNext() {
		response, err := it.Next()
		if err != nil {
			return "", err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return buffer.String(), nil
}

//func (s *SmartContract) GetHistories2(ctx contractapi.TransactionContextInterface, key string) ([]*AssetHistoryItem, error) {
//	res, err := s.AssetExists(ctx, key)
//	if err != nil {
//		return nil, err
//	}
//
//	if !res {
//		return nil, fmt.Errorf("没有该数据")
//	}
//
//	it, err := ctx.GetStub().GetHistoryForKey(key)
//	if err != nil {
//		return nil, err
//	}
//	defer it.Close()
//	var result = make([]*AssetHistoryItem, 0)
//	for it.HasNext() {
//		response, err := it.Next()
//		if err != nil {
//			return nil, err
//		}
//
//		var item = &AssetHistoryItem{
//			TxId: response.TxId,
//			IsDelete: response.IsDelete,
//			Timestamp: response.Timestamp,
//		}
//
//		if !item.IsDelete {
//			item.Value = &Asset{}
//			_ = json.Unmarshal(response.Value, item.Value)
//		}
//
//		result = append(result, item)
//	}
//	return result, nil
//}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	res, err := ctx.GetStub().GetState(key)
	return res != nil, err
}

func (s *SmartContract) getCreatorName(info []byte) error {
	certStart := bytes.IndexAny(info, "-----BEGIN")
	if certStart == -1 {
		return fmt.Errorf("No certificate found")
	}
	certText := info[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return fmt.Errorf("Could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return fmt.Errorf("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	fmt.Println("Name:" + uname)
	return nil
}

// GetAllAssets returns all assets found in world state
// peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// namespace:"basic" key:"pid:skillcamps:asset1" value:"{\"ID\":\"pid:skillcamps:asset1\",\"appraisedValue\":300,\"color\":\"blue\",\"owner\":\"Tomoko\",\"size\":5}"
		fmt.Println(queryResponse.String())
		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAsset2"]}'
func (s *SmartContract) GetAllAsset2(ctx contractapi.TransactionContextInterface) {
	it, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer it.Close()
	var result []Asset
	for it.HasNext() {
		item, err := it.Next()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var i Asset
		_ = json.Unmarshal(item.Value, &i)
		result = append(result, i)
	}

	fmt.Println(result)
	return
}

// peer chaincode query -C mychannel -n basic -c '{"Args":["GetQueryResult"]}'
func (s *SmartContract) GetQueryResult(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()
	queryString := fmt.Sprintf(`{"selector":{"size":15,"appraisedValue":{"$gt":700}},"use_index":"colorSizeDoc"}`)
	it, err := stub.GetQueryResult(queryString)
	if err != nil {
		return fmt.Errorf("获取迭代器失败, " + err.Error())
	}

	defer it.Close()

	for it.HasNext() {
		item, err := it.Next()
		if err != nil {
			return fmt.Errorf("迭代器错误, " + err.Error())
		}

		fmt.Println(item.String())
	}

	return nil
}
