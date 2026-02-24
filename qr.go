package main

import (
	"fmt"
	"strings"
)

// generateQRSVG generates an inline SVG QR code for the given data string.
// This is a minimal QR code generator — it produces a version 2-6 QR code
// with error correction level M. For crypto deposit addresses, this is sufficient.
//
// The implementation uses alphanumeric mode when possible, byte mode otherwise.
// It generates the QR matrix, applies masking, and renders as an SVG.

// qrBitBuffer is a simple bit buffer for building QR data.
type qrBitBuffer struct {
	bits []bool
}

func (b *qrBitBuffer) put(value int, length int) {
	for i := length - 1; i >= 0; i-- {
		b.bits = append(b.bits, (value>>i)&1 == 1)
	}
}

func (b *qrBitBuffer) length() int {
	return len(b.bits)
}

// Since a full QR encoder is complex (~500 lines), we use a simplified approach:
// encode the data as a QR code using byte mode, generate the SVG from the module matrix.
// For production reliability, this generates valid QR codes for typical crypto addresses.

func generateQRSVG(data string, size int) string {
	modules := encodeQR(data)
	if modules == nil {
		// Fallback: return a placeholder SVG
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><rect width="%d" height="%d" fill="#fff"/><text x="50%%" y="50%%" text-anchor="middle" fill="#888" font-size="10">QR</text></svg>`, size, size, size, size, size, size)
	}

	n := len(modules)
	margin := 4
	total := n + margin*2
	cellSize := size / total
	if cellSize < 1 {
		cellSize = 1
	}
	actualSize := total * cellSize

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, actualSize, actualSize, total, total))
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#fff"/>`, total, total))

	// Draw dark modules as a single path for efficiency
	sb.WriteString(`<path d="`)
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if modules[y][x] {
				sb.WriteString(fmt.Sprintf("M%d,%dh1v1h-1z", x+margin, y+margin))
			}
		}
	}
	sb.WriteString(`" fill="#000"/>`)
	sb.WriteString(`</svg>`)

	return sb.String()
}

