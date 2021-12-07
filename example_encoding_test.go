package gmf_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/3d0c/gmf"
)

func Example_encoding() {
	oFilename := "examples/sample-encoding1.mpg"
	dstWidth, dstHeight := 640, 480

	codec, err := gmf.FindEncoder(gmf.AV_CODEC_ID_MPEG1VIDEO)
	if err != nil {
		log.Fatal(err)
	}

	videoEncCtx := gmf.NewCodecCtx(codec)
	if videoEncCtx == nil {
		log.Fatal(errors.New("failed to create a new codec context"))
	}
	defer videoEncCtx.Free()

	oCtx, err := gmf.NewOutputCtx(oFilename)
	if err != nil {
		log.Fatal(errors.New("failed to create a new output context"))
	}
	defer oCtx.Free()

	videoEncCtx.
		SetBitRate(400000).
		SetWidth(dstWidth).
		SetHeight(dstHeight).
		SetTimeBase(gmf.AVR{Num: 1, Den: 25}).
		SetPixFmt(gmf.AV_PIX_FMT_YUV420P).
		SetProfile(gmf.FF_PROFILE_MPEG4_SIMPLE).
		SetMbDecision(gmf.FF_MB_DECISION_RD)

	if oCtx.IsGlobalHeader() {
		videoEncCtx.SetFlag(gmf.CODEC_FLAG_GLOBAL_HEADER)
	}

	s := oCtx.NewStream(codec)
	if s == nil {
		log.Fatal(errors.New(fmt.Sprintf("Unable to create stream for videoEnc [%s]\n", codec.LongName())))
	}
	defer s.Free()

	if err = videoEncCtx.Open(nil); err != nil {
		log.Fatal(err)
	}
	s.DumpContextCodec(videoEncCtx)

	oCtx.SetStartTime(0)

	if err = oCtx.WriteHeader(); err != nil {
		log.Fatal(err)
	}
	oCtx.Dump()

	var p *gmf.Packet
	i := int64(0)
	for frame := range SyntheticVideoNewFrame(videoEncCtx.Width(), videoEncCtx.Height(), videoEncCtx.PixFmt()) {
		frame.SetPts(i)

		p, err = frame.Encode(videoEncCtx)
		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}
			log.Printf("Write frame=%d size=%v pts=%v dts=%v\n", frame.Pts(), p.Size(), p.Pts(), p.Dts())
			p.Free()
		}
		frame.Free()
		i++
	}
	fmt.Println("frames written to examples/sample-encoding1.mpg")
	// Output: frames written to examples/sample-encoding1.mpg
}
