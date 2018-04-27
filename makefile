all:
	rice embed-go
	go  build .
	rm *rice-box.go
	echo 'done'
install: all
	mv ./mcli ${GOPATH}/bin/
