package find

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var days = [8]string{
	0: "",
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
	7: "Sunday",
}

// minFixedLen is the number of bytes consumed by all fixed-length fields
// (system 21 + date-numeric 4 + time/day 7 + MAC 6).
const minFixedLen = 38

type Packet struct {
	data   []byte
	cursor int
}

func NewPacket(payload []byte) (*Packet, error) {
	if len(payload) < 3 {
		return nil, fmt.Errorf("payload too short")
	}

	length := int(payload[2])
	start := 3
	end := start + length

	if len(payload) < end {
		return nil, fmt.Errorf("data length exceeds payload")
	}

	return &Packet{
		data: payload[start:end],
	}, nil
}

func (p *Packet) at(offset int) byte {
	return p.data[p.cursor+offset]
}

func (p *Packet) int16BE(offset int) int16 {
	return int16(binary.BigEndian.Uint16(p.data[p.cursor+offset : p.cursor+offset+2]))
}

func (p *Packet) advance(n int) {
	p.cursor += n
}

func (p *Packet) remaining() []byte {
	return p.data[p.cursor:]
}

func (p *Packet) Parse() (*Device, error) {
	if len(p.data) < minFixedLen {
		return nil, fmt.Errorf("packet too short (have %d, need at least %d bytes)", len(p.data), minFixedLen)
	}

	d := &Device{}

	// getSystem(): reads data[2..4] as Int16BE, advances 21.
	d.System = int(p.int16BE(2))
	p.advance(21)

	// getDateNumeric(): month[0], day[1], year[2..4], advances 4.
	d.Date.Numeric = fmt.Sprintf("%02d/%02d/%d", p.at(0), p.at(1), p.int16BE(2))
	p.advance(4)

	// getTime(): hours[0], minutes[1], seconds[3], advances 3.
	// Note: data[3] is read for seconds but the cursor only moves 3 bytes,
	// leaving data[3] as data[0] for the next call (getDay).
	d.Time = fmt.Sprintf("%02d:%02d:%02d", p.at(0), p.at(1), p.at(3))
	p.advance(3)

	// getDay(): reads data[0] (= old data[3] from time block), advances 3.
	dayVal := p.at(0)
	if int(dayVal) < len(days) && days[dayVal] != "" {
		d.Day = days[dayVal]
	} else {
		d.Day = fmt.Sprintf("%d", dayVal)
	}

	p.advance(3)

	// getDateText(): read until triple-null '\0\0\0', advance past content + 3.
	endIdx := bytes.Index(p.remaining(), []byte{0, 0, 0})
	if endIdx < 0 {
		return nil, fmt.Errorf("packet missing date text terminator")
	}

	d.Date.Text = string(p.remaining()[:endIdx])
	p.advance(endIdx + 3)

	// getMac(): 6 bytes as colon-separated lower-case hex, advances 6.
	if p.cursor+6 > len(p.data) {
		return nil, fmt.Errorf("packet too short for MAC")
	}

	d.MAC = fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		p.at(0), p.at(1), p.at(2), p.at(3), p.at(4), p.at(5))
	p.advance(6)

	// getHostname(): read until '\0', advance past content + 1.
	nullIdx := bytes.IndexByte(p.remaining(), 0)
	if nullIdx < 0 {
		return nil, fmt.Errorf("packet missing hostname terminator")
	}

	d.Hostname = string(p.remaining()[:nullIdx])
	p.advance(nullIdx + 1)

	// getId(): read until '\0' (or end of buffer).
	nullIdx = bytes.IndexByte(p.remaining(), 0)
	if nullIdx < 0 {
		d.ID = string(p.remaining())
	} else {
		d.ID = string(p.remaining()[:nullIdx])
	}

	return d, nil
}
