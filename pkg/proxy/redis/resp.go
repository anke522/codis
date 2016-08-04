// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package redis

import "fmt"

type RespType byte

const (
	TypeString    RespType = '+'
	TypeError     RespType = '-'
	TypeInt       RespType = ':'
	TypeBulkBytes RespType = '$'
	TypeArray     RespType = '*'
)

func (t RespType) String() string {
	switch t {
	case TypeString:
		return "<string>"
	case TypeError:
		return "<error>"
	case TypeInt:
		return "<int>"
	case TypeBulkBytes:
		return "<bulkbytes>"
	case TypeArray:
		return "<array>"
	default:
		return fmt.Sprintf("<unknown-0x%02x>", byte(t))
	}
}

type Resp struct {
	Type RespType

	Value []byte
	Array []*Resp
}

func (r *Resp) IsString() bool {
	return r.Type == TypeString
}

func (r *Resp) IsError() bool {
	return r.Type == TypeError
}

func (r *Resp) IsInt() bool {
	return r.Type == TypeInt
}

func (r *Resp) IsBulkBytes() bool {
	return r.Type == TypeBulkBytes
}

func (r *Resp) IsArray() bool {
	return r.Type == TypeArray
}

func NewString(value []byte) *Resp {
	return &Resp{
		Type:  TypeString,
		Value: value,
	}
}

func NewError(value []byte) *Resp {
	return &Resp{
		Type:  TypeError,
		Value: value,
	}
}

func NewInt(value []byte) *Resp {
	return &Resp{
		Type:  TypeInt,
		Value: value,
	}
}

func NewBulkBytes(value []byte) *Resp {
	return &Resp{
		Type:  TypeBulkBytes,
		Value: value,
	}
}

func NewArray(array []*Resp) *Resp {
	return &Resp{
		Type:  TypeArray,
		Array: array,
	}
}

type RespAlloc struct {
	alloc struct {
		buf []Resp
		off int
	}
	slice struct {
		buf []*Resp
		off int
	}
}

func (p *RespAlloc) New() *Resp {
	var d = &p.alloc
	if len(d.buf) == d.off {
		d.buf = make([]Resp, 16)
		d.off = 0
	}
	r := &d.buf[d.off]
	d.off += 1
	return r
}

func (p *RespAlloc) MakeSlice(n int) []*Resp {
	if n >= 32 {
		return make([]*Resp, n)
	}
	var d = &p.slice
	if max := len(d.buf) - d.off; max < n {
		d.buf = make([]*Resp, 512)
		d.off = 0
	}
	n += d.off
	s := d.buf[d.off:n:n]
	d.off = n
	return s
}
