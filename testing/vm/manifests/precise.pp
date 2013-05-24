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

Package <| |> -> Exec["remaining steps"]
File <| |> -> Exec["remaining steps"]

exec { "remaining steps":
  command => '/vagrant/provision.sh',
}


