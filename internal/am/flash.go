package am

const (
	FlashStoreKey = "flash"
)

type (
	FlashSet []FlashItem

	FlashItem struct {
		Msg  string
		Type MsgType
	}

	// MsgType stands for message type
	MsgType string
)

func NewFlashSet() FlashSet {
	return make(FlashSet, 0)
}

func (f FlashSet) IsEmpty() bool {
	return len(f) == 0
}

func (f FlashSet) AddItem(fi FlashItem) FlashSet {
	return append(f, fi)
}

func (f FlashSet) AddItems(fis []FlashItem) FlashSet {
	return append(f, fis...)
}

func (fi FlashItem) IsEmpty() bool {
	return fi == FlashItem{}
}

func NewFlashItem(msg string, msgType MsgType) FlashItem {
	return FlashItem{
		Msg:  msg,
		Type: msgType,
	}
}
