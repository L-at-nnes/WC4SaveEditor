package fileio

import (
	"os"
	"testing"
)

func TestReadUint8AtFileOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_uint8_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write test data: [0x00, 0x7F, 0x80, 0xFF]
	testData := []byte{0x00, 0x7F, 0x80, 0xFF}
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		expected int
	}{
		{
			name:     "Read first byte",
			offset:   0,
			expected: 0x00,
		},
		{
			name:     "Read second byte",
			offset:   1,
			expected: 0x7F,
		},
		{
			name:     "Read third byte",
			offset:   2,
			expected: 0x80,
		},
		{
			name:     "Read fourth byte",
			offset:   3,
			expected: 0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadUint8AtFileOffset(tempFile.Name(), tt.offset)
			if result != tt.expected {
				t.Errorf("ReadUint8AtFileOffset(%d) = 0x%02X, want 0x%02X", tt.offset, result, tt.expected)
			}
		})
	}
}

func TestReadUint16AtFileOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_uint16_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write test data: [0x00, 0x01, 0x7F, 0xFF, 0x80, 0x00]
	testData := []byte{0x00, 0x01, 0x7F, 0xFF, 0x80, 0x00}
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		expected int
	}{
		{
			name:     "Read first uint16 (little endian)",
			offset:   0,
			expected: 0x0100, // 0x00, 0x01 in little endian
		},
		{
			name:     "Read second uint16 (little endian)",
			offset:   2,
			expected: 0xFF7F, // 0x7F, 0xFF in little endian
		},
		{
			name:     "Read third uint16 (little endian)",
			offset:   4,
			expected: 0x0080, // 0x80, 0x00 in little endian
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadUint16AtFileOffset(tempFile.Name(), tt.offset)
			if result != tt.expected {
				t.Errorf("ReadUint16AtFileOffset(%d) = 0x%04X, want 0x%04X", tt.offset, result, tt.expected)
			}
		})
	}
}

func TestReadUint32AtFileOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_uint32_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write test data: [0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC]
	testData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC}
	if _, err := tempFile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		expected int
	}{
		{
			name:     "Read first uint32 (little endian)",
			offset:   0,
			expected: 0x03020100, // 0x00, 0x01, 0x02, 0x03 in little endian
		},
		{
			name:     "Read second uint32 (little endian)",
			offset:   4,
			expected: 0xFCFDFEFF, // 0xFF, 0xFE, 0xFD, 0xFC in little endian
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadUint32AtFileOffset(tempFile.Name(), tt.offset)
			if result != tt.expected {
				t.Errorf("ReadUint32AtFileOffset(%d) = 0x%08X, want 0x%08X", tt.offset, result, tt.expected)
			}
		})
	}
}

func TestWriteUint8AtFileOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_write_uint8_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Initialize file with zeros
	initialData := make([]byte, 10)
	if _, err := tempFile.Write(initialData); err != nil {
		t.Fatalf("Failed to write initial data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		value    int
		expected int
	}{
		{
			name:     "Write 0x7F at offset 0",
			offset:   0,
			value:    0x7F,
			expected: 0x7F,
		},
		{
			name:     "Write 0xFF at offset 5",
			offset:   5,
			value:    0xFF,
			expected: 0xFF,
		},
		{
			name:     "Write 0x00 at offset 9",
			offset:   9,
			value:    0x00,
			expected: 0x00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteUint8AtFileOffset(tempFile.Name(), tt.offset, tt.value)

			// Read back and verify
			result := ReadUint8AtFileOffset(tempFile.Name(), tt.offset)
			if result != tt.expected {
				t.Errorf("WriteUint8AtFileOffset(%d, %d) wrote 0x%02X, want 0x%02X",
					tt.offset, tt.value, result, tt.expected)
			}
		})
	}
}

func TestWriteUint16AtFileOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_write_uint16_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Initialize file with zeros
	initialData := make([]byte, 20)
	if _, err := tempFile.Write(initialData); err != nil {
		t.Fatalf("Failed to write initial data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		value    int
		expected int
	}{
		{
			name:     "Write 0x1234 at offset 0",
			offset:   0,
			value:    0x1234,
			expected: 0x1234,
		},
		{
			name:     "Write 0xFFFF at offset 4",
			offset:   4,
			value:    0xFFFF,
			expected: 0xFFFF,
		},
		{
			name:     "Write 0x0000 at offset 8",
			offset:   8,
			value:    0x0000,
			expected: 0x0000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteUint16AtFileOffset(tempFile.Name(), tt.offset, tt.value)

			// Read back and verify
			result := ReadUint16AtFileOffset(tempFile.Name(), tt.offset)
			if result != tt.expected {
				t.Errorf("WriteUint16AtFileOffset(%d, %d) wrote 0x%04X, want 0x%04X",
					tt.offset, tt.value, result, tt.expected)
			}
		})
	}
}

func TestWriteUint32AtFileOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_write_uint32_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Initialize file with zeros
	initialData := make([]byte, 20)
	if _, err := tempFile.Write(initialData); err != nil {
		t.Fatalf("Failed to write initial data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		value    int
		expected int
	}{
		{
			name:     "Write 0x12345678 at offset 0",
			offset:   0,
			value:    0x12345678,
			expected: 0x12345678,
		},
		{
			name:     "Write 0xFFFFFFFE at offset 4",
			offset:   4,
			value:    0xFFFFFFFE,
			expected: 0xFFFFFFFE,
		},
		{
			name:     "Write 0x00000000 at offset 8",
			offset:   8,
			value:    0x00000000,
			expected: 0x00000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteUint32AtFileOffset(tempFile.Name(), tt.offset, tt.value)

			// Read back and verify
			result := ReadUint32AtFileOffset(tempFile.Name(), tt.offset)
			if result != tt.expected {
				t.Errorf("WriteUint32AtFileOffset(%d, %d) wrote 0x%08X, want 0x%08X",
					tt.offset, tt.value, result, tt.expected)
			}
		})
	}
}

func TestWriteDataAtOffset(t *testing.T) {
	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "test_write_data_*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Initialize file with zeros
	initialData := make([]byte, 20)
	if _, err := tempFile.Write(initialData); err != nil {
		t.Fatalf("Failed to write initial data: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		offset   int
		data     []byte
		expected []byte
	}{
		{
			name:     "Write 4 bytes at offset 0",
			offset:   0,
			data:     []byte{0x01, 0x02, 0x03, 0x04},
			expected: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:     "Write 3 bytes at offset 5",
			offset:   5,
			data:     []byte{0xFF, 0xFE, 0xFD},
			expected: []byte{0xFF, 0xFE, 0xFD},
		},
		{
			name:     "Write empty data",
			offset:   10,
			data:     []byte{},
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteDataAtOffset(tempFile.Name(), tt.offset, tt.data)

			// Read back and verify
			for i, expectedByte := range tt.expected {
				result := ReadUint8AtFileOffset(tempFile.Name(), tt.offset+i)
				if result != int(expectedByte) {
					t.Errorf("WriteDataAtOffset(%d, %v) at position %d wrote 0x%02X, want 0x%02X",
						tt.offset, tt.data, i, result, expectedByte)
				}
			}
		})
	}
}
