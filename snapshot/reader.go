package snapshot

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var magicNumber = []byte{0x50, 0x05, 0x51, 0x30}

type Reader struct {
	Header         *Header
	CurrentSection *Section

	filename   string
	fl         *os.File
	nextOffset uint64

	handlers map[SectionName]sectionHandlerFunc
}

func NewDefaultReader(filename string) (r *Reader, err error) {
	reader, err := NewReader(filename)
	if err != nil {
		return nil, err
	}
	reader.RegisterSectionHandler(SectionNameChainSnapshotHeader, readChainSnapshotHeader)
	reader.RegisterSectionHandler(SectionNameBlockState, readBlockState)
	reader.RegisterSectionHandler(SectionNameAccountObject, readAccountObjects)
	reader.RegisterSectionHandler(SectionNameAccountMetadataObject, readAccountMetadataObjects)
	reader.RegisterSectionHandler(SectionNameAccountRamCorrectionObject, readAccountRAMCorrectionObject)
	reader.RegisterSectionHandler(SectionNameGlobalPropertyObject, readGlobalPropertyObject)
	reader.RegisterSectionHandler(SectionNameProtocolStateObject, readProtocolStateObject)
	reader.RegisterSectionHandler(SectionNameDynamicGlobalPropertyObject, readDynamicGlobalPropertyObject)
	reader.RegisterSectionHandler(SectionNameBlockSummaryObject, readBlockSummary)
	reader.RegisterSectionHandler(SectionNameTransactionObject, readTransactionObject)
	reader.RegisterSectionHandler(SectionNameGeneratedTransactionObject, readGeneratedTransactionObject)
	reader.RegisterSectionHandler(SectionNameCodeObject, readCodeObject)
	reader.RegisterSectionHandler(SectionNameContractTables, readContractTables)
	reader.RegisterSectionHandler(SectionNamePermissionObject, readPermissionObject)
	reader.RegisterSectionHandler(SectionNamePermissionLinkObject, readPermissionLinkObject)
	reader.RegisterSectionHandler(SectionNameResourceLimitsObject, readResourceLimitsObject)
	reader.RegisterSectionHandler(SectionNameResourceUsageObject, readResourceUsageObject)
	reader.RegisterSectionHandler(SectionNameResourceLimitsStateObject, readResourceLimitsStateObject)
	reader.RegisterSectionHandler(SectionNameResourceLimitsConfigObject, readResourceLimitsConfigObject)
	reader.RegisterSectionHandler(SectionNameGenesisState, readGenesisState)

	// Ultra Specific
	reader.RegisterSectionHandler(SectionAccountFreeActionsObject, readAccountFreeActionsObject)

	return reader, nil
}

func NewReader(filename string) (r *Reader, err error) {
	r = &Reader{
		filename: filename,
		handlers: map[SectionName]sectionHandlerFunc{},
	}
	r.fl, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	h, err := r.readHeader()
	if err != nil {
		return nil, err
	}

	beginOffset, err := r.fl.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	r.Header = h
	r.nextOffset = uint64(beginOffset)

	return
}

func (r *Reader) RegisterSectionHandler(s SectionName, h sectionHandlerFunc) {
	r.handlers[s] = h
}

func (r *Reader) readHeader() (*Header, error) {
	buf := make([]byte, 8)
	if _, err := r.fl.Read(buf); err != nil {
		return nil, err
	}

	if !bytes.Equal(buf[:4], magicNumber) {
		return nil, fmt.Errorf("invalid magic number (first 4 bytes): %v, expected %v", buf[:4], magicNumber)
	}

	h := &Header{
		Version: binary.LittleEndian.Uint32(buf[4:8]),
	}

	return h, nil
}

var ErrSectionHandlerNotFound = errors.New("section handler not found")

// Deprecated: Use ErrSectionHandlerNotFound instead
var SectionHandlerNotFound = ErrSectionHandlerNotFound

func (r *Reader) HasSectionHandler(s *Section) bool {
	_, found := r.handlers[r.CurrentSection.Name]
	return found
}

func (r *Reader) ProcessCurrentSection(f sectionCallbackFunc) error {
	h, found := r.handlers[r.CurrentSection.Name]
	if !found {
		return ErrSectionHandlerNotFound
	}
	return h(r.CurrentSection, f)
}

// Next retrieves the next section.
func (r *Reader) NextSection() error {
	beginOffset, err := r.fl.Seek(int64(r.nextOffset), io.SeekStart)
	if err != nil {
		return err
	}

	vals := make([]byte, 16)
	bytesRead, err := r.fl.Read(vals)
	if err != nil {
		return err
	}

	// end marker
	if bytesRead == 8 && bytes.Equal(vals[:8], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}) {
		return io.EOF
	}

	sectionSize := binary.LittleEndian.Uint64(vals[:8])
	rowCount := binary.LittleEndian.Uint64(vals[8:16])

	str, err := readZeroTerminatedString(r.fl)
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("EOF while reading string section (partial: %s)", str)
		}
		return err
	}

	r.nextOffset = uint64(beginOffset) + sectionSize + 8 // well well, sectionSize includes the rowCount I guess?

	r.CurrentSection = &Section{
		Name:       SectionName(strings.TrimRight(str, string([]byte{0x00}))),
		Offset:     uint64(beginOffset),
		Size:       sectionSize,
		RowCount:   rowCount,
		BufferSize: sectionSize - uint64(len(str)) - 1 /* str-pad 0x00 byte */ - 8,
		Buffer:     r.fl,
	}
	return nil
}

func (r *Reader) Close() error {
	return r.fl.Close()
}
