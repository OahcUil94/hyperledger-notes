# fabric环境搭建 

时间：2020-12-7

目前官方最新稳定版本是v2.2.1

## 初始环境搭建
 
官方提供了一个便捷的脚本文件
https://raw.githubusercontent.com/hyperledger/fabric/v2.2.1/scripts/bootstrap.sh

curl -# -O https://cdn.jsdelivr.net/gh/hyperledger/fabric@v2.2.1/scripts/bootstrap.sh

由于国内的网络环境, 先把文件下载下来, 然后对其中的内容进行修改

1. cloneSamplesRepo函数中, 修改克隆sample项目的地址`git clone -b master https://github.com.cnpmjs.org/hyperledger/fabric-samples.git && cd fabric-samples && git checkout v${VERSION}`
2. pullBinaries函数中, 修改二进制文件的下载地址 `download "${BINARY_FILE}" "https://gh.api.99988866.xyz/https://github.com/hyperledger/fabric/releases/download/v${VERSION}/${BINARY_FILE}"`
`download "${CA_BINARY_FILE}" "https://gh.api.99988866.xyz/https://github.com/hyperledger/fabric-ca/releases/download/v${CA_VERSION}/${CA_BINARY_FILE}"`

脚本所执行的操作:

1. 克隆sample项目, 并进入到samples目录中
2. 下载二进制文件
3. 下载docker镜像

通过指定参数来绕开(bypass)某一步操作, -d绕过docker镜像下载, -s绕过samples仓库的clone, -b绕过二进制文件下载

二进制文件下载下来之后, 将二进制文件路径添加到环境变量中

fabric-samples实验的commit id: 36b5788bad3d43884e557e11f0db4e60660616d2 

FABRIC_CFG_PATH=$PWD

## 搭建测试网络

curl -# -O https://gh.api.99988866.xyz/https://github.com/hyperledger/fabric/releases/download/v2.3.0/hyperledger-fabric-linux-amd64-2.3.0.tar.gz

在上面配置了二进制文件环境变量之后, 执行`./network.sh up`报错: 

```
Starting nodes with CLI timeout of '5' tries and CLI delay of '3' seconds and using database 'leveldb' with crypto from 'cryptogen'
Peer binary and configuration files not found..

Follow the instructions in the Fabric docs to install the Fabric Binaries:
https://hyperledger-fabric.readthedocs.io/en/latest/install.html 
```

peer节点是可以访问的, 就是缺少了config文件夹, 原因是`v2.2.1`版本的二进制压缩包中, 没有config目录, 而在其他版本例如`2.3.0`版本中是有的, 所以需要下载其他版本的压缩包, 把其中的config目录复制出来, 放到bin同级的目录中

现在执行就没有问题了:

```bash
[vagrant@localhost test-network]$ ./network.sh up
Starting nodes with CLI timeout of '5' tries and CLI delay of '3' seconds and using database 'leveldb' with crypto from 'cryptogen'
LOCAL_VERSION=2.2.1
DOCKER_IMAGE_VERSION=2.2.1
/home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/bin/cryptogen
Generate certificates using cryptogen tool
Create Org1 Identities
+ cryptogen generate --config=./organizations/cryptogen/crypto-config-org1.yaml --output=organizations
org1.example.com
+ res=0
Create Org2 Identities
+ cryptogen generate --config=./organizations/cryptogen/crypto-config-org2.yaml --output=organizations
org2.example.com
+ res=0
Create Orderer Org Identities
+ cryptogen generate --config=./organizations/cryptogen/crypto-config-orderer.yaml --output=organizations
+ res=0
Generate CCP files for Org1 and Org2
/home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/bin/configtxgen
Generating Orderer Genesis block
+ configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block
2020-12-07 18:09:36.047 CST [common.tools.configtxgen] main -> INFO 001 Loading configuration
2020-12-07 18:09:36.091 CST [common.tools.configtxgen.localconfig] completeInitialization -> INFO 002 orderer type: etcdraft
2020-12-07 18:09:36.092 CST [common.tools.configtxgen.localconfig] completeInitialization -> INFO 003 Orderer.EtcdRaft.Options unset, setting to tick_interval:"500ms" election_tick:10 heartbeat_tick:1 max_inflight_blocks:5 snapshot_interval_size:16777216
2020-12-07 18:09:36.092 CST [common.tools.configtxgen.localconfig] Load -> INFO 004 Loaded configuration: /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/configtx/configtx.yaml
2020-12-07 18:09:36.108 CST [common.tools.configtxgen] doOutputBlock -> INFO 005 Generating genesis block
2020-12-07 18:09:36.109 CST [common.tools.configtxgen] doOutputBlock -> INFO 006 Writing genesis block
+ res=0
Creating network "net_test" with the default driver
Creating volume "net_orderer.example.com" with default driver
Creating volume "net_peer0.org1.example.com" with default driver
Creating volume "net_peer0.org2.example.com" with default driver
Creating peer0.org2.example.com ... done
Creating orderer.example.com    ... done
Creating peer0.org1.example.com ... done
CONTAINER ID        IMAGE                               COMMAND             CREATED             STATUS                  PORTS                              NAMES
2a0e28c2ffc9        hyperledger/fabric-peer:latest      "peer node start"   1 second ago        Up Less than a second   0.0.0.0:7051->7051/tcp             peer0.org1.example.com
56865b61c236        hyperledger/fabric-orderer:latest   "orderer"           1 second ago        Up Less than a second   0.0.0.0:7050->7050/tcp             orderer.example.com
0bcd4ab311ca        hyperledger/fabric-peer:latest      "peer node start"   1 second ago        Up Less than a second   7051/tcp, 0.0.0.0:9051->9051/tcp   peer0.org2.example.com
```

创建channel: 

