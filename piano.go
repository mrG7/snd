package main

import (
	"dasa.cc/piano/snd"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
)

// TODO this entire file is a hack job

var hasmn = []struct {
	left, right bool
	pos         int
}{
	{false, true, 0},  // C
	{true, true, 2},   // D
	{true, false, 4},  // E
	{false, true, 5},  // F
	{true, true, 7},   // G
	{true, true, 9},   // A
	{true, false, 11}, // B
}

type Piano struct {
	snd.Sound

	keys []float64
	idx  int
}

func NewPiano() *Piano {
	wf := &Piano{}
	wf.Sound = snd.Mono(nil)

	wf.keys = make([]float64, wf.Sound.FrameLen()*4)

	space := 16

	nkeys := len(hasmn)
	mj := len(wf.keys) / nkeys

	// dinky piano signal
	for i := 0; i < len(wf.keys); i += 2 {
		if i <= space {
			// marker for signal alignment
			wf.keys[i] = -0.999
			wf.keys[i+1] = -0.999
			continue
		} else if i >= (len(wf.keys) - space) {
			wf.keys[i] = -0.98
			wf.keys[i+1] = -0.98
			continue
		}

		key := i / mj
		j := i % mj
		if j <= space || (mj-j) >= (mj-space) {
			// spacing for major keys
			wf.keys[i] = -1
		} else if (j <= (mj/4) && hasmn[key].left) || (j >= mj-(mj/4) && hasmn[key].right) {
			// minor key
			wf.keys[i] = -0.3
		} else {
			// major key
			wf.keys[i] = 1
		}

		wf.keys[i+1] = -1
	}

	return wf
}

func (wf *Piano) KeyAt(ev touch.Event, sz size.Event) int {
	// piano is made up of 1024 points in width and half screen height
	x := int(ev.X / float32(sz.WidthPx) * float32(len(wf.keys)))
	y := (float64(sz.HeightPx)-float64(ev.Y))/float64(sz.HeightPx/2)*2 - 1
	if y < -1 || 1 < y {
		return -1
	}

	nkeys := len(hasmn)
	mj := len(wf.keys) / nkeys
	key := x / mj
	if key >= len(hasmn) {
		key = len(hasmn) - 1
	}
	if key < 0 {
		key = 0
	}
	j := x % mj
	if j <= (mj/4) && hasmn[key].left && y > -0.3 {
		return hasmn[key].pos - 1
	} else if j >= mj-(mj/4) && hasmn[key].right && y > -0.3 {
		return hasmn[key].pos + 1
	} else {
		return hasmn[key].pos
	}
}

func (wf *Piano) Prepare() {
	wf.Sound.Prepare()
	out := wf.Sound.Output()
	for i := range out {
		out[i] = wf.Sound.Amp(i) * wf.keys[wf.idx]
		wf.idx = (wf.idx + 1) % len(wf.keys)
	}
}