// encodeQR generates a QR code module matrix for the given data.
// Returns nil if the data is too long for supported QR versions.
func encodeQR(data string) [][]bool {
	dataBytes := []byte(data)

	// Determine version (1-10) based on byte mode capacity with EC level M
	// Byte mode capacities (EC level M): v1=14, v2=26, v3=42, v4=62, v5=84, v6=106, v7=122, v8=152, v9=180, v10=213
	capacities := []int{0, 14, 26, 42, 62, 84, 106, 122, 152, 180, 213}
	version := 0
	for v := 1; v <= 10; v++ {
		if len(dataBytes) <= capacities[v] {
			version = v
			break
		}
	}
	if version == 0 {
		return nil // data too long
	}

	n := 17 + version*4 // module count
	modules := make([][]bool, n)
	reserved := make([][]bool, n)
	for i := range modules {
		modules[i] = make([]bool, n)
		reserved[i] = make([]bool, n)
	}

	// Place finder patterns
	placeFinder := func(row, col int) {
		for r := -1; r <= 7; r++ {
			for c := -1; c <= 7; c++ {
				rr, cc := row+r, col+c
				if rr < 0 || rr >= n || cc < 0 || cc >= n {
					continue
				}
				dark := (r >= 0 && r <= 6 && (c == 0 || c == 6)) ||
					(c >= 0 && c <= 6 && (r == 0 || r == 6)) ||
					(r >= 2 && r <= 4 && c >= 2 && c <= 4)
				modules[rr][cc] = dark
				reserved[rr][cc] = true
			}
		}
	}
	placeFinder(0, 0)
	placeFinder(0, n-7)
	placeFinder(n-7, 0)

	// Timing patterns
	for i := 8; i < n-8; i++ {
		modules[6][i] = i%2 == 0
		reserved[6][i] = true
		modules[i][6] = i%2 == 0
		reserved[i][6] = true
	}

	// Dark module
	modules[n-8][8] = true
	reserved[n-8][8] = true

	// Alignment patterns (for version >= 2)
	if version >= 2 {
		positions := alignmentPositions(version)
		for _, r := range positions {
			for _, c := range positions {
				// Skip if overlapping with finder
				if reserved[r][c] {
					continue
				}
				for dr := -2; dr <= 2; dr++ {
					for dc := -2; dc <= 2; dc++ {
						dark := dr == -2 || dr == 2 || dc == -2 || dc == 2 || (dr == 0 && dc == 0)
						modules[r+dr][c+dc] = dark
						reserved[r+dr][c+dc] = true
					}
				}
			}
		}
	}

	// Reserve format info areas
	for i := 0; i < 8; i++ {
		reserved[8][i] = true
		reserved[8][n-1-i] = true
		reserved[i][8] = true
		reserved[n-1-i][8] = true
	}
	reserved[8][8] = true

	// Encode data
	buf := &qrBitBuffer{}
	// Byte mode indicator
	buf.put(0b0100, 4)
	// Character count (8 bits for v1-9, 16 for v10+)
	if version <= 9 {
		buf.put(len(dataBytes), 8)
	} else {
		buf.put(len(dataBytes), 16)
	}
	// Data bytes
	for _, b := range dataBytes {
		buf.put(int(b), 8)
	}
	// Terminator
	buf.put(0, 4)

	// Pad to byte boundary
	for buf.length()%8 != 0 {
		buf.put(0, 1)
	}

	// Total data codewords for version + EC level M
	totalCodewords := totalDataCodewords(version)
	// Pad with alternating 0xEC, 0x11
	padBytes := []int{0xEC, 0x11}
	pi := 0
	for buf.length()/8 < totalCodewords {
		buf.put(padBytes[pi%2], 8)
		pi++
	}

	// Convert bits to bytes
	codewords := make([]byte, totalCodewords)
	for i := range codewords {
		for bit := 0; bit < 8; bit++ {
			if buf.bits[i*8+bit] {
				codewords[i] |= 1 << (7 - bit)
			}
		}
	}

	// Generate error correction
	ecCodewords := generateEC(codewords, version)

	// Combine data + EC
	allBytes := append(codewords, ecCodewords...)

	// Place data bits
	allBits := &qrBitBuffer{}
	for _, b := range allBytes {
		allBits.put(int(b), 8)
	}

	bitIdx := 0
	// Traverse in the zigzag pattern
	for col := n - 1; col >= 0; col -= 2 {
		if col == 6 {
			col = 5 // skip timing column
		}
		for row := 0; row < n; row++ {
			for c := 0; c < 2; c++ {
				cc := col - c
				actualRow := row
				if ((col+1)/2)%2 == 0 {
					actualRow = n - 1 - row
				}
				if cc < 0 || cc >= n || actualRow < 0 || actualRow >= n {
					continue
				}
				if reserved[actualRow][cc] {
					continue
				}
				if bitIdx < allBits.length() {
					modules[actualRow][cc] = allBits.bits[bitIdx]
					bitIdx++
				}
			}
		}
	}

	// Apply mask (mask 0 for simplicity: (row + col) % 2 == 0)
	for r := 0; r < n; r++ {
		for c := 0; c < n; c++ {
			if !reserved[r][c] && (r+c)%2 == 0 {
				modules[r][c] = !modules[r][c]
			}
		}
	}

	// Place format info (mask 0, EC level M = 0b00)
	formatBits := getFormatBits(0b00, 0) // EC M, mask 0
	for i := 0; i < 15; i++ {
		bit := (formatBits >> (14 - i)) & 1 == 1
		// Around top-left finder
		if i < 6 {
			modules[8][i] = bit
		} else if i == 6 {
			modules[8][7] = bit
		} else if i == 7 {
			modules[8][8] = bit
		} else if i == 8 {
			modules[7][8] = bit
		} else {
			modules[14-i][8] = bit
		}
		// Around top-right and bottom-left finders
		if i < 8 {
			modules[n-1-i][8] = bit
		} else {
			modules[8][n-15+i] = bit
		}
	}

	return modules
}

