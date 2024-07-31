package socomarchive

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type ZdbHeader interface {
	GetEntryCount() uint32
	GetEntrySize() uint32
	GetDataOffset() uint32
}

// ZdbEntryHeaderRaw Based off of reverse engineer work done by Linblow (https://github.com/Linblow)
type ZdbEntryHeaderRaw struct {
	Name       [64]byte
	EntryType  [4]byte
	Unk44      uint32
	DataOffset uint32
	DataSize   uint32
	DataHash   uint32
	Unk54      uint32
}

type ZdbEntryHeader struct {
	Name       string
	EntryType  string
	Unk44      uint32
	DataOffset uint32
	DataSize   uint32
	DataHash   uint32
	Unk54      uint32
}

type SocomArchive struct {
	data         []byte
	Header       ZdbHeader
	EntryHeaders []ZdbEntryHeader
}

func LoadSocomArchive(data []byte) (*SocomArchive, error) {
	archive := &SocomArchive{
		data: data,
	}

	err := archive.load()
	if err != nil {
		return nil, err
	}

	return archive, nil
}

func (a *SocomArchive) GetEntry(index int) ([]byte, error) {
	if len(a.EntryHeaders)-1 < index {
		return nil, errors.New("index out of range")
	}

	if a.EntryHeaders[index].DataSize <= 0 {
		return nil, errors.New("entry size out of range")
	}

	start := int(a.Header.GetDataOffset() + a.EntryHeaders[index].DataOffset)

	data := make([]byte, a.EntryHeaders[index].DataSize)
	copy(data, a.data[start:start+int(a.EntryHeaders[index].DataSize)])

	return data, nil
}

func (a *SocomArchive) load() error {
	headerVersion := binary.LittleEndian.Uint32(a.data[8:12])

	reader := bytes.NewReader(a.data)

	switch headerVersion {
	case 6:
		var header ZdbHeaderV6
		if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
			return err
		}
		a.Header = &header
	case 7:
		var header ZdbHeaderV7
		if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
			return err
		}
		a.Header = &header
	}

	for i := 0; i < int(a.Header.GetEntryCount()); i++ {
		var entryHeader ZdbEntryHeaderRaw
		if err := binary.Read(reader, binary.LittleEndian, &entryHeader); err != nil {
			return err
		}
		li := 0
		for idx, v := range entryHeader.Name {
			if v == 0 {
				li = idx
				break
			}
		}

		var name = make([]byte, li)
		copy(name, entryHeader.Name[:li])

		a.EntryHeaders = append(a.EntryHeaders, ZdbEntryHeader{
			Name:       string(name),
			EntryType:  string(entryHeader.EntryType[:]),
			Unk44:      entryHeader.Unk44,
			DataOffset: entryHeader.DataOffset,
			DataSize:   entryHeader.DataSize,
			DataHash:   entryHeader.DataHash,
			Unk54:      entryHeader.Unk54,
		})
	}

	return nil
}
