test: vm_is_up
	testing/ssh.sh -t "source /data/O_O/conf/rc && cd /data/O_O && make test_server"
	testing/ssh.sh -t "source /data/O_O/conf/rc && cd /data/O_O && make integration"

ssh: vm_is_up
	cd testing/vm && vagrant ssh

halt:
	cd testing/vm && vagrant halt

clean: reload
reload: vm_is_up
	cd testing/vm && vagrant reload

vm_is_up:
	testing/vm_ok.sh
