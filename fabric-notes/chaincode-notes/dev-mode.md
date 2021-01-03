# chaincode开发模式

基于`fabric-samples/test-network`现成的测试网络, 做一些修改

[TOC]

## 版本信息

- 时间: 2021-01-03
- fabric: v2.2.1
- fabric-ca: v1.4.9
- fabric-samples: 27ac653c746093dbe33eadf4c8a39b371939a805(commitid)
- couchdb: v3.1.1
- golang: v1.15.6

## 修改test-network里的内容

### Order服务关闭TLS

修改`docker/docker-compose-test-net.yaml`文件: 

```
- ORDERER_GENERAL_TLS_ENABLED=true
+ ORDERER_GENERAL_TLS_ENABLED=false
```

由于测试网络Order服务处于`etcdraft`模式下, 是无法关闭TLS的, 所以需要调成为`solo`模式, 修改`configtx/configtx.yaml`文件:

```
  Orderer: &OrdererDefaults
 
      # Orderer Type: The orderer implementation to start
-     OrdererType: etcdraft
+     OrdererType: solo
```

### Peer服务关闭TLS

修改`docker/docker-compose-test-net.yaml`文件, 修改以下几点内容:

1. 开启链码级Debug日志
2. 关闭TLS
3. peer node 启动命令增加`--peer-chaincodedev`参数
4. 暴露链码监听端口

```
  peer0.org1.example.com:
      environment:
        ...
-       - FABRIC_LOGGING_SPEC=INFO
+       - FABRIC_LOGGING_SPEC=chaincode=debug
        #- FABRIC_LOGGING_SPEC=DEBUG
-       - CORE_PEER_TLS_ENABLED=true
+       - CORE_PEER_TLS_ENABLED=false
        - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
        - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
        ...
-     command: peer node start
+     command: peer node start --peer-chaincodedev
      ports:
        - 7051:7051
+       - 7052:7052
 
  peer0.org2.example.com:
      environment:
        ...
-       - FABRIC_LOGGING_SPEC=INFO
+       - FABRIC_LOGGING_SPEC=chaincode=debug
        #- FABRIC_LOGGING_SPEC=DEBUG
-       - CORE_PEER_TLS_ENABLED=true
+       - CORE_PEER_TLS_ENABLED=false
        - CORE_PEER_LISTENADDRESS=0.0.0.0:9051
        - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:9052
        ...
-     command: peer node start
+     command: peer node start --peer-chaincodedev
      ports:
        - 9051:9051
+       - 9052:9052
```

### 修改环境变量脚本文件scripts/envVar.sh

```
  source scriptUtils.sh

- export CORE_PEER_TLS_ENABLED=true
+ export CORE_PEER_TLS_ENABLED=false

  parsePeerConnectionParameters() {
      ...
      PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"
      ## Set path to TLS certificate
      TLSINFO=$(eval echo "--tlsRootCertFiles \$PEER0_ORG$1_CA")
-     PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
+     # PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
      ...
  }
```

### 修改创建通道脚本文件scripts/createChannel

```
  createChannel() {
          ...
-		  peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.example.com -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock ./channel-artifacts/${CHANNEL_NAME}.block --tls --cafile $ORDERER_CA >&log.txt
+		  # peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.example.com -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock ./channel-artifacts/${CHANNEL_NAME}.block --tls --cafile $ORDERER_CA >&log.txt
+		  peer channel create -o localhost:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock ./channel-artifacts/${CHANNEL_NAME}.block >&log.txt
		  ...
  }


  updateAnchorPeers() {
          ...
-		  peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA >&log.txt
+		  # peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA >&log.txt
+		  peer channel update -o localhost:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
          ...
  }
```

### 修改部署链码脚本文件scripts/deployCC.sh

```
  # approveForMyOrg VERSION PEER ORG
  approveForMyOrg() {
    ...
-   peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
+   # peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
+   peer lifecycle chaincode approveformyorg -o localhost:7050 --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${CC_NAME}:${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
    ...
  }

  commitChaincodeDefinition() {
    ...
-   peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
+   # peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
+   peer lifecycle chaincode commit -o localhost:7050 --channelID $CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
    ...
  }

  chaincodeInvokeInit() {
    ...
-   peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call} >&log.txt
+   # peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call} >&log.txt
+   peer chaincode invoke -o localhost:7050 -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call} >&log.txt
    ...
  }

  ...
  ## package the chaincode
- packageChaincode
+ # packageChaincode

  ## Install chaincode on peer0.org1 and peer0.org2
  # infoln "Installing chaincode on peer0.org1..."
- installChaincode 1
+ # installChaincode 1
  # infoln "Install chaincode on peer0.org2..."
- installChaincode 2
+ # installChaincode 2

  ## query whether the chaincode is installed
- queryInstalled 1
+ # queryInstalled 1

  ## approve the definition for org1
  approveForMyOrg 1
  ...
```

## 启动测试网络

`./network.sh up createChannel -c mychannel`

## 注册链码到peer中

这里使用`fabric-samples/asset-transfer-basic/chaincode-go`中的链码, 和以往不同的是链码无需打包, 可在本地直接运行:

```
export CORE_CHAINCODE_LOGLEVEL=debug 
export CORE_PEER_TLS_ENABLED=false
# 链码名+版本号  
export CORE_CHAINCODE_ID_NAME=basic:v1.0.0 

# 需在两个组织的peer上都进行注册
go run . -peer.address 127.0.0.1:7052
go run . -peer.address 127.0.0.1:9052
```

> 注意: 这里指定的peer.address是peer服务链码的监听地址

## 部署链码

```
cd fabric-samples/test-network
./network.sh deployCC -ccn basic -ccp . -ccv v1.0.0 -ccs 1 -ccl go
```

## 调用链码测试

```
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

export CORE_PEER_TLS_ENABLED=false
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer chaincode invoke -o localhost:7050 -C mychannel -n basic --peerAddresses localhost:7051 --peerAddresses localhost:9051 -c '{"function":"InitLedger","Args":[]}'
peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
peer chaincode invoke -o localhost:7050 -C mychannel -n basic --peerAddresses localhost:7051 --peerAddresses localhost:9051 -c '{"function":"TransferAsset","Args":["asset6","Christopher"]}'
peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","asset6"]}'
``` 

## 参考资料

> 注意: chaincode开发模式的文档在1.4和2.3版本可以查看到, 中间其它版本无法找到 
 
- [Tutorials » Running chaincode in development mode](https://hyperledger-fabric.readthedocs.io/en/release-2.3/peer-chaincode-devmode.html)
- [Commands Reference » peer node](https://hyperledger-fabric.readthedocs.io/en/release-2.2/commands/peernode.html)
