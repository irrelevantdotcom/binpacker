package binpacker

import (
	"bytes"
	"testing"
	"math"

	"github.com/stretchr/testify/assert"
	"encoding/binary"
)

func TestPushByte(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushByte(0x01)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{1}, "byte error.")
}

func TestPushBytes(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushBytes([]byte{0x01, 0x002})
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0x01, 0x02}, "bytes error.")
}

func TestPushUint8(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushUint8(1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{1}, "uint8 error.")
}

func TestPushUint16(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushUint16(1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 1}, "uint16 error.")
}

func TestPushInt16(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushInt16(-1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{255, 255}, "uint16 error.") // -1 eq 255 255
}

func TestPushUint32(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushUint32(1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 0, 0, 1}, "uint32 error.")
}

func TestPushInt32(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushInt32(-1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{255, 255, 255, 255}, "int32 error.")
}

func TestPushUint64(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushUint64(1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 0, 0, 0, 0, 0, 0, 1}, "uint64 error.")
}

func TestPushInt64(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushInt64(-1)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{255, 255, 255, 255, 255, 255, 255, 255}, "int64 error.")
}

func TestPushFloat32(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushFloat32(math.SmallestNonzeroFloat32)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 0, 0, 1}, "float32 error.")
}

func TestPushFloat64(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushFloat64(math.SmallestNonzeroFloat64)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 0, 0, 0, 0, 0, 0, 1}, "float64 error.")
}

func TestPushString(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushString("Hi")
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{'H', 'i'}, "string error.")
}

func TestCombinedPush(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	p.PushUint16(1).PushString("Hi")
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 1, 'H', 'i'}, "combine push error.")
}


func TestPack(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	b = p.Pack("nNa5xa1", uint16(1), uint32(2), "Hello", "!")
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, b.Bytes(), []byte{0, 1, 0, 0, 0, 2, 'H', 'e', 'l', 'l', 'o', 0, '!' }, "Pack function error.")
}

func TestUnpack(t *testing.T) {
	b := new(bytes.Buffer)
	p := NewPacker(binary.BigEndian, b)
	
	d := bytes.NewBuffer([]byte{0, 1, 0, 0, 0, 2, 'H', 'e', 'l', 'l', 'o'})
	up := "nfred/Njim/a5shiela"
	r := p.Unpack(up, d)
	assert.Equal(t, p.Error(), nil, "Has error.")
	assert.Equal(t, r, map[string]interface{}{"fred": uint16(1), "jim": uint32(2), "shiela": "Hello"})
}
