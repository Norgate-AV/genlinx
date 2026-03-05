package find

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ---------------------------------------------------------------------------
// Suite
// ---------------------------------------------------------------------------

type FindTestSuite struct {
	suite.Suite
}

func TestFindSuite(t *testing.T) {
	suite.Run(t, new(FindTestSuite))
}

// ---------------------------------------------------------------------------
// Test payload
// ---------------------------------------------------------------------------

// testPayload builds a valid ICSP identify-reply payload that decodes to a
// known Device. Wire layout (all offsets relative to the data slice, i.e.
// payload[3:]):
//
//	[0..20]   system block   (21 bytes; int16BE system id at [2..4])
//	[21..24]  date numeric   (4 bytes; month, day, year int16BE)
//	[25..28]  time block     (4 bytes; h, m, pad, s — cursor advances only 3)
//	[28..30]  day block      (3 bytes; data[28] is shared with time seconds)
//	[31..]    date text      (until \x00\x00\x00)
//	          MAC            (6 bytes)
//	          hostname       (until \x00)
//	          id             (until \x00 or EOF)
func testPayload() []byte {
	var data []byte

	// system block: 21 bytes, system=1 encoded as int16BE at offsets [2..4].
	sys := make([]byte, 21)
	sys[2] = 0x00
	sys[3] = 0x01
	data = append(data, sys...)

	// date numeric: month=3, day=5, year=2026 (0x07EA).
	data = append(data, 3, 5, 0x07, 0xEA)

	// time block: hours=10, minutes=30, pad=0, seconds/dayVal=4.
	// getTime() advances only 3, so data[28]=4 is re-read by getDay().
	data = append(data, 10, 30, 0, 4)

	// day block padding: data[28]=4 already set above; two more bytes required.
	data = append(data, 0, 0)

	// date text: arbitrary string terminated by triple-null.
	data = append(data, append([]byte("Thursday, 05 March 2026"), 0, 0, 0)...)

	// MAC: 6 bytes → "00:60:9f:01:02:03"
	data = append(data, 0x00, 0x60, 0x9F, 0x01, 0x02, 0x03)

	// hostname: null-terminated.
	data = append(data, append([]byte("NX-3200"), 0)...)

	// id: no trailing null (EOF path).
	data = append(data, []byte("ICSPS70B")...)

	// prefix with [0x00, 0x00, length].
	return append([]byte{0x00, 0x00, byte(len(data))}, data...)
}

// ---------------------------------------------------------------------------
// newPkt
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestNewPkt_ValidPayload() {
	p, err := NewPacket(testPayload())
	s.Require().NoError(err)
	s.Require().NotNil(p)
}

func (s *FindTestSuite) TestNewPkt_TooShort() {
	_, err := NewPacket([]byte{0x00, 0x00})
	s.Require().Error(err)
	s.Contains(err.Error(), "too short")
}

func (s *FindTestSuite) TestNewPkt_DataLengthExceedsPayload() {
	// payload[2] claims 100 bytes of data but only 2 follow.
	_, err := NewPacket([]byte{0x00, 0x00, 100, 0x01, 0x02})
	s.Require().Error(err)
	s.Contains(err.Error(), "data length exceeds payload")
}

// ---------------------------------------------------------------------------
// parse — happy path
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestParse_FullPacket() {
	p, err := NewPacket(testPayload())
	s.Require().NoError(err)

	d, err := p.Parse()
	s.Require().NoError(err)
	s.Require().NotNil(d)

	s.Equal(1, d.System)
	s.Equal("03/05/2026", d.Date.Numeric)
	s.Equal("10:30:04", d.Time)
	s.Equal("Thursday", d.Day)
	s.Equal("Thursday, 05 March 2026", d.Date.Text)
	s.Equal("00:60:9f:01:02:03", d.MAC)
	s.Equal("NX-3200", d.Hostname)
	s.Equal("ICSPS70B", d.ID)
}

// ---------------------------------------------------------------------------
// parse — error paths
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestParse_TruncatedAtSystem() {
	// Only 3 bytes of data — can't satisfy the 21-byte system requirement.
	payload := []byte{0x00, 0x00, 3, 0x01, 0x02, 0x03}
	p, err := NewPacket(payload)
	s.Require().NoError(err)

	_, err = p.Parse()
	s.Require().Error(err)
	s.Contains(err.Error(), "too short")
}

func (s *FindTestSuite) TestParse_MissingDateTextTerminator() {
	var data []byte
	sys := make([]byte, 21)
	sys[3] = 0x01
	data = append(data, sys...)
	data = append(data, 3, 5, 0x07, 0xEA)
	data = append(data, 10, 30, 0, 4)
	data = append(data, 0, 0)
	data = append(data, []byte("No triple null here")...)

	payload := append([]byte{0x00, 0x00, byte(len(data))}, data...)
	p, err := NewPacket(payload)
	s.Require().NoError(err)

	_, err = p.Parse()
	s.Require().Error(err)
	s.Contains(err.Error(), "date text terminator")
}

