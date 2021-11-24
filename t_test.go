package oas

import (
	"net/http"
	r "reflect"
	"testing"
	"time"
	u "unsafe"
)

func TestSchemaOf(t *testing.T) {
	test := func(expSchema Schema, expSchemas Schemas, typ interface{}) {
		t.Helper()
		var doc Doc
		eq(t, expSchema, doc.Sch(typ))
		eq(t, expSchemas, doc.Comps.Schemas)
	}

	test(
		Schema{Title: `string`, Type: []string{TypeStr}},
		nil,
		``,
	)

	test(
		Schema{Title: `*string`, Type: []string{TypeStr, TypeNull}},
		nil,
		(*string)(nil),
	)

	test(
		Schema{Title: `**string`, Type: []string{TypeStr, TypeNull}},
		nil,
		(**string)(nil),
	)

	test(
		Schema{Title: `int`, Type: []string{TypeInt}},
		nil,
		0,
	)

	test(
		Schema{Title: `*int`, Type: []string{TypeInt, TypeNull}},
		nil,
		(*int)(nil),
	)

	test(
		Schema{Title: `oas.Str`, Type: []string{TypeStr}},
		nil,
		Str(``),
	)

	test(
		Schema{Title: `*oas.Str`, Type: []string{TypeStr, TypeNull}},
		nil,
		(*Str)(nil),
	)

	test(
		Schema{Title: `oas.NullStr`, Type: []string{TypeStr, TypeNull}},
		nil,
		NullStr(``),
	)

	test(
		Schema{Title: `*oas.NullStr`, Type: []string{TypeStr, TypeNull}},
		nil,
		(*NullStr)(nil),
	)

	test(
		Schema{Title: `oas.NonNullStr`, Type: []string{TypeStr}},
		nil,
		NonNullStr(``),
	)

	/**
	Even though the type always encodes a non-null string, even when the pointer
	is nil, we mimic the behavior of "encoding/json" which considers nil
	pointers to be automatically null, which makes pointer types always nullable.
	*/
	test(
		Schema{Title: `*oas.NonNullStr`, Type: []string{TypeStr, TypeNull}},
		nil,
		(*NonNullStr)(nil),
	)

	test(
		Schema{Title: `oas.IntStr`, Type: []string{TypeNum}},
		nil,
		IntStr(``),
	)

	test(
		Schema{Title: `*oas.IntStr`, Type: []string{TypeNum, TypeNull}},
		nil,
		(*IntStr)(nil),
	)

	test(
		Schema{Title: `oas.IntStrPtr`, Type: []string{TypeStr}},
		nil,
		IntStrPtr(``),
	)

	test(
		Schema{Title: `*oas.IntStrPtr`, Type: []string{TypeNum, TypeNull}},
		nil,
		(*IntStrPtr)(nil),
	)

	test(
		RefSchema(`oas.WrapStr`),
		Schemas{
			`oas.WrapStr`: {
				Title: `oas.WrapStr`,
				Type:  []string{TypeObj},
				Props: Schemas{`Str`: {Title: `oas.Str`, Type: []string{TypeStr}}},
			},
		},
		WrapStr{},
	)

	test(
		NullSchema(`*oas.WrapStr`, RefSchema(`oas.WrapStr`)),
		Schemas{
			`oas.WrapStr`: {
				Title: `oas.WrapStr`,
				Type:  []string{TypeObj},
				Props: Schemas{`Str`: {Title: `oas.Str`, Type: []string{TypeStr}}},
			},
		},
		(*WrapStr)(nil),
	)

	test(
		Schema{Title: `oas.WrapNullStr`, Type: []string{TypeStr, TypeNull}},
		nil,
		WrapNullStr{},
	)

	test(
		Schema{Title: `*oas.WrapNullStr`, Type: []string{TypeStr, TypeNull}},
		nil,
		(*WrapNullStr)(nil),
	)

	test(
		Schema{Title: `time.Time`, Type: []string{TypeStr}, Format: FormatDateTime},
		nil,
		time.Time{},
	)

	test(
		Schema{Title: `*time.Time`, Type: []string{TypeStr, TypeNull}, Format: FormatDateTime},
		nil,
		(*time.Time)(nil),
	)

	test(
		Schema{Title: `oas.TerByte`, Type: []string{TypeBool, TypeNull}},
		nil,
		TerByte(0),
	)

	test(
		Schema{Title: `*oas.TerByte`, Type: []string{TypeBool, TypeNull}},
		nil,
		(*TerByte)(nil),
	)

	test(
		Schema{Title: `oas.Ter8`, Type: []string{TypeBool, TypeNull}},
		nil,
		Ter8{},
	)

	test(
		Schema{Title: `*oas.Ter8`, Type: []string{TypeBool, TypeNull}},
		nil,
		(*Ter8)(nil),
	)

	test(
		Schema{Title: `oas.Ter32`, Type: []string{TypeBool, TypeNull}},
		nil,
		Ter32(0),
	)

	test(
		Schema{Title: `*oas.Ter32`, Type: []string{TypeBool, TypeNull}},
		nil,
		(*Ter32)(nil),
	)

	test(
		RefSchema(`[]string`),
		Schemas{
			`[]string`: {
				Title: `[]string`,
				Type:  []string{TypeArr, TypeNull},
				Items: &Schema{Title: `string`, Type: []string{TypeStr}},
			},
		},
		[]string(nil),
	)

	test(
		RefSchema(`[]string`),
		Schemas{
			`[]string`: {
				Title: `[]string`,
				Type:  []string{TypeArr, TypeNull},
				Items: &Schema{Title: `string`, Type: []string{TypeStr}},
			},
		},
		(*[]string)(nil),
	)

	test(
		RefSchema(`[]*string`),
		Schemas{
			`[]*string`: {
				Title: `[]*string`,
				Type:  []string{TypeArr, TypeNull},
				Items: &Schema{Title: `*string`, Type: []string{TypeStr, TypeNull}},
			},
		},
		([]*string)(nil),
	)

	test(
		Schema{Title: `oas.Uuid`, Type: []string{TypeStr}, Format: FormatUuid},
		nil,
		Uuid{},
	)

	test(
		Schema{Title: `*oas.Uuid`, Type: []string{TypeStr, TypeNull}, Format: FormatUuid},
		nil,
		(*Uuid)(nil),
	)

	test(
		Schema{Title: `oas.NullUuid`, Type: []string{TypeStr, TypeNull}, Format: FormatUuid},
		nil,
		NullUuid{},
	)

	test(
		Schema{Title: `*oas.NullUuid`, Type: []string{TypeStr, TypeNull}, Format: FormatUuid},
		nil,
		(*NullUuid)(nil),
	)

	test(
		Schema{Title: `oas.NullTime`, Type: []string{TypeStr, TypeNull}, Format: FormatDateTime},
		nil,
		NullTime{},
	)

	test(
		Schema{Title: `*oas.NullTime`, Type: []string{TypeStr, TypeNull}, Format: FormatDateTime},
		nil,
		(*NullTime)(nil),
	)

	test(
		RefSchema(`oas.Unit`),
		Schemas{
			`oas.Unit`: {
				Title: `oas.Unit`,
				Type:  []string{TypeObj},
				Props: Schemas{
					`one_json`: {Title: `string`, Type: []string{TypeStr}},
				},
			},
		},
		Unit{},
	)

	test(
		NullSchema(`*oas.Unit`, RefSchema(`oas.Unit`)),
		Schemas{
			`oas.Unit`: {
				Title: `oas.Unit`,
				Type:  []string{TypeObj},
				Props: Schemas{
					`one_json`: {Title: `string`, Type: []string{TypeStr}},
				},
			},
		},
		(*Unit)(nil),
	)

	test(
		RefSchema(`oas.UnitWith`),
		Schemas{
			`oas.UnitWith`: {
				Title: `oas.UnitWith`,
				Type:  []string{TypeObj},
				Props: Schemas{
					`one_json`: {Title: `string`, Type: []string{TypeStr}},
					`Untagged`: {Title: `int`, Type: []string{TypeInt}},
				},
			},
		},
		UnitWith{},
	)

	test(
		NullSchema(`*oas.UnitWith`, RefSchema(`oas.UnitWith`)),
		Schemas{
			`oas.UnitWith`: {
				Title: `oas.UnitWith`,
				Type:  []string{TypeObj},
				Props: Schemas{
					`one_json`: {Title: `string`, Type: []string{TypeStr}},
					`Untagged`: {Title: `int`, Type: []string{TypeInt}},
				},
			},
		},
		(*UnitWith)(nil),
	)

	test(
		RefSchema(`oas.Pair`),
		Schemas{
			`oas.Pair`: {
				Title: `oas.Pair`,
				Type:  []string{TypeObj},
				Props: Schemas{
					`one_json`: {Title: `string`, Type: []string{TypeStr}},
					`two_json`: {Title: `int`, Type: []string{TypeInt}},
				},
			},
		},
		Pair{},
	)

	test(
		NullSchema(`*oas.Pair`, RefSchema(`oas.Pair`)),
		Schemas{
			`oas.Pair`: {
				Title: `oas.Pair`,
				Type:  []string{TypeObj},
				Props: Schemas{
					`one_json`: {Title: `string`, Type: []string{TypeStr}},
					`two_json`: {Title: `int`, Type: []string{TypeInt}},
				},
			},
		},
		(*Pair)(nil),
	)

	test(
		RefSchema(`oas.Outer`),
		outerSchemas(),
		Outer{},
	)

	test(
		NullSchema(`*oas.Outer`, RefSchema(`oas.Outer`)),
		outerSchemas(),
		(*Outer)(nil),
	)

	test(
		RefSchema(`map[string]int`),
		Schemas{
			`map[string]int`: {
				Title:    `map[string]int`,
				Type:     []string{TypeObj, TypeNull},
				AddProps: &Schema{Title: `int`, Type: []string{TypeInt}},
			},
		},
		map[string]int(nil),
	)

	test(
		RefSchema(`map[string]int`),
		Schemas{
			`map[string]int`: {
				Title:    `map[string]int`,
				Type:     []string{TypeObj, TypeNull},
				AddProps: &Schema{Title: `int`, Type: []string{TypeInt}},
			},
		},
		(*map[string]int)(nil),
	)
}

