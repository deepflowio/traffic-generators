# build
go build -o tcpserver tcpserver.go
go build -o tcpclient tcpclient.go

# deploy
scp tcpclient 10.1.19.9:~
scp tcpserver 10.1.19.9:~
