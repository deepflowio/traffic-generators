# 简介

用于在两个服务器之间生成如下流量：
- 接近 100K 的活跃 TCP Flow
- 每个 Flow 的生存时间超过 60 秒，使得产生 ForceReport
- 接近 10K 的 PPS，即并不是每个 Flow 每一秒都有包，不制造 PPS 压力
- 接近 1K/s 的新建 Flow 速率，且每秒开始时突发完成
- 从 Flow 中提取的 (ClientIP, ServerIP, ServerPort) 基数够大，制造指标数据的压力

# 客户端、服务端配置

```bash
ulimit -n 400000
```

# 服务端配置

```bash
echo 10000 > /proc/sys/net/core/somaxconn
echo 1 > /proc/sys/net/ipv4/tcp_tw_reuse
echo 1 > /proc/sys/net/ipv4/tcp_tw_recycle
echo 0 > /proc/sys/net/ipv4/tcp_syncookies
```

# 客户端、服务端增加 IP

```bash
# 客户端
ip addr add dev lo 192.168.10.100/24
ip addr add dev lo 192.168.10.101/24
ip addr add dev lo 192.168.10.102/24
ip addr add dev lo 192.168.10.103/24
ip addr add dev lo 192.168.10.104/24
ip addr add dev lo 192.168.10.105/24

# 服务端
ip addr add dev lo 192.168.10.200/24
ip addr add dev lo 192.168.10.201/24
ip addr add dev lo 192.168.10.202/24
ip addr add dev lo 192.168.10.203/24
ip addr add dev lo 192.168.10.204/24
ip addr add dev lo 192.168.10.205/24
```

# 客户端运行

```bash
./tcpclient \
    192.168.10.52,192.168.10.200,192.168.10.201,192.168.10.202,192.168.10.203,192.168.10.204,192.168.10.205 \
    192.168.10.51,192.168.10.100,192.168.10.101,192.168.10.102,192.168.10.103,192.168.10.104,192.168.10.105
```

# 服务端运行

```bash
./tcpserver
```
