# 链码

链码之间的互相调用

- 链码的安装是安装在channel内每个组织的背书节点上的, 所以需要安装多次
- 链码的实例化是在channel上进行的, 所以只需要实例化一次
- 如果链码使用了fabric-chaincode-shim的api, 则需要有init函数

## 链码开发

[https://hyperledger-fabric.readthedocs.io/en/release-2.2/developapps/smartcontract.html](https://hyperledger-fabric.readthedocs.io/en/release-2.2/developapps/smartcontract.html)

## 链码部署

[https://hyperledger-fabric.readthedocs.io/en/release-2.2/deploy_chaincode.html](https://hyperledger-fabric.readthedocs.io/en/release-2.2/deploy_chaincode.html)

## 链码编写

### 目录结构

- chaincode-go
  - chaincode
    - smartcontract.go
    - smartcontract_test.go
  - main.go
  - go.mod
  
### 链码包

https://github.com/hyperledger/fabric-contract-api-go

## 打包步骤

假设当前目录是fabric-samples/test-network, 且peer这些二进制文件放在了fabric-samples/bin目录中

- 把peer这些二进制文件加入到环境变量中, export PATH=${PWD}/../bin:$PATH
- 配置文件路径FABRIC_CFG_PATH指向config/core.yaml文件, export FABRIC_CFG_PATH=$PWD/../config/
- 执行打包命令: peer lifecycle chaincode package basic.tar.gz --path ../asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0 

## 安装链码

安装链码需要在每一个组织的背书节点上进行安装链码, 所以有几个组织, 就需要安装几次

```bash
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install basic.tar.gz

export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

peer lifecycle chaincode install basic.tar.gz
```

### 组织认可链码

- 首先要确认背书策略, Application/Channel/lifeycleEndorsement, 默认情况下要求channel上大多数成员批准链码后, 才能在channel上使用
- 由于渠道上只有两个组织, 两个的大多数是2, 所以需要让channel上的组织1和组织2都批准才可以  
- 查询链码的id: peer lifecycle chaincode queryinstalled
- 设置环境变量: export CC_PACKAGE_ID=basic_1.0:69de748301770f6ef64b42aa6bb6cb291df20aa39542c3ef94008615704007f3
- 配置组织1, 组织2的环境变量分别执行下面的命令

```bash
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

### 将链码提交给通道

- 检查通道成员已批准相同的链码定义
- `peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json`

- 需要指定哪些组织批准了链码定义
- 将channel成员的链码定义认可提交给order服务 

peer lifecycle chaincode commit -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \ 
  --channelID mychannel --name basic --version 1.0 --sequence 1 --tls \
  --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --peerAddresses localhost:7051 \ 
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
 
查询是否已经提交过去了  
peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

## 调用链码

peer chaincode invoke -o localhost:7050 \ 
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel -n basic \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"InitLedger","Args":[]}'

peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'

