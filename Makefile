all:
	mkdir bin
	GOOS=linux GOARCH=amd64 go build -o bin/nsrlferret
	GOOS=windows GOARCH=amd64 go build -o bin/nsrlferret.exe

linux:
	mkdir bin
	GOOS=linux GOARCH=amd64 go build -o bin/nsrlferret

windows:
	mkdir bin
	GOOS=windows GOARCH=amd64 go build -o bin/nsrlferret.exe

run:
	go run .

clean:
	rm -rf bin
