package gmf

/*
#cgo pkg-config: libavcodec

#include "libavcodec/avcodec.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	AV_PKT_FLAG_KEY     = C.AV_PKT_FLAG_KEY     // The packet contains a keyframe
	AV_PKT_FLAG_CORRUPT = C.AV_PKT_FLAG_CORRUPT // The packet content is corrupted
)

type Packet struct {
	avPacket *C.struct_AVPacket
}

func NewPacket() *Packet {
	p := &Packet{}

	p.avPacket = C.av_packet_alloc()
	if p.avPacket == nil {
		return nil
	}

	p.avPacket.data = nil
	p.avPacket.size = 0

	return p
}

func (p *Packet) Pts() int64 {
	return int64(p.avPacket.pts)
}

func (p *Packet) SetPts(pts int64) *Packet {
	p.avPacket.pts = C.int64_t(pts)
	return p
}

func (p *Packet) Dts() int64 {
	return int64(p.avPacket.dts)
}

func (p *Packet) SetDts(val int64) *Packet {
	p.avPacket.dts = C.int64_t(val)
	return p
}

func (p *Packet) Flags() int {
	return int(p.avPacket.flags)
}

func (p *Packet) SetFlags(flags int) *Packet {
	p.avPacket.flags = C.int(flags)
	return p
}

func (p *Packet) Duration() int64 {
	return int64(p.avPacket.duration)
}

func (p *Packet) SetDuration(duration int64) *Packet {
	p.avPacket.duration = C.int64_t(duration)
	return p
}

func (p *Packet) StreamIndex() int {
	return int(p.avPacket.stream_index)
}

func (p *Packet) Size() int {
	return int(p.avPacket.size)
}

func (p *Packet) Pos() int64 {
	return int64(p.avPacket.pos)
}

func (p *Packet) Data() []byte {
	return C.GoBytes(unsafe.Pointer(p.avPacket.data), C.int(p.avPacket.size))
}

// SetData [NOT SUGGESTED] should free data later
func (p *Packet) SetData(data []byte) *Packet {
	p.avPacket.size = C.int(len(data))
	p.avPacket.data = (*C.uint8_t)(C.CBytes(data))
	return p
}

// FreeData free data when use SetData
func (p *Packet) FreeData() *Packet {
	if p.avPacket.data != nil {
		C.free(unsafe.Pointer(p.avPacket.data))
		p.avPacket.data = nil
		p.avPacket.size = 0
	}
	return p
}

func (p *Packet) Ref() *Packet {
	np := NewPacket()
	if np == nil {
		return np
	}

	if C.av_packet_ref(np.avPacket, p.avPacket) != 0 {
		np.Free()
	}

	return np
}

func (p *Packet) Dump() {
	if p.avPacket != nil {
		fmt.Printf("idx: %d\npts: %d\ndts: %d\nsize: %d\nduration:%d\npos:%d\ndata: % x\n", p.StreamIndex(), p.Pts(), p.Dts(), p.Size(), p.Duration(), p.Pos(), C.GoBytes(unsafe.Pointer(p.avPacket.data), 128))
		fmt.Println("------------------------------")
	}

}

func (p *Packet) SetStreamIndex(val int) *Packet {
	p.avPacket.stream_index = C.int(val)
	return p
}

// Free Free the packet, if the packet is reference counted, it will be unreferenced first.
func (p *Packet) Free() {
	if p.avPacket == nil {
		return
	}
	C.av_packet_free(&p.avPacket)
}

func (p *Packet) Unref() {
	if p.avPacket == nil {
		return
	}
	C.av_packet_unref(p.avPacket)
}

func (p *Packet) Time(timebase AVRational) int {
	return int(float64(timebase.AVR().Num) / float64(timebase.AVR().Den) * float64(p.Pts()))
}
