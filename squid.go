package squid

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

const (
	timestampLength = 10
	randomLength    = 16
	counterLength   = 2
	squidLength     = timestampLength + randomLength + counterLength
)

// Generator represents a SQUID generator
type Generator struct {
	machineID     []byte
	processID     int
	encoder       *base32.Encoding
	counter       uint32
	lastTimestamp int64
}

// NewGenerator creates a new SQUID generator
func NewGenerator() (*Generator, error) {
	machineID := make([]byte, 6)
	_, err := rand.Read(machineID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate machine ID: %v", err)
	}

	return &Generator{
		machineID: machineID,
		processID: os.Getpid(),
		encoder:   base32.NewEncoding("0123456789abcdefghijklmnopqrstuv").WithPadding(base32.NoPadding),
	}, nil
}

// New generates a new SQUID
func (g *Generator) New() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Ensure the timestamp is always increasing
	for {
		lastTS := atomic.LoadInt64(&g.lastTimestamp)
		if timestamp <= lastTS {
			timestamp = lastTS + 1
		}
		if atomic.CompareAndSwapInt64(&g.lastTimestamp, lastTS, timestamp) {
			break
		}
	}

	count := atomic.AddUint32(&g.counter, 1)

	randomBytes := make([]byte, randomLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	squid := make([]byte, squidLength)

	// Timestamp (10 bytes)
	for i := 0; i < timestampLength; i++ {
		squid[i] = byte(timestamp >> (8 * (timestampLength - 1 - i)))
	}

	// Random (16 bytes)
	copy(squid[timestampLength:], randomBytes)

	// Counter (2 bytes)
	squid[squidLength-2] = byte(count >> 8)
	squid[squidLength-1] = byte(count)

	fmt.Printf("squid byte length: %d\n", len(squid))
	encoded := g.encoder.EncodeToString(squid)
	fmt.Printf("encoded squid length: %d\n", len(encoded))

	return encoded
}

// Parse decodes a SQUID string into its components
func (g *Generator) Parse(s string) (time.Time, []byte, uint16, error) {
	decoded, err := g.encoder.DecodeString(s)
	if err != nil {
		return time.Time{}, nil, 0, err
	}

	if len(decoded) != squidLength {
		return time.Time{}, nil, 0, fmt.Errorf("invalid SQUID length")
	}

	timestamp := int64(0)
	for i := 0; i < timestampLength; i++ {
		timestamp = (timestamp << 8) | int64(decoded[i])
	}

	random := decoded[timestampLength : timestampLength+randomLength]
	counter := uint16(decoded[squidLength-2])<<8 | uint16(decoded[squidLength-1])

	return time.Unix(0, timestamp*int64(time.Millisecond)), random, counter, nil
}