func (s *FindTestSuite) TestParse_MissingHostnameTerminator() {
	var data []byte
	sys := make([]byte, 21)
	sys[3] = 0x01
	data = append(data, sys...)
	data = append(data, 3, 5, 0x07, 0xEA)
	data = append(data, 10, 30, 0, 4)
	data = append(data, 0, 0)
	data = append(data, append([]byte("Thursday, 05 March 2026"), 0, 0, 0)...)
	data = append(data, 0x00, 0x60, 0x9F, 0x01, 0x02, 0x03)
	data = append(data, []byte("NX-3200")...) // no null terminator

	payload := append([]byte{0x00, 0x00, byte(len(data))}, data...)
	p, err := NewPacket(payload)
	s.Require().NoError(err)

	_, err = p.Parse()
	s.Require().Error(err)
	s.Contains(err.Error(), "hostname terminator")
}

// ---------------------------------------------------------------------------
// parse — edge cases
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestParse_UnknownDayFallsBackToNumeric() {
	// dayVal=0 is not in the days map; should produce the string "0".
	var data []byte
	sys := make([]byte, 21)
	sys[3] = 0x01
	data = append(data, sys...)
	data = append(data, 3, 5, 0x07, 0xEA)
	data = append(data, 10, 30, 0, 0) // seconds/dayVal = 0 (not in map)
	data = append(data, 0, 0)
	data = append(data, append([]byte("Unknown Day"), 0, 0, 0)...)
	data = append(data, 0x00, 0x60, 0x9F, 0x01, 0x02, 0x03)
	data = append(data, append([]byte("host"), 0)...)
	data = append(data, []byte("id")...)

	payload := append([]byte{0x00, 0x00, byte(len(data))}, data...)
	p, err := NewPacket(payload)
	s.Require().NoError(err)

	d, err := p.Parse()
	s.Require().NoError(err)
	s.Equal("0", d.Day)
}

func (s *FindTestSuite) TestParse_IDWithNullTerminator() {
	// getId() should strip a trailing null when one is present.
	var data []byte
	sys := make([]byte, 21)
	sys[3] = 0x01
	data = append(data, sys...)
	data = append(data, 3, 5, 0x07, 0xEA)
	data = append(data, 10, 30, 0, 4)
	data = append(data, 0, 0)
	data = append(data, append([]byte("Thursday, 05 March 2026"), 0, 0, 0)...)
	data = append(data, 0x00, 0x60, 0x9F, 0x01, 0x02, 0x03)
	data = append(data, append([]byte("NX-3200"), 0)...)
	data = append(data, append([]byte("ICSPS70B"), 0)...) // trailing null

	payload := append([]byte{0x00, 0x00, byte(len(data))}, data...)
	p, err := NewPacket(payload)
	s.Require().NoError(err)

	d, err := p.Parse()
	s.Require().NoError(err)
	s.Equal("ICSPS70B", d.ID)
}

// ---------------------------------------------------------------------------
// isTimeout
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestIsTimeout_TrueForTimeoutError() {
	s.True(isTimeout(&fakeNetErr{timeout: true}))
}

func (s *FindTestSuite) TestIsTimeout_FalseForNonTimeoutNetError() {
	s.False(isTimeout(&fakeNetErr{timeout: false}))
}

func (s *FindTestSuite) TestIsTimeout_FalseForNonNetError() {
	s.False(isTimeout(errors.New("some other error")))
}

// fakeNetErr is a minimal net.Error implementation for testing isTimeout.
type fakeNetErr struct{ timeout bool }

func (e *fakeNetErr) Error() string   { return "fake net error" }
func (e *fakeNetErr) Timeout() bool   { return e.timeout }
func (e *fakeNetErr) Temporary() bool { return false }

// ---------------------------------------------------------------------------
// PrintJSON
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestPrintJSON_EmptySlice() {
	err := PrintJSON([]Device{})
	s.NoError(err)
}

func (s *FindTestSuite) TestPrintJSON_SingleDevice() {
	err := PrintJSON([]Device{
		{IP: "192.168.1.1", System: 1, MAC: "00:60:9f:01:02:03", Hostname: "NX-3200", ID: "ICSPS70B"},
	})
	s.NoError(err)
}

// ---------------------------------------------------------------------------
// PrintTable
// ---------------------------------------------------------------------------

func (s *FindTestSuite) TestPrintTable_EmptySlice() {
	s.NotPanics(func() {
		PrintTable([]Device{})
	})
}

func (s *FindTestSuite) TestPrintTable_WithDevices() {
	s.NotPanics(func() {
		PrintTable([]Device{
			{IP: "192.168.1.1", System: 1, MAC: "00:60:9f:01:02:03", Hostname: "NX-3200", ID: "ICSPS70B"},
			{IP: "192.168.1.2", System: 2, MAC: "00:60:9f:04:05:06", Hostname: "NX-4200", ID: "ICSPS70C"},
		})
	})
}
