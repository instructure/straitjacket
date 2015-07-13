# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define 'straitjacket'
  config.vm.box = "ubuntu/vivid64"

  config.vm.network "forwarded_port", guest: 8081, host: 8081

  config.vm.provider 'virtualbox' do |v|
    v.memory = 2048
    v.cpus = 2
  end

  config.vm.provider "virtualbox" do |vb|
    lfs_disk = "btrfs-fs.vdi"
    unless File.exist?(lfs_disk)
      vb.customize ['createhd', '--filename', lfs_disk, '--size', 20 * 1024]
    end
    vb.customize ['storageattach', :id, '--storagectl', 'SATAController', '--port', 1, '--device', 0, '--type', 'hdd', '--medium', lfs_disk]
  end

  config.vm.synced_folder ".", "/home/vagrant/straitjacket"

  config.vm.provision "shell", inline: <<-SCRIPT
  set -e
  apt-get update
  apt-get install -y btrfs-tools apparmor-utils
  mkfs.btrfs /dev/sdb
  mkdir /var/lib/docker
  mount /dev/sdb /var/lib/docker
  echo /dev/sdb /var/lib/docker btrfs defaults 0 0 >> /etc/fstab
  wget -qO- https://get.docker.com/ | sh
  usermod -aG docker vagrant
  SCRIPT
end
