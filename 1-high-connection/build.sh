# go build -o trafficgen main.go
go build -o tcpserver tcpserver.go
go build -o tcpclient tcpclient.go

# deploy
scp tcpclient 10.50.10.51:~
scp tcpserver 10.50.10.52:~
