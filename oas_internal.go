package oas

import (
	"encoding"
	"encoding/json"
	"fmt"
	r "reflect"
	"strconv"
	"strings"
	"time"
	u "unsafe"
)

const (
	formatDateIso8601          = `2006-01-02`
	formatTimeIso8601ExtendedT = `T15:04:05`
	formatTimeIso8601Extended  = `15:04:05`
)

var (
	ifaceTextMarshaler = r.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	ifaceJsonMarshaler = r.TypeOf((*json.Marshaler)(nil)).Elem()
)

func typeDeref(typ r.Type) r.Type {
	for typ != nil && typ.Kind() == r.Ptr {
		typ = typ.Elem()
	}
	return typ
}

var (
	errMissingTitle = fmt.Errorf(`[oas] missing schema title`)
)

func errSchemaUnsupported(typ r.Type) error {
	return fmt.Errorf(
		`[oas] can't generate schema for type %q of kind %q`,
		typ, typ.Kind(),
	)
}

func errSchemaMissing(name string) error {
	return fmt.Errorf(`[oas] missing schema component %q`, name)
}

func errSchemaUnexpectedRef(name, ref string) error {
	return fmt.Errorf(
		`[oas] double indirection: schema referenced by %q unexpectedly has reference %q`,
		name, ref,
	)
}

func errSchemaRedundant(name string) error {
	return fmt.Errorf(`[oas] redundant schema %q`, name)
}

func validKeyFor(mapType, keyType r.Type, keySch Schema) {
	if !keySch.TypeIs(TypeStr) {
		panic(fmt.Errorf(
			`[oas] can't generate schema for map type %q: key type %q has representation type %q instead of required %q`,
			mapType, keyType, keySch.Type, TypeStr,
		))
	}
}

func isDateTimeRfc3339(val string) bool {
	_, err := time.Parse(time.RFC3339, val)
	return err == nil
}

func isDateIso8601(val string) bool {
	_, err := time.Parse(formatDateIso8601, val)
	return err == nil
}

func isTimeIso8601ExtendedT(val string) bool {
	_, err := time.Parse(formatTimeIso8601ExtendedT, val)
	return err == nil
}

func isTimeIso8601Extended(val string) bool {
	_, err := time.Parse(formatTimeIso8601Extended, val)
	return err == nil
}

func isUuid(val string) bool {
	return isUuidCanon(val) || isUuidSimple(val)
}

func isUuidCanon(val string) bool {
	if len(val) != (len(uuidCanonHyphens) + len(uuidCanonDigits)) {
		return false
	}
	for _, ind := range uuidCanonHyphens {
		if val[ind] != '-' {
			return false
		}
	}
	for _, ind := range uuidCanonDigits {
		if !isHexDigit(val[ind]) {
			return false
		}
	}
	return true
}

var uuidCanonDigits = [...]int{
	0, 1, 2, 3, 4, 5, 6, 7,
	9, 10, 11, 12,
	14, 15, 16, 17,
	19, 20, 21, 22,
	24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35,
}

var uuidCanonHyphens = [...]int{8, 13, 18, 23}

func isUuidSimple(val string) bool {
	if len(val) != 32 {
		return false
	}
	for _, char := range []byte(val) {
		if !isHexDigit(char) {
			return false
		}
	}
	return true
}

func isDurationIso8601(val string) bool { return durations[val] }

var durations = map[string]bool{
	// Possible encodings of zero duration.
	`P`:    true,
	`P0`:   true,
	`P0Y`:  true,
	`PT0S`: true,

	// Possible encodings of non-zero duration tweaked by `nonZero`.
	`PT1S`: true,
	`P1Y`:  true,
}

func isTypeSkippable(typ r.Type) bool {
	if typ == nil {
		return true
	}

	switch typ.Kind() {
	case r.Chan, r.Func, r.Interface, r.UnsafePointer:
		return true
	case r.Array, r.Slice, r.Map, r.Ptr:
		return isTypeSkippable(typ.Elem())
	default:
		return false
	}
}

func isDecDigit(val byte) bool { return decDigits[val] }

/*
Could probably be replaced with `(val >= '0' && val <= '9')`,
but I don't trust myself to remember if 0 comes before 9.
*/
var decDigits = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
}

func isHexDigit(val byte) bool { return hexDigits[val] }

var hexDigits = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true,
}

/*
Allocation-free conversion. Reinterprets a byte slice as a string. Borrowed from
the standard library. Reasonably safe. Should not be used when the underlying
byte array is volatile, for example part of a scratch buffer in SQL scanning.
*/
func bytesString(val []byte) string {
	return *(*string)(u.Pointer(&val))
}

func unquote(src string) string {
	out, err := strconv.Unquote(src)
	if err != nil {
		return src
	}
	return out
}

