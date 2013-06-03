exec { "apt-update":
  command => "/usr/bin/apt-get update"
}

Exec["apt-update"] -> Package <| |>

package { ["vim", "curl", "git", "bzr", "make", "g++", "gcc", "stud" ]:
  ensure => present,
}

file { ['/data/gopath', '/data/gopath/src', '/data/gopath/src/github.com', '/data/gopath/src/github.com/racker']:
  ensure => "directory",
  owner => "vagrant",
  group => "vagrant",
}

exec { "O_O":
  command => "/bin/mkdir -p /data/O_O && /bin/mount --bind /data/gopath/src/github.com/racker/go-agent-endpoint/testing/vm/O_O /data/O_O",
  creates => "/data/O_O/conf",
}

exec { "install_go":
  command => "/usr/bin/curl 'https://go.googlecode.com/files/go1.1.linux-amd64.tar.gz' | /bin/tar zx -C /data/O_O",
  creates => "/data/O_O/go",
  require => [ Exec['O_O'], Package['curl'] ],
  user => "vagrant",
}

$GO_ENV = ["GOROOT=/data/O_O/go", "GOPATH=/data/gopath", "GOBIN=/data/gopath/bin"]

exec { "install_gocheck":
  command => "/data/O_O/go/bin/go get -u launchpad.net/gocheck",
  environment => $GO_ENV,
  require => [ File['/data/gopath/src/github.com'], File['/home/vagrant/.profile'], Exec['O_O'], Package['bzr'] ],
  user => "vagrant",
}

exec { "install_colorgo":
  command => "/data/O_O/go/bin/go get -u github.com/songgao/colorgo",
  environment => $GO_ENV,
  require => [ File['/data/gopath/src/github.com'], File['/home/vagrant/.profile'], Exec['O_O'] ],
  user => "vagrant",
}

exec { "install_pingpong":
  command => "/data/O_O/go/bin/go get -u github.com/songgao/pingpong",
  environment => $GO_ENV,
  require => [ File['/data/gopath/src/github.com'], File['/home/vagrant/.profile'], Exec['O_O'] ],
  user => "vagrant",
}

file { '/home/vagrant/.profile':
  ensure => file,
  content => template('/vagrant/manifests/templates/profile.erb'),
}

exec { "repo_init":
  command => '/bin/bash -c "cd /data/gopath/src/github.com/racker/go-agent-endpoint && git submodule update --init --recursive"',
  require => [ Package['git'] ],
}
