package gmf_test

import (
	"fmt"
	"log"
	"sync"

	"github.com/3d0c/gmf"
)

type output struct {
	filename string
	codec    int
}

func Example_encodingMultiple() {
	o := []output{
		{"examples/sample-enc-mpeg1.mpg", gmf.AV_CODEC_ID_MPEG1VIDEO},
		{"examples/sample-enc-mpeg2.mpg", gmf.AV_CODEC_ID_MPEG2VIDEO},
		{"examples/sample-enc-mpeg4.mp4", gmf.AV_CODEC_ID_MPEG4},
	}

	var wg sync.WaitGroup
	for _, v := range o {
		wg.Add(1)
		go encodingMultipleWorker(v, &wg)
	}
	wg.Wait()

	fmt.Println("frames written to examples/sample-enc-mpeg4.mp4")
	// Output: frames written to examples/sample-enc-mpeg4.mp4
}

func encodingMultipleWorker(item output, wg *sync.WaitGroup) {
	defer wg.Done()
	codec, err := gmf.FindEncoder(item.codec)
	if err != nil {
		log.Fatalln(err)
	}
	videoEncCtx := gmf.NewCodecCtx(codec)
	if videoEncCtx == nil {
		log.Fatalln(err)
	}
	defer videoEncCtx.Free()

	oCtx, err := gmf.NewOutputCtx(item.filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer oCtx.Free()

	videoEncCtx.
		SetBitRate(400000).
		SetWidth(320).
		SetHeight(200).
		SetTimeBase(gmf.AVR{Num: 1, Den: 25}).
		SetPixFmt(gmf.AV_PIX_FMT_YUV420P)

	switch item.codec {
	case gmf.AV_CODEC_ID_MPEG1VIDEO:
		videoEncCtx.SetMbDecision(gmf.FF_MB_DECISION_RD)
	case gmf.AV_CODEC_ID_MPEG4:
		videoEncCtx.SetProfile(gmf.FF_PROFILE_MPEG4_SIMPLE)
	}

	if oCtx.IsGlobalHeader() {
		videoEncCtx.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}

	s := oCtx.NewStream(codec)
	if s == nil {
		log.Fatalln(fmt.Sprintf("Unable to create stream for videoEnc [%s]\n", codec.LongName()))
	}
	defer s.Free()

	err = videoEncCtx.Open(nil)
	if err != nil {
		log.Fatalln(err)
	}
	s.DumpContextCodec(videoEncCtx)

	oCtx.SetStartTime(0)

	err = oCtx.WriteHeader()
	if err != nil {
		log.Fatalln(err)
	}
	oCtx.Dump()

	var p *gmf.Packet
	i := int64(0)
	for frame := range SyntheticVideoNewFrame(320, 200, gmf.AV_PIX_FMT_YUV420P) {
		frame.SetPts(i)
		p, err = frame.Encode(videoEncCtx)
		if err != nil {
			log.Fatalln("frame.Encode error", err)
		}
		if p != nil {
			if p.Pts() != gmf.AV_NOPTS_VALUE {
				p.SetPts(gmf.RescaleQ(p.Pts(), videoEncCtx.TimeBase(), s.TimeBase()))
			}

			if p.Dts() != gmf.AV_NOPTS_VALUE {
				p.SetDts(gmf.RescaleQ(p.Dts(), videoEncCtx.TimeBase(), s.TimeBase()))
			}

			err = oCtx.WritePacket(p)
			if err != nil {
				log.Fatalln("oCtx.WritePacket", err)
			}
			log.Printf("Write %s [%d] size=%v pts=%v dts=%v\n", item.filename, item.codec, p.Size(), p.Pts(), p.Dts())
			p.Free()
		}
		frame.Free()
		i++
	}
}