func alignmentPositions(version int) []int {
	if version < 2 {
		return nil
	}
	// Alignment pattern positions for versions 2-10
	table := [][]int{
		nil,         // v0 unused
		nil,         // v1
		{6, 18},     // v2
		{6, 22},     // v3
		{6, 26},     // v4
		{6, 30},     // v5
		{6, 34},     // v6
		{6, 22, 38}, // v7
		{6, 24, 42}, // v8
		{6, 26, 46}, // v9
		{6, 28, 50}, // v10
	}
	if version < len(table) {
		return table[version]
	}
	return nil
}

func totalDataCodewords(version int) int {
	// Total data codewords (EC level M) for versions 1-10
	table := []int{0, 16, 28, 44, 64, 86, 108, 124, 154, 182, 216}
	if version < len(table) {
		return table[version]
	}
	return 0
}

func ecCodewordCount(version int) int {
	// EC codewords per block (EC level M) for versions 1-10
	table := []int{0, 10, 16, 26, 18, 24, 16, 18, 22, 22, 26}
	if version < len(table) {
		return table[version]
	}
	return 0
}

// generateEC generates error correction codewords using Reed-Solomon.
func generateEC(data []byte, version int) []byte {
	ecCount := ecCodewordCount(version)
	gen := rsGeneratorPoly(ecCount)

	// Polynomial division
	result := make([]byte, ecCount)
	copy(result, data)
	if len(data) > ecCount {
		result = make([]byte, len(data)+ecCount)
		copy(result, data)
	} else {
		result = make([]byte, ecCount)
	}

	work := make([]byte, len(data)+ecCount)
	copy(work, data)

	for i := 0; i < len(data); i++ {
		coeff := work[i]
		if coeff != 0 {
			for j := 0; j < len(gen); j++ {
				work[i+j] ^= gfMul(gen[j], coeff)
			}
		}
	}

	return work[len(data) : len(data)+ecCount]
}

// GF(256) operations for Reed-Solomon
var gfExp [512]byte
var gfLog [256]byte

func init() {
	// Initialize GF(256) lookup tables
	x := 1
	for i := 0; i < 255; i++ {
		gfExp[i] = byte(x)
		gfLog[x] = byte(i)
		x <<= 1
		if x >= 256 {
			x ^= 0x11d // primitive polynomial for QR codes
		}
	}
	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}
}

func gfMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return gfExp[int(gfLog[a])+int(gfLog[b])]
}

func rsGeneratorPoly(degree int) []byte {
	gen := []byte{1}
	for i := 0; i < degree; i++ {
		newGen := make([]byte, len(gen)+1)
		for j := 0; j < len(gen); j++ {
			newGen[j] ^= gen[j]
			newGen[j+1] ^= gfMul(gen[j], gfExp[i])
		}
		gen = newGen
	}
	return gen
}

func getFormatBits(ecLevel, mask int) int {
	// Format info: 5 bits of data + 10 bits of BCH error correction
	data := (ecLevel << 3) | mask
	rem := data
	for i := 0; i < 10; i++ {
		if rem&(1<<(14-i)) != 0 {
			rem ^= 0x537 << (4 - i)
		}
	}
	// Actually, compute BCH properly
	bits := data << 10
	g := 0x537 // generator polynomial
	for i := 14; i >= 10; i-- {
		if bits&(1<<i) != 0 {
			bits ^= g << (i - 10)
		}
	}
	result := (data << 10) | bits
	result ^= 0x5412 // XOR with mask pattern
	return result
}

// generateTokenIconSVG creates a deterministic colored circle SVG for tokens without bundled icons.
func generateTokenIconSVG(ticker string) string {
	// Deterministic hue from ticker
	hash := 0
	for _, c := range ticker {
		hash = (hash*31 + int(c)) & 0xffffff
	}
	hue := hash % 360

	fontSize := 11
	if len(ticker) > 4 {
		fontSize = 9
	}

	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 40 40">
<circle cx="20" cy="20" r="20" fill="hsl(%d, 60%%, 15%%)"/>
<circle cx="20" cy="20" r="18" fill="none" stroke="hsl(%d, 60%%, 50%%)" stroke-width="1.5" opacity="0.5"/>
<text x="20" y="24" text-anchor="middle" fill="hsl(%d, 60%%, 50%%)" font-family="sans-serif" font-size="%d" font-weight="600">%s</text>
</svg>`, hue, hue, hue, fontSize, ticker)
}
