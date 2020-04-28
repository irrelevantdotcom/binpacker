package binpacker

import (
	"encoding/binary"
	"bytes"
	"io"
	"math"
	"strings"
	"strconv"
	"unicode"
	"unicode/utf8"
)

// Packer is a binary packer helps you pack data into an io.Writer.
type Packer struct {
	writer io.Writer
	endian binary.ByteOrder
	err    error
}

func (p *Packer) Pack(format string, args ...interface{}) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	packer := NewPacker(binary.BigEndian, buffer)
	f := explodePack(format)
	i := 0
	j := 0
//	for _, v := range args {
	for (i < len(f) && j < len(args)) {
		v := args[j]
		packType := f[i]
		switch (packType[0]) {
			case 'n':
				packer.endian = binary.BigEndian
				packer.PushUint16(v.(uint16))
				j++
			case 'N':
				packer.endian = binary.BigEndian
				packer.PushUint32(v.(uint32))
				j++
			case 'x':								// x doesn't consume a value.
				packer.PushByte(0)
			case 'a':																	// "a256" or a* (special case)
				if packType[1:] == "*" {
					packer.PushString(strings.TrimRight(v.(string),"\000"))				// remove trailing NULLs. If needed, add an "x"
				} else {
					w, _ := strconv.Atoi(packType[1:])
					packer.PushString(fixedLengthString(w,v.(string)))	// fixed width.
				}
				j++
		}
		i++
	}
	return buffer  //, nil
}

/*
	format would be something like nfred/Njim/a5shiela
*/
func (p *Packer) Unpack(format string, data *bytes.Buffer) map[string]interface{} {

	f := explodeUnpack(format)
	res := make(map[string]interface{});
	s := 0;
	
	for k,v := range f {
		chr, _ := utf8.DecodeRuneInString(v)
		switch (chr) {
			case 'n':
				res[k] = binary.BigEndian.Uint16(data.Bytes()[s:])
				s += 2
			case 'N':	
				res[k] = binary.BigEndian.Uint32(data.Bytes()[s:])
				s += 4
			case 'a':
				l := 0
				if v[1:] == "*x" {
					l = clen([]byte(v))
				} else {
					l, _ = strconv.Atoi(v[1:])
				}
				t := s + l
				res[k] = string(data.Bytes()[s:t])
				s = t
		}
	}
	return res

}


func fixedLengthString(length int, str string) string {
	if len(str) < length {
		str = str + strings.Repeat("\000",length - len(str))
	}
	return str[0:length]
}

//https://stackoverflow.com/a/27834860
func clen(n []byte) int {
    for i := 0; i < len(n); i++ {
        if n[i] == 0 {
            return i
        }
    }
    return len(n)
}



// explode splits s into a slice of UTF-8 strings,
// one string per Unicode character + numerics 
// Invalid UTF-8 sequences become correct encodings of U+FFFD.
// based on https://golang.org/src/strings/strings.go?m=text

func explodePack(s string) []string {
	n := utf8.RuneCountInString(s)
	a := make([]string, n)
	for i := 0; i < n-1; i++ {
		ch := ""
		for  {
			chr, size := utf8.DecodeRuneInString(s)
			ch = ch + string(chr)
			s = s[size:]
			chr, size = utf8.DecodeRuneInString(s)
			if !unicode.IsNumber(chr) {
				break
			}
		}
		a[i] = ch

//		if chr == utf8.RuneError {
//			a[i] = string(utf8.RuneError)
//		}
		if s == "" {
			break
		}
	}
	if n > 0 {
		a[n-1] = s
	}
	return a
}

