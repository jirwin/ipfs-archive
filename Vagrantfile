Vagrant.configure(2) do |config|
  config.vm.box = "debian/jessie64"

  config.vm.provider "vmware_fusion" do |v|
    v.vmx["memsize"] = "2048"
    v.vmx["numvcpus"] = 2
  end

  config.vm.provider "virtualbox" do |v|
    v.memory = 2048
    v.cpus = 2
  end

  if Vagrant.has_plugin?("vagrant-cachier")
    config.cache.scope = :box
  end

  config.vm.provision "shell" do |s|
    s.path = "scripts/provision.sh"
  end

  $script = <<-SCRIPT
  echo "export GOPATH=$HOME/go" >> /home/vagrant/.profile
  echo "export IPFS_PATH=/data/ipfs" >> /home/vagrant/.profile
  echo "export GOBIN=/home/vagrant/go/bin" >> /home/vagrant/.profile
  echo "export PATH=$PATH:/opt/go/bin:/home/vagrant/go/bin" >> /home/vagrant/.profile
  source /home/vagrant/.profile
  sudo chown vagrant:vagrant /home/vagrant/go
  mkdir $GOPATH/pkg $GOPATH/bin
  if [ ! -d "/data/ipfs" ]; then
    sudo mkdir -p /data/ipfs
    sudo chown -R vagrant:vagrant /data/ipfs
    ipfs init
  fi
  SCRIPT

  config.vm.synced_folder "../../../", "/home/vagrant/go/src", type: "sshfs"
  config.vm.network "forwarded_port", guest: 7002, host: 7002


  config.vm.provision "shell", inline: $script, privileged: false
end
