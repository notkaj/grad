package playback

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"

	"github.com/gopxl/beep"
)

// ffmpegStreamer reads raw PCM data from ffmpeg's stdout and converts it to beep samples.
type ffmpegStreamer struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	err    error
}

func newFFmpegStreamer(url string, sampleRate beep.SampleRate) (*ffmpegStreamer, error) {
	// We request raw s16le PCM to avoid WAV header seeking issues over pipes.
	cmd := exec.Command("ffmpeg",
		"-i", url,
		"-f", "s16le",
		"-acodec", "pcm_s16le",
		"-ac", "2",
		"-ar", fmt.Sprintf("%d", sampleRate),
		"-loglevel", "quiet",
		"-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &ffmpegStreamer{
		cmd:    cmd,
		stdout: stdout,
	}, nil
}

func (s *ffmpegStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.err != nil {
		return 0, false
	}

	buf := make([]byte, len(samples)*4) // 2 channels * 2 bytes (s16le)
	nBytes, err := io.ReadFull(s.stdout, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		s.err = err
		return 0, false
	}

	// Calculate how many full samples (4 bytes each) we actually read
	n = nBytes / 4

	for i := 0; i < n; i++ {
		left := int16(binary.LittleEndian.Uint16(buf[i*4 : i*4+2]))
		right := int16(binary.LittleEndian.Uint16(buf[i*4+2 : i*4+4]))
		samples[i][0] = float64(left) / 32768.0
		samples[i][1] = float64(right) / 32768.0
	}

	if n == 0 && (err == io.EOF || err == io.ErrUnexpectedEOF) {
		return 0, false
	}

	return n, true
}

func (s *ffmpegStreamer) Err() error {
	return s.err
}

func (s *ffmpegStreamer) Close() error {
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	if s.stdout != nil {
		return s.stdout.Close()
	}
	return nil
}
