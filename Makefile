proto:
	cd mirror && protoc -I=. --go_out=. *.proto