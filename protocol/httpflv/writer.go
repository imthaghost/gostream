package httpflv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/imthaghost/gostream/avv"
	"github.com/imthaghost/gostream/protocol/amf"
	"github.com/imthaghost/gostream/utils/pio"
	"github.com/imthaghost/gostream/utils/uid"

	log "github.com/sirupsen/logrus"
)

const (
	headerLen   = 11
	maxQueueNum = 1024
)

type FLVWriter struct {
	Uid string
	avv.RWBaser
	app, title, url string
	buf             []byte
	closed          bool
	closedChan      chan struct{}
	ctx             http.ResponseWriter
	packetQueue     chan *avv.Packet
}

func NewFLVWriter(app, title, url string, ctx http.ResponseWriter) *FLVWriter {
	ret := &FLVWriter{
		Uid:         uid.NewId(),
		app:         app,
		title:       title,
		url:         url,
		ctx:         ctx,
		RWBaser:     avv.NewRWBaser(time.Second * 10),
		closedChan:  make(chan struct{}),
		buf:         make([]byte, headerLen),
		packetQueue: make(chan *avv.Packet, maxQueueNum),
	}

	ret.ctx.Write([]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09})
	pio.PutI32BE(ret.buf[:4], 0)
	ret.ctx.Write(ret.buf[:4])
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Error("SendPacket error: ", err)
			ret.closed = true
		}
	}()
	return ret
}

func (flvWriter *FLVWriter) DropPacket(pktQue chan *avv.Packet, info avv.Info) {
	log.Warningf("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(avv.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				log.Debug("insert keyframe to queue")
				pktQue <- tmpPkt
			}

			if len(pktQue) > maxQueueNum-10 {
				<-pktQue
			}
			// drop other packet
			<-pktQue
		}
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			log.Debug("insert audio to queue")
			pktQue <- tmpPkt
		}
	}
	log.Debug("packet queue len: ", len(pktQue))
}

func (flvWriter *FLVWriter) Write(p *avv.Packet) (err error) {
	err = nil
	if flvWriter.closed {
		err = fmt.Errorf("flvwrite source closed")
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("FLVWriter has already been closed:%v", e)
		}
	}()

	if len(flvWriter.packetQueue) >= maxQueueNum-24 {
		flvWriter.DropPacket(flvWriter.packetQueue, flvWriter.Info())
	} else {
		flvWriter.packetQueue <- p
	}

	return
}

func (flvWriter *FLVWriter) SendPacket() error {
	for {
		p, ok := <-flvWriter.packetQueue
		if ok {
			flvWriter.RWBaser.SetPreTime()
			h := flvWriter.buf[:headerLen]
			typeID := avv.TAG_VIDEO
			if !p.IsVideo {
				if p.IsMetadata {
					var err error
					typeID = avv.TAG_SCRIPTDATAAMF0
					p.Data, err = amf.MetaDataReform(p.Data, amf.DEL)
					if err != nil {
						return err
					}
				} else {
					typeID = avv.TAG_AUDIO
				}
			}
			dataLen := len(p.Data)
			timestamp := p.TimeStamp
			timestamp += flvWriter.BaseTimeStamp()
			flvWriter.RWBaser.RecTimeStamp(timestamp, uint32(typeID))

			preDataLen := dataLen + headerLen
			timestampbase := timestamp & 0xffffff
			timestampExt := timestamp >> 24 & 0xff

			pio.PutU8(h[0:1], uint8(typeID))
			pio.PutI24BE(h[1:4], int32(dataLen))
			pio.PutI24BE(h[4:7], int32(timestampbase))
			pio.PutU8(h[7:8], uint8(timestampExt))

			if _, err := flvWriter.ctx.Write(h); err != nil {
				return err
			}

			if _, err := flvWriter.ctx.Write(p.Data); err != nil {
				return err
			}

			pio.PutI32BE(h[:4], int32(preDataLen))
			if _, err := flvWriter.ctx.Write(h[:4]); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("closed")
		}

	}

	return nil
}

func (flvWriter *FLVWriter) Wait() {
	select {
	case <-flvWriter.closedChan:
		return
	}
}

func (flvWriter *FLVWriter) Close(error) {
	log.Debug("http flv closed")
	if !flvWriter.closed {
		close(flvWriter.packetQueue)
		close(flvWriter.closedChan)
	}
	flvWriter.closed = true
}

func (flvWriter *FLVWriter) Info() (ret avv.Info) {
	ret.UID = flvWriter.Uid
	ret.URL = flvWriter.url
	ret.Key = flvWriter.app + "/" + flvWriter.title
	ret.Inter = true
	return
}
