all:
	[ -d bin ] || mkdir bin
	GOOS=linux GOARCH=amd64 go build -o bin/nsrlferret -ldflags "-s -w"
	GOOS=windows GOARCH=amd64 go build -o bin/nsrlferret.exe -ldflags "-s -w"

linux:
	[ -d bin ] || mkdir bin
	GOOS=linux GOARCH=amd64 go build -o bin/nsrlferret -ldflags "-s -w"

windows:
	[ -d bin ] || mkdir bin
	GOOS=windows GOARCH=amd64 go build -o bin/nsrlferret.exe -ldflags "-s -w"

run:
	go run .

clean:
	rm -rf bin
