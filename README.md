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



