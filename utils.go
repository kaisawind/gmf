package gmf

/*
#cgo pkg-config: libavcodec libavutil

#include "libavutil/avutil.h"
#include "libavutil/error.h"
#include "libavutil/mathematics.h"
#include "libavutil/rational.h"
#include "libavutil/samplefmt.h"
#include "libavcodec/avcodec.h"
#include "libavutil/imgutils.h"

uint32_t return_int (int num) {
	return (uint32_t)(num);
}

uint8_t * gmf_alloc_buffer(int32_t fmt, int width, int height) {
	int numBytes = av_image_get_buffer_size(fmt, width, height, 0);
	return (uint8_t *) av_malloc(numBytes*sizeof(uint8_t));
}

*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

type AVRational C.struct_AVRational

type AVR struct {
	Num int
	Den int
}

const (
	AVERROR_EOF = -541478725
	// AV_ROUND_PASS_MINMAX = 8192
)

var (
	AV_TIME_BASE   int        = C.AV_TIME_BASE
	AV_TIME_BASE_Q AVRational = AVRational{1, C.int(AV_TIME_BASE)}
)

func (a AVR) AVRational() AVRational {
	return AVRational{C.int(a.Num), C.int(a.Den)}
}

func (a AVR) String() string {
	return fmt.Sprintf("%d/%d", a.Num, a.Den)
}

func (a AVR) Av2qd() float64 {
	return float64(a.Num) / float64(a.Den)
}

func (a AVR) Invert() AVR {
	return AVR{Num: a.Den, Den: a.Num}
}

func (a AVRational) AVR() AVR {
	return AVR{Num: int(a.num), Den: int(a.den)}
}

func AvError(averr int) error {
	errlen := 1024
	b := make([]byte, errlen)

	C.av_strerror(C.int(averr), (*C.char)(unsafe.Pointer(&b[0])), C.size_t(errlen))

	return errors.New(string(b[:bytes.Index(b, []byte{0})]))
}

func AvErrno(ret int) syscall.Errno {
	if ret < 0 {
		ret = -ret
	}
	return syscall.Errno(ret)
}

func RescaleQ(a int64, encBase AVRational, stBase AVRational) int64 {
	return int64(C.av_rescale_q(C.int64_t(a), C.struct_AVRational(encBase), C.struct_AVRational(stBase)))
}

func RescaleQRnd(a int64, encBase AVRational, stBase AVRational) int64 {
	return int64(C.av_rescale_q_rnd(C.int64_t(a), C.struct_AVRational(encBase), C.struct_AVRational(stBase), C.AV_ROUND_NEAR_INF|C.AV_ROUND_PASS_MINMAX))
}

func CompareTimeStamp(aTimestamp int, aTimebase AVRational, bTimestamp int, bTimebase AVRational) int {
	return int(C.av_compare_ts(C.int64_t(aTimestamp), C.struct_AVRational(aTimebase),
		C.int64_t(bTimestamp), C.struct_AVRational(bTimebase)))
}
func RescaleDelta(inTb AVRational, inTs int64, fsTb AVRational, duration int, last *int64, outTb AVRational) int64 {
	return int64(C.av_rescale_delta(C.struct_AVRational(inTb), C.int64_t(inTs), C.struct_AVRational(fsTb), C.int(duration), (*C.int64_t)(unsafe.Pointer(&last)), C.struct_AVRational(outTb)))
}

func Rescale(a, b, c int64) int64 {
	return int64(C.av_rescale(C.int64_t(a), C.int64_t(b), C.int64_t(c)))
}

func RescaleTs(pkt *Packet, encBase AVRational, stBase AVRational) {
	C.av_packet_rescale_ts(pkt.avPacket, C.struct_AVRational(encBase), C.struct_AVRational(stBase))
}

func GetSampleFmtName(fmt int32) string {
	return C.GoString(C.av_get_sample_fmt_name(fmt))
}

func AvInvQ(q AVRational) AVRational {
	avr := q.AVR()
	return AVRational{C.int(avr.Den), C.int(avr.Num)}
}
