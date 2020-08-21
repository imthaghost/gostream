<p align="center">
Simple and efficient live broadcast API. Currently supports commonly used transmission protocols, file formats, and encoding formats. Built in pure Golang, for high performance, and cross-platform.
</p>
<br>
<p align="center">
   <a href="https://goreportcard.com/report/github.com/imthaghost/goclone"><img src="https://goreportcard.com/badge/github.com/imthaghost/goclone"></a>
</p>
<br>

#### Supported transport protocols
- RTMP
- AMF
- HLS
- HTTP-FLV

#### Supported container formats
- FLV
- TS

#### Supported encoding formats
- H264
- AAC
- MP3

#### Boot from Docker
```bash
# run
docker run -p 1935:1935 -p 7001:7001 -p 7002:7002 -p 8090:8090 -d imthaghost/gostream
``` 
#### Compile from source
```bash
# download the source code
git clone https://github.com/imthaghost/gostream.git
# change to gostream directory and build
go build 
# or
make build
```

#### Options
```bash
./gostream  -h
Usage of ./gostream:
      --api_addr string       HTTP manage interface server listen address (default ":8090")
      --config_file string    configure filename (default "livego.yaml")
      --flv_dir string        output flv file at flvDir/APP/KEY_TIME.flv (default "tmp")
      --gop_num int           gop num (default 1)
      --hls_addr string       HLS server listen address (default ":7002")
      --hls_keep_after_end    Maintains the HLS after the stream ends
      --httpflv_addr string   HTTP-FLV server listen address (default ":7001")
      --level string          Log level (default "info")
      --read_timeout int      read time out (default 10)
      --rtmp_addr string      RTMP server listen address
```
