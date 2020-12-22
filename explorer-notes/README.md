# blockchain explorer

hyperledger下的[区块链浏览器](https://github.com/hyperledger/blockchain-explorer)

- 记录时间: 2020-12-22
- 版本信息: v1.1.3

## docker部署blockchain explorer

### 先决条件

部署起`fabric-samples/test-network`提供的fabric测试网络, 并部署一个链码: 

1. `./network.sh up createChannel -c mychannel -s couchdb`
2. `./network.sh deployCC`
3. 初始化链码并写入相关的数据
4. 以上操作完成后, 假设organizations的目录位置是`fabric-notes/fabric-samples/test-network/organizations`

```bash
.
├── cryptogen
│ ├── crypto-config-orderer.yaml
│ ├── crypto-config-org1.yaml
│ └── crypto-config-org2.yaml
├── fabric-ca
│ ├── ordererOrg
│ ├── org1
│ ├── org2
│ └── registerEnroll.sh
├── ordererOrganizations
│ └── example.com
└── peerOrganizations
    ├── org1.example.com
    └── org2.example.com
```

### 下载必要文件

```bash
wget https://cdn.jsdelivr.net/gh/hyperledger/blockchain-explorer@v1.1.3/examples/net1/config.json
wget https://cdn.jsdelivr.net/gh/hyperledger/blockchain-explorer@v1.1.3/examples/net1/connection-profile/first-network.json -P connection-profile
wget https://cdn.jsdelivr.net/gh/hyperledger/blockchain-explorer@v1.1.3/docker-compose.yaml
```

生成的目录结构:

```bash
.
├── config.json
├── connection-profile
│ └── first-network.json
└── docker-compose.yaml
```

- `docker-compose.yaml`中的volumes配置了浏览器运行所需要的文件
- `connection-profile/first-network.json`文件包含了浏览器登录的账号密码, 默认是`exploreradmin:exploreradminpw`

### 修改docker-compose.yaml文件

```yaml
services:
  ...
  explorer.mynetwork.com:
    image: hyperledger/explorer:latest
    environment:
      ...
      - DISCOVERY_AS_LOCALHOST=false
    volumes:
      - ./config.json:/opt/explorer/app/platform/fabric/config.json
      - ./connection-profile:/opt/explorer/app/platform/fabric/connection-profile
      - ../fabric-notes/fabric-samples/test-network/organizations:/tmp/crypto
      - walletstore:/opt/wallet
    ports:
      - 8080:8080
    ...   
```

### 启动服务

`docker-compose up -d`

### 停止服务

docker-compose down --volumes --remove-orphans