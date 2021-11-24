package oas

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	r "reflect"
)

type Str string

type NullStr string

func (self NullStr) MarshalJSON() ([]byte, error) {
	if self == `` {
		return []byte(`null`), nil
	}
	return json.Marshal(string(self))
}

type NonNullStr string

func (self *NonNullStr) MarshalJSON() ([]byte, error) {
	if self == nil {
		return json.Marshal(``)
	}
	return json.Marshal(string(*self))
}

type IntStr string

func (IntStr) MarshalJSON() ([]byte, error) {
	return json.Marshal(-123)
}

type IntStrPtr string

func (*IntStrPtr) MarshalJSON() ([]byte, error) {
	return json.Marshal(-123)
}

type WrapStr struct{ Str }

type WrapNullStr struct{ NullStr }

type TerByte byte

func (self TerByte) MarshalJSON() ([]byte, error) {
	if self == 0 {
		return []byte(`null`), nil
	}
	if self == 1 {
		return []byte(`false`), nil
	}
	if self == 2 {
		return []byte(`true`), nil
	}
	return nil, fmt.Errorf(`invalid value of %[1]T: %[1]v`, self)
}

type Ter8 [1]byte

func (self Ter8) MarshalJSON() ([]byte, error) {
	if self == (Ter8{}) {
		return []byte(`null`), nil
	}
	if self == (Ter8{1}) {
		return []byte(`false`), nil
	}
	if self == (Ter8{2}) {
		return []byte(`true`), nil
	}
	return nil, fmt.Errorf(`invalid value of %[1]T: %[1]v`, self)
}

type Ter32 int32

func (self Ter32) MarshalJSON() ([]byte, error) {
	if self == 0 {
		return []byte(`null`), nil
	}
	if self == 1 {
		return []byte(`false`), nil
	}
	if self == 2 {
		return []byte(`true`), nil
	}
	return nil, fmt.Errorf(`invalid value of %[1]T: %[1]v`, self)
}

type Uuid [16]byte

func (self Uuid) String() string {
	var buf [32]byte
	hex.Encode(buf[:], self[:])
	return string(buf[:])
}

func (self Uuid) MarshalText() ([]byte, error) {
	var buf [32]byte
	hex.Encode(buf[:], self[:])
	return buf[:], nil
}

type NullUuid Uuid

func (self NullUuid) IsZero() bool { return self == NullUuid{} }

func (self NullUuid) String() string {
	if self.IsZero() {
		return ``
	}
	return Uuid(self).String()
}

func (self NullUuid) MarshalText() ([]byte, error) {
	if self.IsZero() {
		return nil, nil
	}
	return Uuid(self).MarshalText()
}

func (self NullUuid) MarshalJSON() ([]byte, error) {
	if self.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(Uuid(self))
}

type NullTime time.Time

func (self NullTime) IsZero() bool { return self == NullTime{} }

func (self NullTime) MarshalJSON() ([]byte, error) {
	if self.IsZero() {
		return []byte(`null`), nil
	}
	return json.Marshal(time.Time(self))
}

type Unit struct {
	One string `json:"one_json" db:"one_db"`
}

// nolint:structcheck,unused,govet
type UnitWith struct {
	Untagged      int
	One           string        `json:"one_json" db:"one_db"`
	private       bool          `json:"private_json" db:"private_db"`
	SkippableChan chan struct{} `json:"skippable_chan_json" db:"skippable_chan_db"`
	SkippableFunc func()        `json:"skippable_func_json" db:"skippable_func_db"`
}

type Pair struct {
	One string `json:"one_json" db:"one_db"`
	Two int    `json:"two_json" db:"two_db"`
}

type Listed struct {
	Scalar Uuid   `json:"scalar"`
	List   []Unit `json:"list"`
}

type Outer struct {
	Embed
	OuterOne   string `json:"outer_one"`
	OuterInner *Inner `json:"outer_inner"`
	OuterSlice []Pair `json:"outer_slice"`
}

type Inner struct {
	InnerTwo string `json:"inner_two"`
}

type Embed struct {
	EmbedThree NullUuid `json:"embed_three"`
}

func outerSchemas() Schemas {
	return Schemas{
		`oas.Outer`: {
			Title: `oas.Outer`,
			Type:  []string{TypeObj},
			Props: Schemas{
				`embed_three`: {
					Title:  `oas.NullUuid`,
					Type:   []string{TypeStr, TypeNull},
					Format: FormatUuid,
				},
				`outer_one`:   Schema{Title: `string`, Type: []string{TypeStr}},
				`outer_inner`: NullSchema(`*oas.Inner`, RefSchema(`oas.Inner`)),
				`outer_slice`: RefSchema(`[]oas.Pair`),
			},
		},
		`oas.Inner`: {
			Title: `oas.Inner`,
			Type:  []string{TypeObj},
			Props: Schemas{
				`inner_two`: {Title: `string`, Type: []string{TypeStr}},
			},
		},
		`oas.Pair`: {
			Title: `oas.Pair`,
			Type:  []string{TypeObj},
			Props: Schemas{
				`one_json`: {Title: `string`, Type: []string{TypeStr}},
				`two_json`: {Title: `int`, Type: []string{TypeInt}},
			},
		},
		`[]oas.Pair`: {
			Title: `[]oas.Pair`,
			Type:  []string{TypeArr, TypeNull},
			Items: RefSchema(`oas.Pair`).Opt(),
		},
	}
}

func eq(t testing.TB, exp, act interface{}) {
	t.Helper()
	if !r.DeepEqual(exp, act) {
		t.Fatalf(`
expected (simple):
	%[1]v
actual (simple):
	%[2]v
expected (JSON):
	%[3]s
actual (JSON):
	%[4]s
`, exp, act, jsonStr(exp), jsonStr(act))
	}
}

func jsonStr(val interface{}) string {
	chunk, err := json.MarshalIndent(val, ``, `  `)
	try(err)
	return string(chunk)
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func writeFile(path, body string) {
	try(os.WriteFile(path, []byte(body), os.ModePerm))
}

func intPtr(val int) *int          { return &val }
func stringPtr(val string) *string { return &val }
func boolPtr(val bool) *bool       { return &val }
