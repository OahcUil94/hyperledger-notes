Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"
  config.vm.box_check_update = false
  config.vm.network "private_network", ip: "192.168.33.11"
  config.vm.synced_folder "~/code/go", "/home/vagrant/code/go", type: "nfs"
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "4096"
    vb.cpus = 2
    vb.name = "dev-node"
  end
  config.vm.provision "shell", path: "install.sh"
end
