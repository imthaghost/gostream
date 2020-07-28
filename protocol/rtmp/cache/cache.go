package cache

import (
	"github.com/imthaghost/gostream/avv"
	"github.com/imthaghost/gostream/configure"
)

type Cache struct {
	gop      *GopCache
	videoSeq *SpecialCache
	audioSeq *SpecialCache
	metadata *SpecialCache
}

func NewCache() *Cache {
	return &Cache{
		gop:      NewGopCache(configure.Config.GetInt("gop_num")),
		videoSeq: NewSpecialCache(),
		audioSeq: NewSpecialCache(),
		metadata: NewSpecialCache(),
	}
}

func (cache *Cache) Write(p avv.Packet) {
	if p.IsMetadata {
		cache.metadata.Write(&p)
		return
	} else {
		if !p.IsVideo {
			ah, ok := p.Header.(avv.AudioPacketHeader)
			if ok {
				if ah.SoundFormat() == avv.SOUND_AAC &&
					ah.AACPacketType() == avv.AAC_SEQHDR {
					cache.audioSeq.Write(&p)
					return
				} else {
					return
				}
			}

		} else {
			vh, ok := p.Header.(avv.VideoPacketHeader)
			if ok {
				if vh.IsSeq() {
					cache.videoSeq.Write(&p)
					return
				}
			} else {
				return
			}

		}
	}
	cache.gop.Write(&p)
}

func (cache *Cache) Send(w avv.WriteCloser) error {
	if err := cache.metadata.Send(w); err != nil {
		return err
	}

	if err := cache.videoSeq.Send(w); err != nil {
		return err
	}

	if err := cache.audioSeq.Send(w); err != nil {
		return err
	}

	if err := cache.gop.Send(w); err != nil {
		return err
	}

	return nil
}