// a256acctName/a16oneTimePasswd/NresVersion/nlocale/narch/nversion to
// "acctName"=> "a256"
func explodeUnpack(s string) map[string]string {
	res:= make(map[string]string)
	fields := strings.Split(s, "/")
	
	for i:=0; i<len(fields); i++ {
		n := fields[i]
		chr, size := utf8.DecodeRuneInString(n)
		if (chr == 'n' || chr == 'N') {
			fn := fields[i][size:]
			res[fn] = string(chr)
		} else if (chr == 'a') {							// a22name or a*xname. aname should be treated as a1name
			ch := ""
			typ := string(chr)
			n = n[size:]
			chr, size := utf8.DecodeRuneInString(n)
			if unicode.IsNumber(chr) || chr == '*' {
				// add characters to ch until there are no numbers ..
				for  {
					ch = ch + string(chr)
					n = n[size:]
					chr, size = utf8.DecodeRuneInString(n)
					if !unicode.IsNumber(chr) {
						break
					}
				}
				if chr == 'x' {							// special case to grab trailing x after a number or *
					ch = ch + string(chr)
					n = n[size:]
				}
			}
			
			res[n] = typ + ch
		} else {
			// unknown type letter
			
		}
	}
	return res
}


// NewPacker returns a *Packer hold an io.Writer. User must provide the byte order explicitly.
func NewPacker(endian binary.ByteOrder, writer io.Writer) *Packer {
	return &Packer{
		writer: writer,
		endian: endian,
	}
}

// Error returns an error if any errors exists
func (p *Packer) Error() error {
	return p.err
}

// PushByte write a single byte into writer.
func (p *Packer) PushByte(b byte) *Packer {
	return p.errFilter(func() {
		_, p.err = p.writer.Write([]byte{b})
	})
}

// PushBytes write a bytes array into writer.
func (p *Packer) PushBytes(bytes []byte) *Packer {
	return p.errFilter(func() {
		_, p.err = p.writer.Write(bytes)
	})
}

// PushUint8 write a uint8 into writer.
func (p *Packer) PushUint8(i uint8) *Packer {
	return p.errFilter(func() {
		_, p.err = p.writer.Write([]byte{byte(i)})
	})
}

// PushUint16 write a uint16 into writer.
func (p *Packer) PushUint16(i uint16) *Packer {
	return p.errFilter(func() {
		buffer := make([]byte, 2)
		p.endian.PutUint16(buffer, i)
		_, p.err = p.writer.Write(buffer)
	})
}

// PushUint16 write a int16 into writer.
func (p *Packer) PushInt16(i int16) *Packer {
	return p.PushUint16(uint16(i))
}

// PushUint32 write a uint32 into writer.
func (p *Packer) PushUint32(i uint32) *Packer {
	return p.errFilter(func() {
		buffer := make([]byte, 4)
		p.endian.PutUint32(buffer, i)
		_, p.err = p.writer.Write(buffer)
	})
}

// PushInt32 write a int32 into writer.
func (p *Packer) PushInt32(i int32) *Packer {
	return p.PushUint32(uint32(i))
}

// PushUint64 write a uint64 into writer.
func (p *Packer) PushUint64(i uint64) *Packer {
	return p.errFilter(func() {
		buffer := make([]byte, 8)
		p.endian.PutUint64(buffer, i)
		_, p.err = p.writer.Write(buffer)
	})
}

// PushInt64 write a int64 into writer.
func (p *Packer) PushInt64(i int64) *Packer {
	return p.PushUint64(uint64(i))
}

// PushFloat32 write a float32 into writer.
func (p *Packer) PushFloat32(i float32) *Packer {
	return p.PushUint32(math.Float32bits(i))
}

// PushFloat64 write a float64 into writer.
func (p *Packer) PushFloat64(i float64) *Packer {
	return p.PushUint64(math.Float64bits(i))
}

// PushString write a string into writer.
func (p *Packer) PushString(s string) *Packer {
	return p.errFilter(func() {
		_, p.err = p.writer.Write([]byte(s))
	})
}

func (p *Packer) errFilter(f func()) *Packer {
	if p.err == nil {
		f()
	}
	return p
}
