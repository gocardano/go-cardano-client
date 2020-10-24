package cbor

import (
	"github.com/gocardano/go-cardano-client/errors"
)

const (
	defaultUint8ValueError  = uint8(0)
	defaultUint16ValueError = uint16(0)
	defaultUint32ValueError = uint32(0)
	defaultUint64ValueError = uint64(0)
)

var bitMaskPositions = map[uint8]uint8{
	0: uint8(0x80),
	1: uint8(0x40),
	2: uint8(0x20),
	3: uint8(0x10),
	4: uint8(0x08),
	5: uint8(0x04),
	6: uint8(0x02),
	7: uint8(0x01),
}

// BitstreamReader provides ability to read bits/byte(s) from a byte array
type BitstreamReader struct {
	data               []byte
	currentBitPosition uint64
	lengthInBytes      uint32
	lengthInBits       uint64
}

// NewBitstreamReader returns instance of bitstream reader
func NewBitstreamReader(data []byte) *BitstreamReader {
	return &BitstreamReader{
		data:               data,
		currentBitPosition: 0,
		lengthInBytes:      uint32(len(data)),
		lengthInBits:       uint64(len(data) * 8),
	}
}

// HasMoreBits returns flag to indicate if there's more bits to parse
func (r *BitstreamReader) HasMoreBits() bool {
	return r.lengthInBits > r.currentBitPosition
}

// ReadBitsAsUint64 returns the next x bits as a uint64
func (r *BitstreamReader) ReadBitsAsUint64(bitCount uint8) (uint64, error) {
	result, err := r.doReadBitsAsUint64(0, r.currentBitPosition, bitCount)
	r.currentBitPosition += uint64(bitCount)
	return result, err
}

// ReadBitsAsUint32 returns the next x bits as a uint32
func (r *BitstreamReader) ReadBitsAsUint32(bitCount uint8) (uint32, error) {
	result, err := r.doReadBitsAsUint32(0, r.currentBitPosition, bitCount)
	r.currentBitPosition += uint64(bitCount)
	return result, err
}

// ReadBitsAsUint16 returns the next x bits as a uint16
func (r *BitstreamReader) ReadBitsAsUint16(bitCount uint8) (uint16, error) {
	result, err := r.doReadBitsAsUint16(0, r.currentBitPosition, bitCount)
	r.currentBitPosition += uint64(bitCount)
	return result, err
}

// ReadBitsAsUint8 returns the next x bits as a uint8
func (r *BitstreamReader) ReadBitsAsUint8(bitCount uint8) (uint8, error) {
	result, err := r.doReadBitsAsUint8(0, r.currentBitPosition, bitCount)
	r.currentBitPosition += uint64(bitCount)
	return result, err
}

// ReadBytes returns the next x bytes as byte array
func (r *BitstreamReader) ReadBytes(byteCount uint64) ([]byte, error) {

	if r.currentBitPosition%8 != 0 {
		return nil, errors.NewError(errors.ErrCborUnhandledReadBytesInTermsOfBits)
	}

	if r.currentBitPosition+(byteCount*8) > r.lengthInBits {
		return nil, errors.NewError(errors.ErrBitstreamReaderEOF)
	}

	startPosition := r.currentBitPosition / 8
	endPosition := startPosition + byteCount

	r.currentBitPosition += byteCount * 8

	return r.data[startPosition:endPosition], nil
}

// doReadBitsAsUint64 returns the bits values as a uint
func (r *BitstreamReader) doReadBitsAsUint64(bytePosition uint32, bitPosition uint64, bitCount uint8) (uint64, error) {
	if bitCount > 64 {
		return defaultUint64ValueError, errors.NewError(errors.ErrBitstreamVarInsufficientCapacity)
	}
	return r.doReadBits(bytePosition, bitPosition, bitCount)
}

// doReadBitsAsUint32 returns the bits values as a uint
func (r *BitstreamReader) doReadBitsAsUint32(bytePosition uint32, bitPosition uint64, bitCount uint8) (uint32, error) {
	if bitCount > 32 {
		return defaultUint32ValueError, errors.NewError(errors.ErrBitstreamVarInsufficientCapacity)
	}
	val, err := r.doReadBitsAsUint64(bytePosition, bitPosition, bitCount)
	return uint32(val), err
}

// doReadBitsAsUint16 returns the bits values as a uint
func (r *BitstreamReader) doReadBitsAsUint16(bytePosition uint32, bitPosition uint64, bitCount uint8) (uint16, error) {
	if bitCount > 16 {
		return defaultUint16ValueError, errors.NewError(errors.ErrBitstreamVarInsufficientCapacity)
	}
	val, err := r.doReadBitsAsUint64(bytePosition, bitPosition, bitCount)
	return uint16(val), err
}

// doReadBitsAsUint8 returns the bits values as a uint
func (r *BitstreamReader) doReadBitsAsUint8(bytePosition uint32, bitPosition uint64, bitCount uint8) (uint8, error) {
	if bitCount > 8 {
		return defaultUint8ValueError, errors.NewError(errors.ErrBitstreamVarInsufficientCapacity)
	}
	val, err := r.doReadBitsAsUint64(bytePosition, bitPosition, bitCount)
	return uint8(val), err
}

// doReadBitsAsUint64 return the bits value as a uint64
func (r *BitstreamReader) doReadBits(bytePosition uint32, bitPosition uint64, bitCount uint8) (uint64, error) {

	accessingMaxBits := uint64(bytePosition*8) + bitPosition + uint64(bitCount)

	if r.lengthInBits < accessingMaxBits {
		return defaultUint64ValueError, errors.NewMessageErrorf(errors.ErrBitstreamReaderEOF,
			"BytePosition: [%d], BitPosition: [%d], BitCount: [%d], TotalBits: [%d], AccessingMaxBit: [%d]",
			bytePosition, bitPosition, bitCount, r.lengthInBits, accessingMaxBits)
	}
	result := uint64(0)
	for i := uint8(0); i < bitCount; i++ {
		bit, err := r.doReadBit(bytePosition, bitPosition+uint64(i))
		if err != nil {
			return defaultUint64ValueError, err
		}
		result = result | uint64(bit)
		if i < bitCount-1 {
			result = result << 1
		}
	}
	return result, nil
}

// doReadBit returns the bit value at position bytePosition and bitPosition
func (r *BitstreamReader) doReadBit(bytePosition uint32, bitPosition uint64) (uint8, error) {
	if r.lengthInBits <= uint64(bytePosition*8)+bitPosition {
		return defaultUint8ValueError, errors.NewError(errors.ErrBitstreamReaderEOF)
	}
	if bitPosition > 7 {
		bytePosition += uint32(bitPosition / 8)
		bitPosition = bitPosition % 8
	}
	return (r.data[bytePosition] & bitMaskPositions[uint8(bitPosition)]) >> (7 - bitPosition), nil
}
