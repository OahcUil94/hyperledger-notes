# Hyperledger 超级账本

- 记录时间: 2020-12-15
- Fabric版本: 2.2
- CouchDB版本: 3.1.1   

## 开发环境

### virtualbox

### vagrant

### git

### golang

### docker

## Fabric

docker镜像:

docker pull hyperledger/fabric-peer:2.3
docker pull hyperledger/fabric-orderer:2.3
docker pull hyperledger/fabric-baseos:2.3
docker pull hyperledger/fabric-tools:2.3

## 基础组件

fabric-ca, 数据的增删改查, 需要有数字签名
peer: 存储ledger和blockchain存储的位置
（背书策略, 决定ledger是否更新）

order服务: 排序服务, 创建区块

背书策略, 如何达成共识

msp

## 阿里云服务器启动网络的时候失败

## 相关资料

- [https://github.com/IBM-Blockchain-Archive/marbles](https://github.com/IBM-Blockchain-Archive/marbles)


可能是/etc/resolv.conf文件的问题
options timeout:2 attempts:3 rotate single-request-reopen
