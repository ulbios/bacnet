Vagrant.configure("2") do |config|
    config.vm.box = "ubuntu/xenial64"
    config.ssh.forward_agent = true
    config.ssh.forward_x11 = true

    config.vm.synced_folder "../bacnet/", "/home/vagrant/bacnet"

    config.vm.provider "virtualbox" do |v|
        v.customize ["modifyvm", :id, "--memory", 1024]
    end

    config.vm.define "foovm" do |devvm|
        devvm.vm.hostname = 'foovm'
        devvm.vm.network :private_network, ip: "10.0.123.2"
    end
end