func nonZero(val r.Value) bool {
	typ := val.Type()
	if typ.Size() == 0 {
		return false
	}

	switch typ.Kind() {
	case r.Int8, r.Int16, r.Int32, r.Int64, r.Int:
		val.SetInt(1)
		return true

	case r.Uint8, r.Uint16, r.Uint32, r.Uint64, r.Uint:
		val.SetUint(1)
		return true

	case r.Float32, r.Float64:
		val.SetFloat(1)
		return true

	case r.Bool:
		val.SetBool(true)
		return true

	case r.String:
		val.SetString(` `)
		return true

	case r.Array:
		return nonZero(val.Index(0))

	case r.Map:
		return nonZeroMap(val)

	case r.Slice:
		return nonZeroSlice(val)

	case r.Struct:
		return nonZeroStruct(val)

	case r.Ptr:
		return nonZeroPtr(val)

	default:
		return false
	}
}

func nonZeroMap(val r.Value) bool {
	typ := val.Type()
	key := r.New(typ.Key()).Elem()
	elem := r.New(typ.Elem()).Elem()

	if nonZero(key) && nonZero(elem) {
		val.Set(r.MakeMap(typ))
		val.SetMapIndex(key, elem)
		return true
	}

	return false
}

func nonZeroSlice(val r.Value) bool {
	typ := val.Type()
	elem := r.New(typ.Elem()).Elem()

	if nonZero(elem) {
		val.Set(r.Append(val, elem))
		return true
	}

	return false
}

func nonZeroStruct(val r.Value) bool {
	return nonZeroStructPublic(val) || nonZeroStructAny(val)
}

func nonZeroStructPublic(val r.Value) bool {
	typ := val.Type()

	for ind := range iter(typ.NumField()) {
		if !isPublic(typ.Field(ind).PkgPath) {
			continue
		}
		if nonZero(val.Field(ind)) {
			return true
		}
	}

	return false
}

func nonZeroStructAny(val r.Value) bool {
	typ := val.Type()

	for ind := range iter(typ.NumField()) {
		elem := r.New(typ.Field(ind).Type).Elem()

		if nonZero(elem) {
			/**
			The "clearer" alternative would be to unset the "unexported" flag in the
			field obtained by `val.Field(ind)` and call `nonZero` on it directly. But
			it seems to require hacks that involve magic constants, which may change
			between language versions:

				field := val.Field(ind)
				(*[3]uintptr)(u.Pointer(&field))[2] &^= uintptr((1 << 5) | (1 << 6))
				if nonZero(field) {return true}

			The memcpy approach seems less likely to break between language releases.
			Assumes indirection; `r.Value.UnsafeAddr` should panic if the value
			isn't indirect. If that ever happens, we'll need to figure out how to
			skip before even trying.
			*/
			memcpy(
				val.Field(ind).UnsafeAddr(),
				elem.UnsafeAddr(),
				elem.Type().Size(),
			)

			return true
		}
	}

	return false
}

func memcpy(tar, src, len uintptr) {
	copy(
		*(*[]byte)(u.Pointer(&[3]uintptr{tar, len, len})),
		*(*[]byte)(u.Pointer(&[3]uintptr{src, len, len})),
	)
}

func nonZeroPtr(val r.Value) bool {
	if val.IsNil() {
		val.Set(r.New(val.Type().Elem()))
	}
	return nonZero(val.Elem())
}

func typeName(typ r.Type) string {
	if typ == nil {
		return ``
	}

	switch typ.Kind() {
	case r.Struct:
		if typ.Name() == `` {
			panic(fmt.Errorf(`[oas] unexpected anonymous struct type %q`, typ))
		}
		return typ.String()

	default:
		return typ.String()
	}
}

func iter(count int) []struct{} { return make([]struct{}, count) }

func isPublic(pkgPath string) bool { return pkgPath == `` }

func jsonName(field r.StructField) string {
	return tagIdent(field.Tag.Get(`json`))
}

func tagIdent(tag string) string {
	index := strings.IndexRune(tag, ',')
	if index >= 0 {
		tag = tag[:index]
	}
	if tag == `-` {
		return ``
	}
	return tag
}

func someSchema(vals []Schema, fun func(Schema) bool) bool {
	if fun == nil {
		return false
	}
	for _, val := range vals {
		if fun(val) {
			return true
		}
	}
	return false
}

// Considers `[]string(nil) == []string{}` and that's intentional.
func stringsEq(one, two []string) bool {
	if len(one) != len(two) {
		return false
	}
	for ind := range one {
		if one[ind] != two[ind] {
			return false
		}
	}
	return true
}

func stringsContain(vals []string, exp string) bool {
	for _, val := range vals {
		if val == exp {
			return true
		}
	}
	return false
}

func toJson(val json.Marshaler) (_ []byte, err error) {
	return val.MarshalJSON()
}

func toText(val encoding.TextMarshaler) (_ []byte, err error) {
	return val.MarshalText()
}

func unprefix(base, prefix string) (string, bool) {
	if strings.HasPrefix(base, prefix) {
		return base[len(prefix):], true
	}
	return base, false
}
