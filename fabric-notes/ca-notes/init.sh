# 注册目录
export REGISTRAR_DIR=$PWD
# ca客户端目录
export FABRIC_CA_CLIENT_HOME=$REGISTRAR_DIR

# 组织目录
export ORG_DIR=$PWD/crypto-config/peerOrganizations/dummyOrg.com
export PEER_TLS=$PWD/peertls
# 组织中的peer目录
export PEER_DIR=$ORG_DIR/peers/peer0.dummyOrg.com
#
export REGISTRAR_DIR=$ORG_DIR/users/admin
export ADMIN_DIR=$ORG_DIR/users/Admin@dummyOrg.com
export TLS=$ORG_DIR/tlsca
mkdir -p $ORG_DIR/ca $ORG_DIR/msp $PEER_DIR $REGISTRAR_DIR $ADMIN_DIR $TLS
mkdir certsICA

fabric-ca-client enroll -m admin -u http://adminCA:adminpw@ca.root:7054


fabric-ca-client register --id.name ica.dummyOrg --id.type client --id.secret adminpw --csr.names C=ES,ST=Madrid,L=Madrid,O=dummyOrg.com --csr.cn ica.dummyOrg -m ica.dummyOrg --id.attrs '"hf.IntermediateCA=true"' -u http://ca.root:7054

docker-compose up -d ica.dummyOrg

cp -r $PWD/certsICA/* $ORG_DIR/ca