```bash
[vagrant@localhost test-network]$ ./network.sh createChannel
Creating channel 'mychannel'.
If network is not up, starting nodes with CLI timeout of '5' tries and CLI delay of '3' seconds and using database 'leveldb
Generating channel create transaction 'mychannel.tx'
+ configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID mychannel
2020-12-07 18:12:36.015 CST [common.tools.configtxgen] main -> INFO 001 Loading configuration
2020-12-07 18:12:36.033 CST [common.tools.configtxgen.localconfig] Load -> INFO 002 Loaded configuration: /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/configtx/configtx.yaml
2020-12-07 18:12:36.033 CST [common.tools.configtxgen] doOutputChannelCreateTx -> INFO 003 Generating new channel configtx
2020-12-07 18:12:36.048 CST [common.tools.configtxgen] doOutputChannelCreateTx -> INFO 004 Writing new channel tx
+ res=0
Generating anchor peer update transactions
Generating anchor peer update transaction for Org1MSP
+ configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
2020-12-07 18:12:36.075 CST [common.tools.configtxgen] main -> INFO 001 Loading configuration
2020-12-07 18:12:36.094 CST [common.tools.configtxgen.localconfig] Load -> INFO 002 Loaded configuration: /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/configtx/configtx.yaml
2020-12-07 18:12:36.094 CST [common.tools.configtxgen] doOutputAnchorPeersUpdate -> INFO 003 Generating anchor peer update
2020-12-07 18:12:36.100 CST [common.tools.configtxgen] doOutputAnchorPeersUpdate -> INFO 004 Writing anchor peer update
+ res=0
Generating anchor peer update transaction for Org2MSP
+ configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
2020-12-07 18:12:36.125 CST [common.tools.configtxgen] main -> INFO 001 Loading configuration
2020-12-07 18:12:36.142 CST [common.tools.configtxgen.localconfig] Load -> INFO 002 Loaded configuration: /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/configtx/configtx.yaml
2020-12-07 18:12:36.142 CST [common.tools.configtxgen] doOutputAnchorPeersUpdate -> INFO 003 Generating anchor peer update
2020-12-07 18:12:36.148 CST [common.tools.configtxgen] doOutputAnchorPeersUpdate -> INFO 004 Writing anchor peer update
+ res=0
Creating channel mychannel
Using organization 1
+ peer channel create -o localhost:7050 -c mychannel --ordererTLSHostnameOverride orderer.example.com -f ./channel-artifacts/mychannel.tx --outputBlock ./channel-artifacts/mychannel.block --tls --cafile /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
+ res=0
2020-12-07 18:12:39.223 CST [channelCmd] InitCmdFactory -> INFO 001 Endorser and orderer connections initialized
2020-12-07 18:12:39.254 CST [cli.common] readBlock -> INFO 002 Expect block, but got status: &{SERVICE_UNAVAILABLE}
2020-12-07 18:12:39.258 CST [channelCmd] InitCmdFactory -> INFO 003 Endorser and orderer connections initialized
2020-12-07 18:12:39.459 CST [cli.common] readBlock -> INFO 004 Expect block, but got status: &{SERVICE_UNAVAILABLE}
2020-12-07 18:12:39.464 CST [channelCmd] InitCmdFactory -> INFO 005 Endorser and orderer connections initialized
2020-12-07 18:12:39.665 CST [cli.common] readBlock -> INFO 006 Expect block, but got status: &{SERVICE_UNAVAILABLE}
2020-12-07 18:12:39.672 CST [channelCmd] InitCmdFactory -> INFO 007 Endorser and orderer connections initialized
2020-12-07 18:12:39.874 CST [cli.common] readBlock -> INFO 008 Expect block, but got status: &{SERVICE_UNAVAILABLE}
2020-12-07 18:12:39.878 CST [channelCmd] InitCmdFactory -> INFO 009 Endorser and orderer connections initialized
2020-12-07 18:12:40.079 CST [cli.common] readBlock -> INFO 00a Expect block, but got status: &{SERVICE_UNAVAILABLE}
2020-12-07 18:12:40.084 CST [channelCmd] InitCmdFactory -> INFO 00b Endorser and orderer connections initialized
2020-12-07 18:12:40.287 CST [cli.common] readBlock -> INFO 00c Received block: 0
Channel 'mychannel' created
Join Org1 peers to the channel...
Using organization 1
+ peer channel join -b ./channel-artifacts/mychannel.block
+ res=0
2020-12-07 18:12:43.450 CST [channelCmd] InitCmdFactory -> INFO 001 Endorser and orderer connections initialized
2020-12-07 18:12:43.472 CST [channelCmd] executeJoin -> INFO 002 Successfully submitted proposal to join channel
Join Org2 peers to the channel...
Using organization 2
+ peer channel join -b ./channel-artifacts/mychannel.block
+ res=0
2020-12-07 18:12:46.544 CST [channelCmd] InitCmdFactory -> INFO 001 Endorser and orderer connections initialized
2020-12-07 18:12:46.565 CST [channelCmd] executeJoin -> INFO 002 Successfully submitted proposal to join channel
Updating anchor peers for org1...
Using organization 1
+ peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
+ res=0
2020-12-07 18:12:49.693 CST [channelCmd] InitCmdFactory -> INFO 001 Endorser and orderer connections initialized
2020-12-07 18:12:49.707 CST [channelCmd] update -> INFO 002 Successfully submitted channel update
Anchor peers updated for org 'Org1MSP' on channel 'mychannel'
Updating anchor peers for org2...
Using organization 2
+ peer channel update -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com -c mychannel -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
+ res=0
2020-12-07 18:12:55.768 CST [channelCmd] InitCmdFactory -> INFO 001 Endorser and orderer connections initialized
2020-12-07 18:12:55.784 CST [channelCmd] update -> INFO 002 Successfully submitted channel update
Anchor peers updated for org 'Org2MSP' on channel 'mychannel'
Channel successfully joined
```

