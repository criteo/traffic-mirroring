static:
	cd mirror/server && statik -f -src ./public

proto:
	cd mirror && protoc -I=. --go_out=. *.proto
