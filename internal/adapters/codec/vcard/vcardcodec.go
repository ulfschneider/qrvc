package vcardcodec

import (
	"bytes"

	"github.com/emersion/go-vcard"
)

type Codec struct {
}

func NewCodec() Codec {
	return Codec{}
}

func (c *Codec) Encode(card vcard.Card) ([]byte, error) {
	var buf bytes.Buffer
	enc := vcard.NewEncoder(&buf)
	if err := enc.Encode(card); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (c *Codec) Decode(vcf []byte) (vcard.Card, error) {
	dec := vcard.NewDecoder(bytes.NewBuffer(vcf))
	card, err := dec.Decode()
	if err != nil {
		return nil, err
	}
	return card, nil
}
