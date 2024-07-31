package socomarchive

// ZdbHeaderV7 Based off of reverse engineer work done by Linblow (https://github.com/Linblow)
type ZdbHeaderV7 struct {
	HeaderSize     uint32
	EntrySize      uint32
	HeaderVersion  uint32
	VersionMinor   uint32
	DataOffset     uint32
	DataSize       uint32
	BuildTimestamp uint32
	Unk1c          uint32
	BuildErrors    uint16
	BuildWarnings  uint16
	Unk24          [40]uint32
	EntryCount     uint32
}

func (h *ZdbHeaderV7) GetDataOffset() uint32 {
	return h.DataOffset
}

func (h *ZdbHeaderV7) GetEntryCount() uint32 {
	return h.EntryCount
}

func (h *ZdbHeaderV7) GetEntrySize() uint32 {
	return h.EntrySize
}
