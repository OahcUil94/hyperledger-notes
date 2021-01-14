
./network.sh deployCC -ccn basic -ccl go -ccv v1.0 -ccs 1 -ccp ~/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/chaincode-notes/asset-transfer

go run . -peer.address 127.0.0.1:7052
go run . -peer.address 127.0.0.1:9052

peer chaincode invoke -o localhost:7050 -C mychannel -n basic --peerAddresses localhost:7051 --peerAddresses localhost:9051 -c '{"function":"InitLedger","Args":[]}'

## 链码方法

`ctx.GetStub().PutState(asset.ID, assetJSON)`: 添加数据
`ctx.GetStub().GetChannelID()`: 获取链码部署的channel名 


peer chaincode query -C mychannel -n basic -c '{"Args":["GetHistories", "pid:skillcamps:asset1"]}'

```yaml
history:
    # enableHistoryDatabase - options are true or false
    # Indicates if the history of key updates should be stored.
    # All history 'index' will be stored in goleveldb, regardless if using
    # CouchDB or alternate database for the state.
    enableHistoryDatabase: true
```

1. components.schemas..required: Array must have at least 1 items
panic: Error creating asset-transfer-basic chaincode: Cannot use metadata. Metadata did not match schema:
1. components.schemas..required: Array must have at least 1 items

出现类似于上面的错误表示合约的函数返回的对象，数字类型只能是链上的数据类型, 它会做校验

参数不允许是可变参数

```
type ChaincodeStubInterface interface {
	GetArgs() [][]byte √
	GetStringArgs() []string √
	GetFunctionAndParameters() (string, []string) √
	GetArgsSlice() ([]byte, error) √
	GetTxID() string √
	GetChannelID() string √
	GetState(key string) ([]byte, error) √
	PutState(key string, value []byte) error √
	DelState(key string) error √
    CreateCompositeKey(objectType string, attributes []string) (string, error) √
    GetCreator() ([]byte, error) √
    GetTxTimestamp() (*timestamp.Timestamp, error) √
    GetHistoryForKey(key string) (HistoryQueryIteratorInterface, error) √
    GetStateByPartialCompositeKey(objectType string, keys []string) (StateQueryIteratorInterface, error) √
    SplitCompositeKey(compositeKey string) (string, []string, error) √
    GetStateByRange(startKey, endKey string) (StateQueryIteratorInterface, error) √
    GetQueryResult(query string) (StateQueryIteratorInterface, error) √

	GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
		bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
    GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
		pageSize int32, bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
	GetQueryResultWithPagination(query string, pageSize int32,
		bookmark string) (StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
	
    InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response
	SetStateValidationParameter(key string, ep []byte) error
	GetStateValidationParameter(key string) ([]byte, error)
	
	GetPrivateData(collection, key string) ([]byte, error)
	GetPrivateDataHash(collection, key string) ([]byte, error)
	PutPrivateData(collection string, key string, value []byte) error
	DelPrivateData(collection, key string) error
	GetPrivateDataByRange(collection, startKey, endKey string) (StateQueryIteratorInterface, error)
	GetPrivateDataByPartialCompositeKey(collection, objectType string, keys []string) (StateQueryIteratorInterface, error)
	GetPrivateDataQueryResult(collection, query string) (StateQueryIteratorInterface, error)

    SetPrivateDataValidationParameter(collection, key string, ep []byte) error
	GetPrivateDataValidationParameter(collection, key string) ([]byte, error)
	
	GetTransient() (map[string][]byte, error)
	GetBinding() ([]byte, error)
	GetDecorations() map[string][]byte
	GetSignedProposal() (*pb.SignedProposal, error)
	
	SetEvent(name string, payload []byte) error
}
```

https://docs.couchdb.org/en/stable/api/database/find.html
