// Package playback provides audio playback support
package playback

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/vorbis"
)

type Player struct {
	mu         sync.Mutex
	URL        string
	mixer      *beep.Mixer
	volume     *effects.Volume
	current    io.Closer
	cancel     context.CancelFunc
	sampleRate beep.SampleRate
}

func NewPlayer() *Player {
	rate := beep.SampleRate(44100)
	speaker.Init(rate, rate.N(time.Second/10))
	mixer := &beep.Mixer{}
	volume := &effects.Volume{
		Streamer: mixer,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speaker.Play(volume)
	return &Player{sampleRate: rate, mixer: mixer, volume: volume}
}

func (p *Player) Play(url string) tea.Cmd {
	return func() tea.Msg {
		p.Stop()

		p.mu.Lock()
		p.URL = url
		ctx, cancel := context.WithCancel(context.Background())
		p.cancel = cancel
		p.mu.Unlock()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return PlayerErrorMsg(err.Error())
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return PlayerErrorMsg(err.Error())
		}

		contentType := strings.ToLower(res.Header.Get("Content-Type"))
		var (
			streamer beep.StreamSeekCloser
			format   beep.Format
		)

		isOgg := strings.Contains(contentType, "ogg") || strings.Contains(contentType, "vorbis")
		isMp3 := strings.Contains(contentType, "mpeg")
		isAAC := strings.Contains(contentType, "aac") || strings.Contains(contentType, "m4a")
		isHLS := strings.Contains(contentType, "mpegurl") || strings.Contains(contentType, "apple.mpegurl")

		if isOgg {
			streamer, format, err = vorbis.Decode(res.Body)
		} else if isMp3 {
			streamer, format, err = mp3.Decode(res.Body)
		} else if isAAC {
			nativeStreamer, nErr := aacDecode(res.Body)
			if nErr == nil {
				resampled := beep.Resample(4, nativeStreamer.format.SampleRate, p.sampleRate, nativeStreamer)
				if p.setupPlayer(ctx, resampled, nativeStreamer) {
					return PlayerStartedMsg("aac (native)")
				}
				return nil
			}
		}

		if isHLS || (err != nil && (isOgg || isMp3)) || streamer == nil {
			res.Body.Close()
			return p.playWithFFmpeg(ctx, url)
		}

		resampled := beep.Resample(4, format.SampleRate, p.sampleRate, streamer)
		if p.setupPlayer(ctx, resampled, streamer) {
			return PlayerStartedMsg("native")
		}
		return nil
	}
}

func (p *Player) setupPlayer(ctx context.Context, streamer beep.Streamer, closer io.Closer) bool {
	wrapped := &closeWrapper{Streamer: streamer, closer: closer}

	p.mu.Lock()
	defer p.mu.Unlock()
	if ctx.Err() != nil {
		wrapped.Close()
		return false
	}

	speaker.Lock()
	p.current = wrapped
	p.mixer.Add(wrapped)
	speaker.Unlock()
	return true
}

func (p *Player) playWithFFmpeg(ctx context.Context, url string) tea.Msg {
	s, err := newFFmpegStreamer(ctx, url, p.sampleRate)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return PlayerErrorMsg("ffmpeg error: " + err.Error())
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if ctx.Err() != nil {
		s.Close()
		return nil
	}

	speaker.Lock()
	defer speaker.Unlock()
	p.current = s
	p.mixer.Add(s)

	return PlayerStartedMsg("ffmpeg")
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	if p.current != nil {
		p.current.Close()
		p.current = nil
	}
	p.mixer.Clear()
}

func (p *Player) IsPlaying() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.current != nil
}

func (p *Player) SetVolume(v float64) {
	speaker.Lock()
	defer speaker.Unlock()
	p.volume.Volume = v
}

func (p *Player) Volume() float64 {
	speaker.Lock()
	defer speaker.Unlock()
	return p.volume.Volume
}

func (p *Player) ToggleMute() {
	speaker.Lock()
	defer speaker.Unlock()
	p.volume.Silent = !p.volume.Silent
}

func (p *Player) IsMuted() bool {
	speaker.Lock()
	defer speaker.Unlock()
	return p.volume.Silent
}

type (
	PlayerErrorMsg   string
	PlayerStartedMsg string
)

type closeWrapper struct {
	beep.Streamer
	closer io.Closer
}

func (w *closeWrapper) Close() error {
	return w.closer.Close()
}
