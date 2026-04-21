package playback

import (
	"fmt"
	"io"

	"github.com/gopxl/beep"
	"github.com/llehouerou/go-aac"
)

// nativeAACStreamer decodes raw AAC (ADTS) streams using a pure-Go decoder.
type nativeAACStreamer struct {
	closer io.Closer
	reader io.Reader
	dec    *aac.Decoder
	format beep.Format

	buf [][2]float64
	pos int
	err error
}

func aacDecode(reader io.ReadCloser) (*nativeAACStreamer, error) {
	dec := aac.NewDecoder()

	s := &nativeAACStreamer{
		closer: reader,
		reader: reader,
		dec:    dec,
	}

	// 1. Read the first frame to initialize the decoder
	frame, err := s.readNextADTSFrame()
	if err != nil {
		reader.Close()
		return nil, fmt.Errorf("failed to read first ADTS frame: %v", err)
	}

	sampleRate, channels, err := dec.SimpleInit(frame)
	if err != nil {
		reader.Close()
		return nil, fmt.Errorf("aac initialization failed: %v", err)
	}

	s.format = beep.Format{
		SampleRate:  beep.SampleRate(sampleRate),
		NumChannels: int(channels),
		Precision:   2,
	}

	// 2. Decode that first frame so we have some initial samples
	if err := s.decodeFrame(frame); err != nil {
		reader.Close()
		return nil, fmt.Errorf("first frame decode failed: %v", err)
	}

	return s, nil
}

func (s *nativeAACStreamer) readNextADTSFrame() ([]byte, error) {
	// ADTS header is 7-9 bytes.
	// Syncword is 0xFFF (12 bits)
	header := make([]byte, 7)

	// Find syncword
	for {
		if _, err := io.ReadFull(s.reader, header[:1]); err != nil {
			return nil, err
		}
		if header[0] == 0xFF {
			if _, err := io.ReadFull(s.reader, header[1:2]); err != nil {
				return nil, err
			}
			if (header[1] & 0xF0) == 0xF0 {
				// Found 0xFFF!
				break
			}
		}
	}

	// Read remaining 5 bytes of fixed header
	if _, err := io.ReadFull(s.reader, header[2:7]); err != nil {
		return nil, err
	}

	// Extract Frame Length (13 bits starting at bit 30)
	// byte 3: last 2 bits
	// byte 4: all 8 bits
	// byte 5: first 3 bits
	frameLen := int(header[3]&0x03)<<11 | int(header[4])<<3 | int(header[5])>>5

	if frameLen < 7 || frameLen > 8192 {
		return nil, fmt.Errorf("invalid ADTS frame length: %d", frameLen)
	}

	frame := make([]byte, frameLen)
	copy(frame, header)

	if _, err := io.ReadFull(s.reader, frame[7:]); err != nil {
		return nil, err
	}

	return frame, nil
}

func (s *nativeAACStreamer) decodeFrame(frameData []byte) error {
	pcm, err := s.dec.DecodeInt16(frameData)
	if err != nil {
		return err
	}

	// Convert interleaved int16 PCM to Beep's float64 samples
	newSamples := make([][2]float64, len(pcm)/s.format.NumChannels)
	for i := range newSamples {
		if s.format.NumChannels == 2 {
			newSamples[i][0] = float64(pcm[i*2]) / 32768.0
			newSamples[i][1] = float64(pcm[i*2+1]) / 32768.0
		} else {
			val := float64(pcm[i]) / 32768.0
			newSamples[i][0] = val
			newSamples[i][1] = val
		}
	}

	s.buf = append(s.buf, newSamples...)
	return nil
}

func (s *nativeAACStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.err != nil {
		return 0, false
	}

	for n < len(samples) {
		// If buffer is empty, try to fill it
		if s.pos >= len(s.buf) {
			s.buf = s.buf[:0]
			s.pos = 0

			frame, err := s.readNextADTSFrame()
			if err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					s.err = err
				}
				if n > 0 {
					return n, true
				}
				return 0, false
			}

			if err := s.decodeFrame(frame); err != nil {
				s.err = err
				if n > 0 {
					return n, true
				}
				return 0, false
			}
		}

		// Copy from buffer to samples
		copyCount := copy(samples[n:], s.buf[s.pos:])
		n += copyCount
		s.pos += copyCount
	}

	return n, true
}

func (s *nativeAACStreamer) Err() error {
	return s.err
}

func (s *nativeAACStreamer) Close() error {
	return s.closer.Close()
}
