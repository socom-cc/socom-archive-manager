package socomarchive

// ZdbHeaderV6 Based off of reverse engineer work done by Linblow (https://github.com/Linblow)
type ZdbHeaderV6 struct {
	HeaderSize     uint32
	EntrySize      uint32
	HeaderVersion  uint32
	DataOffset     uint32
	DataSize       uint32
	BuildTimestamp uint32
	Unk18          uint32
	BuildErrors    uint16
	BuildWarnings  uint16
	Unk20          [41]uint32
	EntryCount     uint32
}

func (h *ZdbHeaderV6) GetDataOffset() uint32 {
	return h.DataOffset
}

func (h *ZdbHeaderV6) GetEntryCount() uint32 {
	return h.EntryCount
}

func (h *ZdbHeaderV6) GetEntrySize() uint32 {
	return h.EntrySize
}