func TestDoc_Route(t *testing.T) {
	var doc Doc
	doc.Route(`/`, http.MethodGet, Op{ReqBody: doc.JsonBodyOpt(Outer{})})

	eq(
		t,
		Doc{
			Paths: Paths{
				`/`: Path{
					Get: &Op{
						ReqBody: &Body{
							Cont: MediaTypes{
								ConTypeJson: MediaType{
									Schema: RefSchema(`oas.Outer`),
								},
							},
						},
					},
				},
			},
			Comps: Comps{Schemas: outerSchemas()},
		},
		doc,
	)
}

func TestVerboseDocJson(t *testing.T) {
	if !testing.Verbose() {
		t.Skip(`run in verbose mode`)
	}
	writeFile(`mock.json`, jsonStr(tDoc()))
}

func tDoc() Doc {
	doc := Doc{
		Openapi: Ver,
		Info: &Info{
			Title: `API documentation for my server`,
			Desc: `
		Documentation in JSON or YAML format,
		compatible with the OpenAPI specification.
		`,
			Ver: `v3`,
		},
	}
	doc.Route(`/ents`, http.MethodPost, Op{ReqBody: doc.JsonBodyOpt(Inner{})})
	doc.Route(`/ents/{}`, http.MethodGet, Op{ReqBody: doc.JsonBodyOpt(Outer{})})
	return doc
}

