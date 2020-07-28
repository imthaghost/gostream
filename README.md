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
1. Download the source code `git clone https://github.com/imthaghost/gostream.git`
2. Go to the gostream directory and execute `go build` or `make build`

## Use
1. Start the service: execute the gostream binary file or `make run` to start the gostream service;
2. Get a channelkey(used for push the video stream) from `http://localhost:8090/control/get?room=movie` and copy data like your channelkey.
3. Upstream push: Push the video stream to `rtmp://localhost:1935/{appname}/{channelkey}` through the` RTMP` protocol(default appname is `live`), for example, use `ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/{appname}/{channelkey}` push([download demo flv](https://s3plus.meituan.net/v1/mss_7e425c4d9dcb4bb4918bbfa2779e6de1/mpack/default/demo.flv));
4. Downstream playback: The following three playback protocols are supported, and the playback address is as follows:
    - `RTMP`:`rtmp://localhost:1935/{appname}/movie`
    - `FLV`:`http://127.0.0.1:7001/{appname}/movie.flv`
    - `HLS`:`http://127.0.0.1:7002/{appname}/movie.m3u8`
   



