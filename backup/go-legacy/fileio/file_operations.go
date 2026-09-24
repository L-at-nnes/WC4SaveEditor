package fileio

import (
	"encoding/binary"
	"io"
	"log"
	"math"
	"os"
)

// ReadUint8AtFileOffset reads a uint8 value from a file at the specified offset
func ReadUint8AtFileOffset(inputFilename string, offset int) int {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}
	defer inputFile.Close()

	byteData := make([]byte, 1)
	if _, err := inputFile.ReadAt(byteData, int64(offset)); err != nil {
		log.Fatal("Failed to read uint8 from file:", err)
	}

	return int(byteData[0])
}

// ReadUint16AtFileOffset reads a uint16 value from a file at the specified offset
func ReadUint16AtFileOffset(inputFilename string, offset int) int {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}
	defer inputFile.Close()

	byteData := make([]byte, 2)
	if _, err := inputFile.ReadAt(byteData, int64(offset)); err != nil {
		log.Fatal("Failed to read uint16 from file:", err)
	}

	return int(binary.LittleEndian.Uint16(byteData))
}

func ReadUint32AtFileOffset(inputFilename string, offset int) int {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}
	defer inputFile.Close()

	byteData := make([]byte, 4)
	if _, err := inputFile.ReadAt(byteData, int64(offset)); err != nil {
		log.Fatal("Failed to read uint32 from file:", err)
	}

	return int(binary.LittleEndian.Uint32(byteData))
}

// WriteUint8AtFileOffset writes a uint8 value at the specified file offset
func WriteUint8AtFileOffset(inputFilename string, offset int, value int) {
	if value > math.MaxUint8 {
		log.Fatal("Value is too large for uint8")
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}
	defer inputFile.Close()

	if _, err := inputFile.WriteAt([]byte{uint8(value)}, int64(offset)); err != nil {
		log.Fatal("Failed to write uint8 to file:", err)
	}
}

// WriteUint16AtFileOffset writes a uint16 value at the specified file offset
func WriteUint16AtFileOffset(inputFilename string, offset int, updatedValue int) {
	if updatedValue > math.MaxUint16 {
		log.Fatal("Value is too large for uint16")
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}
	defer inputFile.Close()

	byteArrUnitType := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteArrUnitType, uint16(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

// WriteUint32AtFileOffset writes a uint32 value at the specified file offset
func WriteUint32AtFileOffset(inputFilename string, offset int, updatedValue int) {
	if updatedValue > math.MaxUint32 {
		log.Fatal("Value is too large for uint32")
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}
	defer inputFile.Close()

	byteArrUnitType := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteArrUnitType, uint32(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

// WriteDataAtOffset writes data to a file at the specified offset (simple overwrite)
func WriteDataAtOffset(inputFilename string, offset int, newData []byte) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}
	defer inputFile.Close()

	// Simple overwrite - no shifting needed for fixed-size data
	if _, err := inputFile.WriteAt(newData, int64(offset)); err != nil {
		log.Fatal("Failed to write data:", err)
	}
}

// EnsureBackup creates a one-time ".bak" copy of the save file before it is
// modified for the first time, so a corrupted write can be recovered from.
// If a backup already exists, it is left untouched (it should hold the
// original, unmodified save, not an already-edited version).
func EnsureBackup(inputFilename string) error {
	backupFilename := inputFilename + ".bak"

	if _, err := os.Stat(backupFilename); err == nil {
		return nil
	}

	src, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(backupFilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
