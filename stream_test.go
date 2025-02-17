package gmf_test

import (
	"log"
	"testing"

	"github.com/3d0c/gmf"
)

func TestStream(t *testing.T) {
	ctx := gmf.NewCtx()

	vc, err := gmf.FindEncoder("mpeg4")
	if err != nil {
		t.Fatal(err)
	}

	cc := gmf.NewCodecCtx(vc)
	if cc == nil {
		t.Fatal("Unable to allocate codec context")
	}
	defer cc.Free()

	if ctx.NewStream(vc) == nil {
		t.Fatal("Unable to create new stream")
	}

	td := CodecCtxTestData

	cc.SetWidth(td.width).SetHeight(td.height).SetTimeBase(td.timebase).SetPixFmt(td.pixfmt).SetBitRate(td.bitrate)

	if err := cc.Open(nil); err != nil {
		t.Fatal(err)
	}

	st := assert(ctx.GetStream(0)).(*gmf.Stream)

	st.DumpContextCodec(cc)

	if cc.Height() != td.height || cc.Width() != td.width {
		t.Fatalf("Expected dimension = %dx%d, %dx%d got\n", td.width, td.height, st.CodecPar().Width(), st.CodecPar().Height())
	}

	ctx.Free()

	log.Println("Stream is OK")
}

func TestStreamInputCtx(t *testing.T) {
	inputCtx, err := gmf.NewInputCtx(inputSampleFilename)
	if err != nil {
		t.Fatal(err)
	}

	ist := assert(inputCtx.GetStream(0)).(*gmf.Stream)

	if ist.CodecPar().Width() != inputSampleWidth || ist.CodecPar().Height() != inputSampleHeight {
		t.Fatalf("Expected dimension = %dx%d, %dx%d got\n", inputSampleWidth, inputSampleHeight, ist.CodecPar().Width(), ist.CodecPar().Height())
	}

	log.Printf("Input stream is OK, cnt: %d, %dx%d\n", inputCtx.StreamsCnt(), ist.CodecPar().Width(), ist.CodecPar().Height())

	inputCtx.Free()
}
