package json_writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

type T struct {
	Writer   io.Writer
	First    bool
	IsObj    []bool
	WroteKey bool
}

type marker struct{ _ int }

var Obj = &marker{}
var Arr = &marker{}
var End = &marker{}

func Wrap(writer io.Writer) *T {
	return &T{
		Writer: writer,
		First:  true,
		IsObj:  make([]bool, 0, 8),
	}
}

func New() *T {
	return &T{
		Writer: &strings.Builder{},
		First:  true,
		IsObj:  make([]bool, 0, 8),
	}
}

func (j *T) Write(value any) error {
	if value == End {
		return j.end()
	}
	if len(j.IsObj) == 0 && !j.First {
		panic("consecutive top-level values")
	}
	lastIndex := len(j.IsObj) - 1
	isKey := lastIndex >= 0 && j.IsObj[lastIndex] && !j.WroteKey
	if isKey {
		if _, ok := value.(string); !ok {
			panic("expected a key")
		}
	}
	switch value {
	case Obj:
		return j.start('}')
	case Arr:
		return j.start('[')
	}
	if !j.First {
		if _, err := j.Writer.Write([]byte{','}); err != nil {
			return err
		}
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	} else if _, err = j.Writer.Write(b); err != nil {
		return err
	}
	if isKey {
		if _, err := j.Writer.Write([]byte{':'}); err != nil {
			return err
		}
		j.WroteKey = true
		j.First = true
	} else {
		j.WroteKey = false
		j.First = false
	}
	return nil
}

func (j *T) start(b byte) error {
	if !j.First {
		if _, err := j.Writer.Write([]byte{','}); err != nil {
			return err
		}
	}
	if _, err := j.Writer.Write([]byte{b}); err != nil {
		return err
	}
	j.IsObj = append(j.IsObj, b == '{')
	j.First = true
	j.WroteKey = false
	return nil
}

func (j *T) end() error {
	lastIndex := len(j.IsObj) - 1
	if lastIndex < 0 || !j.IsObj[lastIndex] {
		panic("not in an object or array")
	} else if j.WroteKey {
		panic("expected a value after the key")
	}
	isObj := j.IsObj[lastIndex]
	var b byte
	if isObj {
		b = '}'
	} else {
		b = ']'
	}
	j.IsObj = j.IsObj[:lastIndex]
	if _, err := j.Writer.Write([]byte{b}); err != nil {
		return err
	}
	j.First = false
	return nil
}

func (j *T) String() string {
	if s, ok := j.Writer.(fmt.Stringer); ok {
		return s.String()
	} else if m, ok := j.Writer.(json.Marshaler); ok {
		b, err := m.MarshalJSON()
		if err != nil {
			panic(fmt.Errorf("MarshalJSON on io.Writer failed: %w", err))
		}
		return string(b)
	}
	panic(fmt.Errorf("%T neither implements fmt.Stringer nor json.Marshaler", j.Writer))
}

func (j *T) Pretty() string { return j.Indent("", "  ") }

func (j *T) Indent(prefix, indent string) string {
	s := j.String()
	b := unsafe.Slice(unsafe.StringData(s), len(s))
	buf := bytes.NewBuffer(nil)
	if err := json.Indent(buf, b, prefix, indent); err != nil {
		panic(err)
	}
	return buf.String()
}

func (j *T) Close() error {
	if closer, ok := j.Writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
