# -*- mode: ruby -*-
# vi: set ft=ruby :
#
Vagrant.require_version ">= 1.4.0"

BOX_NAME = ENV['BOX_NAME'] || "opscode-ubuntu-1310"
BOX_URI = ENV['BOX_URI'] || "http://opscode-vm-bento.s3.amazonaws.com/vagrant/virtualbox/opscode_ubuntu-13.10_chef-provisionerless.box"
VF_BOX_URI = ENV['BOX_URI'] || "http://opscode-vm-bento.s3.amazonaws.com/vagrant/vmware/opscode_ubuntu-13.10_chef-provisionerless.box"
AWS_REGION = ENV['AWS_REGION']
AWS_AMI    = ENV['AWS_AMI']

Vagrant.configure("2") do |config|
  # Setup virtual machine box. This VM configuration code is always executed.
  config.vm.box = BOX_NAME
  config.vm.box_url = BOX_URI

  config.vm.network "private_network", ip: "192.168.50.4"

  # Provision docker and new kernel if deployment was not done.
  # It is assumed Vagrant can successfully launch the provider instance.
  #if Dir.glob("#{File.dirname(__FILE__)}/.vagrant/machines/default/*/id").empty?
  #  # Add lxc-docker package
  #  pkg_cmd = "wget -q -O - https://get.docker.io/gpg | apt-key add -;" \
  #    "echo deb http://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list;" \
  #    "apt-get update -qq; apt-get install -q -y --force-yes lxc-docker; "
  #  pkg_cmd << "apt-get update -qq; apt-get clean;"
  #  pkg_cmd << "sudo usermod -a -G docker vagrant;"
  #  config.vm.provision :shell, :inline => pkg_cmd
  #end
  config.vm.provision :shell, :inline => 'echo DOCKER_OPTS=\"-H tcp://0.0.0.0:4243 -H unix:///var/run/docker.sock -bip=10.2.0.10/16\" > /etc/default/docker'
  config.vm.provision :shell, :inline => '
  echo service docker restart >> /etc/rc.local
  chmod +x /etc/rc.local 2> /dev/null'
  config.vm.provision "docker", version: "0.7.6"
  config.vm.provision :shell, :inline => '/etc/rc.local'
end


# Providers were added on Vagrant >= 1.1.0
Vagrant::VERSION >= "1.1.0" and Vagrant.configure("2") do |config|
  config.vm.provider :vmware_fusion do |f, override|
    override.vm.box = BOX_NAME
    override.vm.box_url = VF_BOX_URI
    f.vmx["memsize"] = "2048"
    f.vmx["numvcpus"] = "2"
  end

  config.vm.provider :virtualbox do |vb|
    config.vm.box = BOX_NAME
    config.vm.box_url = BOX_URI
    #memory
    vb.customize ["modifyvm", :id, "--memory", "2048"]
  end
end
