// Package playback provides audio playback support
package playback

import (
	"net/http"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Player struct {
	URL        string
	ctrl       *beep.Ctrl
	sampleRate beep.SampleRate
}

func NewPlayer() *Player {
	rate := beep.SampleRate(44100)
	speaker.Init(rate, rate.N(time.Second/10))
	return &Player{sampleRate: rate}
}

func (p *Player) Play(url string) tea.Cmd {
	p.URL = url
	return func() tea.Msg {
		p.Stop()

		res, err := http.Get(url)
		if err != nil {
			return PlayerErrorMsg(err.Error())
		}

		streamer, format, err := mp3.Decode(res.Body)
		if err != nil {
			res.Body.Close()
			return PlayerErrorMsg(err.Error())
		}

		resampled := beep.Resample(4, format.SampleRate, p.sampleRate, streamer)

		speaker.Lock()
		if p.ctrl == nil {
			p.ctrl = &beep.Ctrl{Streamer: resampled, Paused: false}
			speaker.Play(p.ctrl)
		} else {
			p.ctrl.Streamer = resampled
			p.ctrl.Paused = false
		}
		speaker.Unlock()

		speaker.Play(p.ctrl)
		return PlayerStartedMsg("Player Playing")
	}
}

func (p *Player) Stop() {
	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = true
		if p.ctrl.Streamer != nil {
			if closer, ok := p.ctrl.Streamer.(beep.StreamCloser); ok {
				closer.Close()
			}
			p.ctrl.Streamer = nil
		}
		p.ctrl.Streamer = nil
		speaker.Unlock()
	}
}

func (p *Player) IsPlaying() bool {
	speaker.Lock()
	defer speaker.Unlock()
	return p.ctrl != nil && p.ctrl.Streamer != nil
}

type (
	PlayerErrorMsg   string
	PlayerStartedMsg string
)
