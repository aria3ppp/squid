package squid

import (
	"sort"
	"testing"
	"time"
)

func TestSQUIDGeneration(t *testing.T) {
	generator, err := NewGenerator()
	if err != nil {
		t.Fatalf("Failed to create SQUID generator: %v", err)
	}

	id := generator.New()
	if len(id) != 45 {
		t.Errorf("Expected SQUID length to be 45, got %d", len(id))
	}
}

func TestSQUIDParsing(t *testing.T) {
	generator, err := NewGenerator()
	if err != nil {
		t.Fatalf("Failed to create SQUID generator: %v", err)
	}

	id := generator.New()
	timestamp, random, counter, err := generator.Parse(id)

	if err != nil {
		t.Errorf("Failed to parse SQUID: %v", err)
	}

	if timestamp.After(time.Now()) {
		t.Error("Parsed timestamp is in the future")
	}

	if len(random) != randomLength {
		t.Errorf("Expected random bytes length to be %d, got %d", randomLength, len(random))
	}

	if counter == 0 {
		t.Error("Counter should not be zero")
	}
}

func TestSQUIDSortability(t *testing.T) {
	generator, err := NewGenerator()
	if err != nil {
		t.Fatalf("Failed to create SQUID generator: %v", err)
	}

	var ids []string
	for i := 0; i < 1000; i++ {
		ids = append(ids, generator.New())
		time.Sleep(time.Millisecond) // Ensure some time passes between generations
	}

	sorted := make([]string, len(ids))
	copy(sorted, ids)
	sort.Strings(sorted)

	for i := range ids {
		if ids[i] != sorted[i] {
			t.Error("SQUIDs are not naturally sorted")
			break
		}
	}
}

func TestSQUIDUniqueness(t *testing.T) {
	generator, err := NewGenerator()
	if err != nil {
		t.Fatalf("Failed to create SQUID generator: %v", err)
	}

	idMap := make(map[string]bool)
	for i := 0; i < 100000; i++ {
		id := generator.New()
		if idMap[id] {
			t.Errorf("Duplicate SQUID generated: %s", id)
		}
		idMap[id] = true
	}
}

func TestSQUIDTimestampMonotonicity(t *testing.T) {
	generator, err := NewGenerator()
	if err != nil {
		t.Fatalf("Failed to create SQUID generator: %v", err)
	}

	var lastTimestamp time.Time
	for i := 0; i < 1000; i++ {
		id := generator.New()
		timestamp, _, _, err := generator.Parse(id)
		if err != nil {
			t.Errorf("Failed to parse SQUID: %v", err)
		}

		if !lastTimestamp.IsZero() && timestamp.Before(lastTimestamp) {
			t.Errorf("Timestamp monotonicity violated: %v is before %v", timestamp, lastTimestamp)
		}
		lastTimestamp = timestamp
	}
}

func TestMultipleGenerators(t *testing.T) {
	gen1, _ := NewGenerator()
	gen2, _ := NewGenerator()

	id1 := gen1.New()
	id2 := gen2.New()

	if id1 == id2 {
		t.Error("Different generators produced the same SQUID")
	}
}

func BenchmarkSQUIDGeneration(b *testing.B) {
	generator, _ := NewGenerator()
	for i := 0; i < b.N; i++ {
		generator.New()
	}
}
