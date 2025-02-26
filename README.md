# Guide

1. cmd/main.go - http server with generation multiplier logic
2. internal/client - client with generating sequence and calculating RTP

```bash
go run cmd/main.go -rtp=0.95
go run intermal/client.go -seq=100000
```
