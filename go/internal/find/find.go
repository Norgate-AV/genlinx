package find

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

const ICSPPort = 1319

// Date holds numeric and text representations of a date.
type Date struct {
	Numeric string `json:"numeric"`
	Text    string `json:"text"`
}

// Device represents a discovered NetLinx device.
type Device struct {
	IP       string `json:"ip"`
	System   int    `json:"system"`
	Date     Date   `json:"date"`
	Time     string `json:"time"`
	Day      string `json:"day"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
	ID       string `json:"id"`
}

var days = map[byte]string{
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
	7: "Sunday",
}

// ---------------------------------------------------------------------------
// Packet parsing
// ---------------------------------------------------------------------------

// pkt is a stateful cursor-based packet parser that mirrors Packet.ts.
type pkt struct {
	data   []byte
	cursor int
}

func newPkt(payload []byte) (*pkt, error) {
	if len(payload) < 3 {
		return nil, fmt.Errorf("payload too short")
	}

	length := int(payload[2])
	start := 3
	end := start + length

	if len(payload) < end {
		return nil, fmt.Errorf("data length exceeds payload")
	}

	data := payload[start:end]
	if len(data) != length {
		return nil, fmt.Errorf("invalid data length")
	}

	return &pkt{data: data}, nil
}

func (p *pkt) at(offset int) byte {
	return p.data[p.cursor+offset]
}

func (p *pkt) int16BE(offset int) int16 {
	return int16(binary.BigEndian.Uint16(p.data[p.cursor+offset : p.cursor+offset+2]))
}

func (p *pkt) advance(n int) {
	p.cursor += n
}

func (p *pkt) remaining() []byte {
	return p.data[p.cursor:]
}

func (p *pkt) require(n int, field string) error {
	if p.cursor+n > len(p.data) {
		return fmt.Errorf("packet too short for %s (need %d, have %d)", field, n, len(p.data)-p.cursor)
	}

	return nil
}

// parse mirrors Packet.parse() exactly, including field ordering and cursor
// advancement.
func (p *pkt) parse() (*Device, error) {
	d := &Device{}

	// getSystem(): reads data[2..4] as Int16BE, advances 21.
	if err := p.require(21, "system"); err != nil {
		return nil, err
	}

	d.System = int(p.int16BE(2))
	p.advance(21)

	// getDateNumeric(): month[0], day[1], year[2..4], advances 4.
	if err := p.require(4, "date numeric"); err != nil {
		return nil, err
	}

	month := fmt.Sprintf("%02d", p.at(0))
	dayN := fmt.Sprintf("%02d", p.at(1))
	year := p.int16BE(2)
	d.Date.Numeric = fmt.Sprintf("%s/%s/%d", month, dayN, year)
	p.advance(4)

	// getTime(): hours[0], minutes[1], seconds[3], advances 3.
	// Note: data[3] is read for seconds but the cursor only moves 3 bytes,
	// leaving data[3] as data[0] for the next call (getDay). This faithfully
	// mirrors the TS implementation.
	if err := p.require(4, "time"); err != nil {
		return nil, err
	}

	hours := fmt.Sprintf("%02d", p.at(0))
	minutes := fmt.Sprintf("%02d", p.at(1))
	secondsByte := p.at(3) // read ahead; becomes data[0] after advance(3)
	d.Time = fmt.Sprintf("%s:%s:%02d", hours, minutes, secondsByte)
	p.advance(3)

	// getDay(): reads data[0] (= old data[3] from time block), advances 3.
	if err := p.require(3, "day"); err != nil {
		return nil, err
	}

	dayVal := p.at(0)
	if name, ok := days[dayVal]; ok {
		d.Day = name
	} else {
		d.Day = fmt.Sprintf("%d", dayVal)
	}

	p.advance(3)

	// getDateText(): read until triple-null '\0\0\0', advance past content + 3.
	tripleNull := []byte{0, 0, 0}
	endIdx := bytes.Index(p.remaining(), tripleNull)

	if endIdx < 0 {
		return nil, fmt.Errorf("packet missing date text terminator")
	}

	d.Date.Text = string(p.remaining()[:endIdx])
	p.advance(endIdx + 3)

	// getMac(): 6 bytes as colon-separated lower-case hex, advances 6.
	if err := p.require(6, "MAC"); err != nil {
		return nil, err
	}

	macParts := make([]string, 6)
	for i := 0; i < 6; i++ {
		macParts[i] = fmt.Sprintf("%02x", p.at(i))
	}

	d.MAC = strings.Join(macParts, ":")
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

// ---------------------------------------------------------------------------
// Discovery
// ---------------------------------------------------------------------------

// Discover listens on UDP port 1319 for ICSP identify-reply packets and
// returns deduplicated devices after the given timeout, sorted by IP.
func Discover(timeout time.Duration) ([]Device, error) {
	addr := &net.UDPAddr{Port: ICSPPort}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on UDP port %d: %w", ICSPPort, err)
	}

	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	var devices []Device
	buf := make([]byte, 4096)

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			// Any error (including deadline exceeded) ends the listen window.
			break
		}

		p, err := newPkt(buf[:n])
		if err != nil {
			continue
		}

		device, err := p.parse()
		if err != nil {
			continue
		}

		device.IP = remoteAddr.IP.String()
		devices = append(devices, *device)
	}

	// Deduplicate by MAC address (last-seen wins, matching TS Map behaviour).
	seen := make(map[string]Device, len(devices))
	for _, d := range devices {
		seen[d.MAC] = d
	}

	result := make([]Device, 0, len(seen))
	for _, d := range seen {
		result = append(result, d)
	}

	// Sort by IP ascending (mirrors the TS sort comparator).
	sort.Slice(result, func(i, j int) bool {
		return result[i].IP < result[j].IP
	})

	return result, nil
}

// ---------------------------------------------------------------------------
// Output
// ---------------------------------------------------------------------------

// PrintJSON outputs devices as pretty-printed JSON.
func PrintJSON(devices []Device) error {
	b, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

// PrintTable outputs devices as an aligned text table.
func PrintTable(devices []Device) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "Discovered Devices")
	fmt.Fprintln(w, "IP Address\tSystem\tDate\tTime\tMAC Address\tHostname\tID")

	for _, d := range devices {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n",
			d.IP, d.System, d.Date.Text, d.Time, d.MAC, d.Hostname, d.ID)
	}

	w.Flush()
}
