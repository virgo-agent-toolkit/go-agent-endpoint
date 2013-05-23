test: vm_is_up
	cd testing/vm && vagrant ssh -c 'cd /data/O_O && make test_server'
	cd testing/vm && vagrant ssh -c 'cd /data/O_O && make test_agent'

ssh: vm_is_up
	cd testing/vm && vagrant ssh

clean: reload
reload: vm_is_up
	cd testing/vm && vagrant reload

vm_is_up:
	testing/vm_ok.sh
