#!/usr/bin/env bash

echo '====set timezone===='
timedatectl set-timezone Asia/Shanghai
timedatectl set-local-rtc 0

echo '====set aliyun yum repo===='
mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
sed -i -e '/mirrors.cloud.aliyuncs.com/d' -e '/mirrors.aliyuncs.com/d' /etc/yum.repos.d/CentOS-Base.repo
yum clean all && yum makecache -y

echo '====set baidu dns server===='
mv /etc/NetworkManager/NetworkManager.conf /etc/NetworkManager/NetworkManager.conf.backup

cat > /etc/NetworkManager/NetworkManager.conf <<EOF
[main]
dns=none

[logging]
EOF

cat >> /etc/resolv.conf <<EOF
nameserver 180.76.76.76
EOF

systemctl restart NetworkManager.service

echo '====close not need service===='
systemctl stop postfix && systemctl disable postfix

echo '====set only journald log===='
mkdir /var/log/journal
mkdir /etc/systemd/journald.conf.d
cat > /etc/systemd/journald.conf.d/99-prophet.conf <<EOF
[Journal]
Storage=persistent
Compress=yes
SyncIntervalSec=5m
RateLimitInterval=30s
RateLimitBurst=1000
SystemMaxUse=10G
SystemMaxFileSize=200M
MaxRetentionSec=2week
ForwardToSyslog=no
EOF
systemctl restart systemd-journald

echo '====load overlay, br_netfilter module===='
cat > /etc/modules-load.d/containerd.conf <<EOF
overlay
br_netfilter
EOF
modprobe overlay
modprobe br_netfilter

echo '==== yum update and install common package===='
yum update -y
yum install -y vim net-tools telnet bind-utils wget yum-utils device-mapper-persistent-data lvm2

echo '====config system k8s network params===='
cat > /etc/sysctl.d/99-kubernetes-cri.conf <<EOF
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF
sysctl --system

echo '==== install docker===='
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum clean all && yum makecache -y
yum install -y docker-ce docker-ce-cli containerd.io
systemctl start docker
cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "registry-mirrors" : [
    "https://thd69qis.mirror.aliyuncs.com",
    "https://f1361db2.m.daocloud.io",
    "https://mirror.ccs.tencentyun.com",
    "https://reg-mirror.qiniu.com",
    "https://docker.mirrors.ustc.edu.cn/",
    "https://registry.docker-cn.com"
  ]
}
EOF
systemctl daemon-reload
systemctl restart docker
systemctl enable docker

echo '====add vagrant user to docker group===='
egrep "^docker" /etc/group >& /dev/null
if [ $? -ne 0 ]
then
  groupadd docker
fi
usermod -aG docker vagrant

echo '====disable selinux===='
setenforce 0
sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

echo '====disable swap===='
swapoff -a
sed -i '/swap/s/^/#/g' /etc/fstab

echo '====disable firewalld===='
systemctl stop firewalld
systemctl disable firewalld

echo '====clear iptable rules===='
yum install -y iptables-services
systemctl start iptables
systemctl enable iptables
service iptables save
iptables -F

mkdir /etc/systemd/system/docker.service.d
echo "ExecStartPost=/sbin/iptables -P FORWARD ACCEPT" >> /etc/systemd/system/docker.service.d/docker.conf
systemctl daemon-reload
systemctl restart docker
