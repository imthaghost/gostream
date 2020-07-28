package cache

import (
	"bytes"

	"github.com/imthaghost/gostream/avv"
	"github.com/imthaghost/gostream/protocol/amf"

	log "github.com/sirupsen/logrus"
)

const (
	SetDataFrame string = "@setDataFrame"
	OnMetaData   string = "onMetaData"
)

var setFrameFrame []byte

func init() {
	b := bytes.NewBuffer(nil)
	encoder := &amf.Encoder{}
	if _, err := encoder.Encode(b, SetDataFrame, amf.AMF0); err != nil {
		log.Fatal(err)
	}
	setFrameFrame = b.Bytes()
}

type SpecialCache struct {
	full bool
	p    *avv.Packet
}

func NewSpecialCache() *SpecialCache {
	return &SpecialCache{}
}

func (specialCache *SpecialCache) Write(p *avv.Packet) {
	specialCache.p = p
	specialCache.full = true
}

func (specialCache *SpecialCache) Send(w avv.WriteCloser) error {
	if !specialCache.full {
		return nil
	}
	return w.Write(specialCache.p)
}
