package gmf_test

import (
	"github.com/3d0c/gmf"
	"testing"
)

// SyntheticVideoNewFrame Synthetic video generator. It produces 25 iterable frames.
// Used for tests.
func SyntheticVideoNewFrame(w, h int, fmt int32) chan *gmf.Frame {
	yield := make(chan *gmf.Frame)

	go func() {
		defer close(yield)
		for i := 0; i < 25; i++ {
			frame := gmf.NewFrame().SetWidth(w).SetHeight(h).SetFormat(fmt)

			if err := frame.ImgAlloc(); err != nil {
				return
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					frame.SetData(0, y*frame.LineSize(0)+x, x+y+i*3)
				}
			}

			// Cb and Cr
			for y := 0; y < h/2; y++ {
				for x := 0; x < w/2; x++ {
					frame.SetData(1, y*frame.LineSize(1)+x, 128+y+i*2)
					frame.SetData(2, y*frame.LineSize(2)+x, 64+x+i*5)
				}
			}

			yield <- frame
		}
	}()
	return yield
}

// SyntheticVideoN tmp
func SyntheticVideoN(N, w, h int, fmt int32) chan *gmf.Frame {
	yield := make(chan *gmf.Frame)

	go func() {
		defer close(yield)
		for i := 0; i < N; i++ {
			frame := gmf.NewFrame().SetWidth(w).SetHeight(h).SetFormat(fmt)

			if err := frame.ImgAlloc(); err != nil {
				return
			}

			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					frame.SetData(0, y*frame.LineSize(0)+x, x+y+i*3)
				}
			}

			// Cb and Cr
			for y := 0; y < h/2; y++ {
				for x := 0; x < w/2; x++ {
					frame.SetData(1, y*frame.LineSize(1)+x, 128+y+i*2)
					frame.SetData(2, y*frame.LineSize(2)+x, 64+x+i*5)
				}
			}

			yield <- frame
		}
	}()
	return yield
}

func TestAvError(t *testing.T) {
	if err := gmf.AvError(-2); err.Error() != "No such file or directory" {
		t.Fatalf("Expected error is 'No such file or directory', '%s' got\n", err.Error())
	}
}
