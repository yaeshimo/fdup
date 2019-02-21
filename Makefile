
all : build

build :
	go build

run : build
	./fdup -help
	./fdup -version
	./fdup -verbose -paths "testdata"

test :
	go test -race -v

arm7 :
	GOOS=linux GOARCH=arm GOARM=7 go build

win :
	GOOS=windows GOARCH=amd64 go build