## 部署链码

```bash
[vagrant@localhost test-network]$ ./network.sh deployCC
deploying chaincode on channel 'mychannel'
executing with the following
- CHANNEL_NAME: mychannel
- CC_NAME: basic
- CC_SRC_PATH: NA
- CC_SRC_LANGUAGE: go
- CC_VERSION: 1.0
- CC_SEQUENCE: 1
- CC_END_POLICY: NA
- CC_COLL_CONFIG: NA
- CC_INIT_FCN: NA
- DELAY: 3
- MAX_RETRY: 5
- VERBOSE: false
Determining the path to the chaincode
asset-transfer-basic
Vendoring Go dependencies at ../asset-transfer-basic/chaincode-go/
~/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/asset-transfer-basic/chaincode-go ~/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network
~/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network
Finished vendoring Go dependencies
+ peer lifecycle chaincode package basic.tar.gz --path ../asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0
+ res=0
Chaincode is packaged
Installing chaincode on peer0.org1...
Using organization 1
+ peer lifecycle chaincode install basic.tar.gz
+ res=0
2020-12-07 18:16:04.661 CST [cli.lifecycle.chaincode] submitInstallProposal -> INFO 001 Installed remotely: response:<status:200 payload:"\nJbasic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad\022\tbasic_1.0" >
2020-12-07 18:16:04.662 CST [cli.lifecycle.chaincode] submitInstallProposal -> INFO 002 Chaincode code package identifier: basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad
Chaincode is installed on peer0.org1
Install chaincode on peer0.org2...
Using organization 2
+ peer lifecycle chaincode install basic.tar.gz
+ res=0
2020-12-07 18:16:14.694 CST [cli.lifecycle.chaincode] submitInstallProposal -> INFO 001 Installed remotely: response:<status:200 payload:"\nJbasic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad\022\tbasic_1.0" >
2020-12-07 18:16:14.695 CST [cli.lifecycle.chaincode] submitInstallProposal -> INFO 002 Chaincode code package identifier: basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad
Chaincode is installed on peer0.org2
Using organization 1
+ peer lifecycle chaincode queryinstalled
+ res=0
Installed chaincodes on peer:
Package ID: basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad, Label: basic_1.0
Query installed successful on peer0.org1 on channel
Using organization 1
+ peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name basic --version 1.0 --package-id basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad --sequence 1
+ res=0
2020-12-07 18:16:16.885 CST [chaincodeCmd] ClientWait -> INFO 001 txid [097ca06935f7dfe1132b38ee95b6d1654510e41a674b9d27d9473803048ad547] committed with status (VALID) at localhost:7051
Chaincode definition approved on peer0.org1 on channel 'mychannel'
Using organization 1
Checking the commit readiness of the chaincode definition on peer0.org1 on channel 'mychannel'...
Attempting to check the commit readiness of the chaincode definition on peer0.org1, Retry after 3 seconds.
+ peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --output json
+ res=0

        "approvals": {
                "Org1MSP": true,
                "Org2MSP": false
        }
}
Checking the commit readiness of the chaincode definition successful on peer0.org1 on channel 'mychannel'
Using organization 2
Checking the commit readiness of the chaincode definition on peer0.org2 on channel 'mychannel'...
Attempting to check the commit readiness of the chaincode definition on peer0.org2, Retry after 3 seconds.
+ peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --output json
+ res=0

        "approvals": {
                "Org1MSP": true,
                "Org2MSP": false
        }
}
Checking the commit readiness of the chaincode definition successful on peer0.org2 on channel 'mychannel'
Using organization 2
+ peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name basic --version 1.0 --package-id basic_1.0:4ec191e793b27e953ff2ede5a8bcc63152cecb1e4c3f301a26e22692c61967ad --sequence 1
+ res=0
2020-12-07 18:16:25.097 CST [chaincodeCmd] ClientWait -> INFO 001 txid [b75ac55ae365ada5116b7af518a4c90f23f6c27b73ead44ee646fc6a35595dea] committed with status (VALID) at localhost:9051
Chaincode definition approved on peer0.org2 on channel 'mychannel'
Using organization 1
Checking the commit readiness of the chaincode definition on peer0.org1 on channel 'mychannel'...
Attempting to check the commit readiness of the chaincode definition on peer0.org1, Retry after 3 seconds.
+ peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --output json
+ res=0

        "approvals": {
                "Org1MSP": true,
                "Org2MSP": true
        }
}
Checking the commit readiness of the chaincode definition successful on peer0.org1 on channel 'mychannel'
Using organization 2
Checking the commit readiness of the chaincode definition on peer0.org2 on channel 'mychannel'...
Attempting to check the commit readiness of the chaincode definition on peer0.org2, Retry after 3 seconds.
+ peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --output json
+ res=0

        "approvals": {
                "Org1MSP": true,
                "Org2MSP": true
        }
}
Checking the commit readiness of the chaincode definition successful on peer0.org2 on channel 'mychannel'
Using organization 1
Using organization 2
+ peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name basic --peerAddresses localhost:7051 --tlsRootCertFiles /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/vagrant/code/go/src/github.com/OahcUil94/hyperledger-notes/fabric-notes/fabric-samples/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --version 1.0 --sequence 1
+ res=0
2020-12-07 18:16:33.420 CST [chaincodeCmd] ClientWait -> INFO 001 txid [626210823222a9e17401c568c2371eb38b8896e993ad74d43673b66c94465fab] committed with status (VALID) at localhost:7051
2020-12-07 18:16:33.467 CST [chaincodeCmd] ClientWait -> INFO 002 txid [626210823222a9e17401c568c2371eb38b8896e993ad74d43673b66c94465fab] committed with status (VALID) at localhost:9051
Chaincode definition committed on channel 'mychannel'
Using organization 1
Querying chaincode definition on peer0.org1 on channel 'mychannel'...
Attempting to Query committed status on peer0.org1, Retry after 3 seconds.
+ peer lifecycle chaincode querycommitted --channelID mychannel --name basic
+ res=0
Committed chaincode definition for chaincode 'basic' on channel 'mychannel':
Version: 1.0, Sequence: 1, Endorsement Plugin: escc, Validation Plugin: vscc, Approvals: [Org1MSP: true, Org2MSP: true]
Query chaincode definition successful on peer0.org1 on channel 'mychannel'
Using organization 2
Querying chaincode definition on peer0.org2 on channel 'mychannel'...
Attempting to Query committed status on peer0.org2, Retry after 3 seconds.
+ peer lifecycle chaincode querycommitted --channelID mychannel --name basic
+ res=0
Committed chaincode definition for chaincode 'basic' on channel 'mychannel':
Version: 1.0, Sequence: 1, Endorsement Plugin: escc, Validation Plugin: vscc, Approvals: [Org1MSP: true, Org2MSP: true]
Query chaincode definition successful on peer0.org2 on channel 'mychannel'
Chaincode initialization is not required
```

## 使用couchDB

`./network up -s couchdb`

id string, color string, size int, owner string, appraisedValue int

peer chaincode invoke -o localhost:7050  -C mychannel -n basic -c '{"function":"CreateAsset","Args":["assets7","black","32","OahcUil","2000"]}'

peer chaincode invoke  -C mychannel -n basic  -c '{"function":"CreateAsset","Args":["assets7","black","32","OahcUil","2000"]}'

peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'

peer chaincode query -C mychannel -n basic -c '{"function":"CreateAsset","Args":["asset8","black","32","OahcUil","2000"]}'

## ledger

peer chaincode invoke -C mychannel -n ledger -c '{"Args":["InitLedger"]}' -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

## 参考资料
https://www.jianshu.com/p/b131a8503559