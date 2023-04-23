package types

type Bool struct {
	ScaleDecoder
}

func (b *Bool) Process() {
	b.Value = b.getNextBool()
}

func (b *Bool) Encode(value bool) string {
	if value {
		return "01"
	}
	return "00"
}

func (b *Bool) TypeStructString() string {
	return "Bool"
}