func Test_nonZero(t *testing.T) {
	test := func(ok bool, exp interface{}) {
		t.Helper()
		val := r.New(r.TypeOf(exp)).Elem()
		eq(t, ok, nonZero(val))
		eq(t, exp, val.Interface())
	}

	test(false, struct{}{})
	test(false, [0]string{})
	test(false, [2]struct{}{})
	test(false, []struct{}(nil))
	test(false, struct {
		_ struct{}
		_ struct{}
		_ struct{}
	}{})
	test(false, chan struct{}(nil))
	test(false, map[string]struct{}(nil))
	test(false, map[struct{}]struct{}(nil))
	test(false, map[struct{}]string(nil))

	test(true, int8(1))
	test(true, int16(1))
	test(true, int32(1))
	test(true, int64(1))
	test(true, int(1))
	test(true, uint8(1))
	test(true, uint16(1))
	test(true, uint32(1))
	test(true, uint64(1))
	test(true, uint(1))
	test(true, true)
	test(true, ` `)

	test(true, [1]int{1})
	test(true, [1]bool{true})
	test(true, [1]string{` `})

	test(true, [2]int{1})
	test(true, [2]bool{true})
	test(true, [2]string{` `})

	test(true, []int{1})
	test(true, []bool{true})
	test(true, []string{` `})

	test(true, [][]int{{1}})
	test(true, [][]bool{{true}})
	test(true, [][]string{{` `}})

	test(true, [][][2]int{{{1}}})
	test(true, [][][2]bool{{{true}}})
	test(true, [][][2]string{{{` `}}})

	test(true, map[int]string{1: ` `})
	test(true, map[int][]string{1: []string{` `}})
	test(true, map[[2]int][]string{[2]int{1}: []string{` `}})

	// nolint:structcheck
	test(true, struct{ private int }{1})

	test(true, struct {
		Void struct{}
		Val  int
	}{
		Void: struct{}{},
		Val:  1,
	})

	test(true, intPtr(1))
	test(true, boolPtr(true))
	test(true, stringPtr(` `))

	test(true, &[2]int{1})
	test(true, &[2]bool{true})
	test(true, &[2]string{` `})

	// nolint:structcheck
	test(true, struct {
		private int
		Void    struct{}
		Inner   *struct{ Val map[int]*[]*string }
	}{
		Inner: &struct{ Val map[int]*[]*string }{Val: map[int]*[]*string{
			1: &[]*string{stringPtr(` `)},
		}},
	})
}

func Test_memcpy(t *testing.T) {
	var tar [4]byte
	var src [4]byte = [4]byte{10, 20, 30, 40}

	memcpy(
		uintptr(u.Pointer(&tar)),
		uintptr(u.Pointer(&src)),
		2,
	)

	eq(t, [4]byte{10, 20, 0, 0}, tar)
	eq(t, [4]byte{10, 20, 30, 40}, src)
}
