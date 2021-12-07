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

type CodecParameters struct {
	avCodecParameters *C.struct_AVCodecParameters
}

func NewCodecParameters() *CodecParameters {
	return &CodecParameters{
		avCodecParameters: C.avcodec_parameters_alloc(),
	}
}

func (cp *CodecParameters) ExtraData() []byte {
	return C.GoBytes(unsafe.Pointer(cp.avCodecParameters.extradata), C.int(cp.avCodecParameters.extradata_size))
}

func (cp *CodecParameters) Free() {
	C.avcodec_parameters_free(&cp.avCodecParameters)
}

func (cp *CodecParameters) CodecType() int {
	return int(cp.avCodecParameters.codec_type)
}

func (cp *CodecParameters) CodecId() int {
	return int(cp.avCodecParameters.codec_id)
}

func (cp *CodecParameters) BitRate() int64 {
	return int64(cp.avCodecParameters.bit_rate)
}

func (cp *CodecParameters) Width() int {
	return int(cp.avCodecParameters.width)
}

// Format
// video: the pixel format, the value corresponds to enum AVPixelFormat.
// audio: the sample format, the value corresponds to enum AVSampleFormat.
func (cp *CodecParameters) Format() int32 {
	return int32(cp.avCodecParameters.format)
}

func (cp *CodecParameters) Height() int {
	return int(cp.avCodecParameters.height)
}

func (cp *CodecParameters) GetVideoSize() string {
	return fmt.Sprintf("%dx%d", cp.Width(), cp.Height())
}

func (cp *CodecParameters) GetAspectRation() AVRational {
	return AVRational(cp.avCodecParameters.sample_aspect_ratio)
}

func (cp *CodecParameters) FrameSize() int {
	return int(cp.avCodecParameters.frame_size)
}

func (cp *CodecParameters) Channels() int {
	return int(cp.avCodecParameters.channels)
}

func (cp *CodecParameters) GetDefaultChannelLayout() int {
	return int(C.av_get_default_channel_layout(C.int(cp.Channels())))
}

func (cp *CodecParameters) FromContext(cc *CodecCtx) error {
	ret := int(C.avcodec_parameters_from_context(cp.avCodecParameters, cc.avCodecCtx))
	if ret < 0 {
		return AvError(ret)
	}

	return nil
}

func (cp *CodecParameters) ToContext() (cc *CodecCtx, err error) {
	cc = &CodecCtx{}
	ret := int(C.avcodec_parameters_to_context(cc.avCodecCtx, cp.avCodecParameters))
	if ret < 0 {
		return cc, AvError(ret)
	}
	cc.codec, err = FindDecoder(cp.CodecId())
	if err != nil {
		return
	}
	return
}
