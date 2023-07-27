package oas

import (
	"encoding"
	"encoding/json"
	r "reflect"
	"strings"
)

/*
This file contains internal functions related to schema generation,
registration, and modification. It's separate to keep the "public" API more
easily browsable.
*/

func (self *Doc) schemaAny(sch *Schema, typ r.Type) {
	if typ == nil {
		self.schemaNone(sch, typ)
		return
	}

	name := typeName(typ)
	_, ok := self.GotCompSchema(name)
	if ok {
		sch.setRef(name)
		return
	}

	self.schemaCommon(sch, typ)
	if self.schemaIfaces(sch, typ) {
		return
	}

	switch typ.Kind() {
	case r.Int32, r.Uint32:
		sch.Format = FormatInt32
		self.schemaInt(sch, typ)

	case r.Int64, r.Uint64:
		sch.Format = FormatInt64
		self.schemaInt(sch, typ)

	case r.Int8, r.Int16, r.Int, r.Uint8, r.Uint16, r.Uint:
		self.schemaInt(sch, typ)

	case r.Float32:
		sch.Format = FormatFloat32
		self.schemaFloat(sch, typ)

	case r.Float64:
		sch.Format = FormatFloat64
		self.schemaFloat(sch, typ)

	case r.Bool:
		self.schemaBool(sch, typ)

	case r.String:
		self.schemaString(sch, typ)

	case r.Ptr:
		self.schemaPtr(sch, typ)

	case r.Array:
		self.schemaArray(sch, typ)

	case r.Slice:
		self.schemaSlice(sch, typ)

	case r.Map:
		self.schemaMap(sch, typ)

	case r.Struct:
		self.schemaStruct(sch, typ)

	default:
		panic(errSchemaUnsupported(typ))
	}
}

func (*Doc) schemaNone(sch *Schema, _ r.Type)   { sch.Nullable() }
func (*Doc) schemaInt(sch *Schema, _ r.Type)    { sch.Type = []string{TypeInt} }
func (*Doc) schemaFloat(sch *Schema, _ r.Type)  { sch.Type = []string{TypeNum} }
func (*Doc) schemaBool(sch *Schema, _ r.Type)   { sch.Type = []string{TypeBool} }
func (*Doc) schemaString(sch *Schema, _ r.Type) { sch.Type = []string{TypeStr} }

func (self *Doc) schemaPtr(sch *Schema, typ r.Type) {
	self.schemaAny(sch, typ.Elem())

	if sch.Ref == `` {
		sch.Title = typeName(typ)
		sch.Nullable()
		return
	}

	refPath := sch.Ref
	tar, ok := self.GotSchema(refPath)
	if !ok {
		panic(errSchemaMissing(refPath))
	}

	/**
	If the target type is inherently nullable, we simply pretend that the pointer
	is not there. This makes the title inconsistent with the cases when
	nullability is the result of using a pointer, but simplifies the output
	structure, avoiding unnecessary references.
	*/
	if tar.IsNullable() {
		return
	}

	*sch = NullSchema(typeName(typ), *sch)
}

func (self *Doc) schemaArray(sch *Schema, typ r.Type) {
	name := typeName(typ)
	defer self.setSchema(name, Schema{}).outlineSchema(sch)

	sch.MaxItems = uint64(typ.Len())
	sch.MinItems = uint64(typ.Len())
	sch.Items = self.TypeSchema(typ.Elem()).Opt()
}

func (self *Doc) schemaSlice(sch *Schema, typ r.Type) {
	name := typeName(typ)
	defer self.setSchema(name, Schema{}).outlineSchema(sch)

	sch.Type = []string{TypeArr, TypeNull}
	sch.Items = self.TypeSchema(typ.Elem()).Opt()
}

func (self *Doc) schemaMap(sch *Schema, typ r.Type) {
	name := typeName(typ)
	defer self.setSchema(name, Schema{}).outlineSchema(sch)

	keyType := typ.Key()
	elemType := typ.Elem()

	sch.Type = []string{TypeObj, TypeNull}

	if isTypeSkippable(elemType) {
		self.schemaNone(sch, typ)
		return
	}

	validKeyFor(typ, keyType, self.TypeSchema(keyType))
	sch.AddProps = self.TypeSchema(elemType).Opt()
}

func (self *Doc) schemaStruct(sch *Schema, typ r.Type) {
	name := typeName(typ)
	defer self.setSchema(name, Schema{}).outlineSchema(sch)

	sch.Type = []string{TypeObj}
	self.schemaStructProps(sch, typ)
}

func (self *Doc) schemaStructProps(sch *Schema, typ r.Type) {
	for ind := range iter(typ.NumField()) {
		field := typ.Field(ind)

		if !isPublic(field.PkgPath) || isTypeSkippable(field.Type) {
			continue
		}

		name := jsonName(field)
		if name != `` {
			self.schemaStructProp(sch, name, field.Type)
			continue
		}

		if field.Anonymous {
			inner := typeDeref(field.Type)
			if inner.Kind() == r.Struct {
				self.schemaStructProps(sch, inner)
				continue
			}
		}

		name = field.Name
		if name != `` {
			self.schemaStructProp(sch, name, field.Type)
		}
	}
}

func (self *Doc) schemaStructProp(sch *Schema, name string, typ r.Type) {
	sch.Props.Init()[name] = self.TypeSchema(typ)
}

func (self *Doc) schemaCommon(sch *Schema, typ r.Type) {
	self.schemaTitle(sch, typ)
}

func (*Doc) schemaTitle(sch *Schema, typ r.Type) {
	val := typeName(typ)
	if val != `` {
		sch.Title = val
		return
	}
}

func (self *Doc) schemaIfaces(sch *Schema, typ r.Type) bool {
	if typ.Implements(ifaceJsonMarshaler) {
		return self.schemaIfaceJson(sch, typ)
	}
	if typ.Implements(ifaceTextMarshaler) {
		return self.schemaIfaceText(sch, typ)
	}
	return false
}

func (self *Doc) schemaIfaceJson(sch *Schema, typ r.Type) bool {
	for typ.Kind() == r.Ptr {
		/**
		Mimic the behavior of "encoding/json". It considers nil pointers to be
		automatically null, without trying to invoke their `json.Marshaler`,
		because if the method is actually implemented on the value type, invoking
		it on a nil pointer will panic, and there's NO WAY to detect that without
		calling the method and getting a panic. This means every pointer type is
		nullable no matter what, even if encoding method is actually implemented
		on a pointer type (not value type) and has a nil pointer check, returning
		something custom for a nil pointer. That case is never invoked
		by "encoding/json".
		*/
		sch.TypeAdd(TypeNull)
		typ = typ.Elem()
	}

	typ = typeDeref(typ)
	val := r.New(typ)
	if self.schemaJsonVal(sch, val) {
		return true
	}

	return nonZero(val.Elem()) && self.schemaJsonVal(sch, val)
}

func (self *Doc) schemaJsonVal(sch *Schema, val r.Value) bool {
	chunk, err := toJson(val.Convert(ifaceJsonMarshaler).Interface().(json.Marshaler))
	return err == nil && self.schemaJsonInspect(sch, bytesString(chunk))
}

/*
TODO: consider supporting the entire JSON syntax. Missing features:

	* Detecting list types and their element types, recursively.
	  (Stop at the first element).

	* Detecting dict types and their element types, recursively.
	  (Stop at the first element).

Inspecting lists and dicts must be done ONLY when the JSON kind doesn't match
the Go kind. For example, given a struct that encodes as a JSON list, we're
better off inspecting its JSON output. But given a slice that implements custom
JSON marshaling but nevertheless encodes as a list, we should probably skip
JSON inspection and inspect it like any other Go slice, because that will give
us more information about its element type.
*/
func (self *Doc) schemaJsonInspect(sch *Schema, val string) bool {
	val = strings.TrimSpace(val)

	if val == `null` {
		sch.Nullable()
		return false
	}

	if val == `true` || val == `false` {
		sch.TypeAdd(TypeBool)
		sch.Format = ``
		return true
	}

	if len(val) > 0 && val[0] == '"' {
		sch.TypeAdd(TypeStr)
		self.schemaTextInspectFormat(sch, unquote(val))
		return true
	}

	if len(val) > 0 && (val[0] == '-' || isDecDigit(val[0])) {
		sch.TypeAdd(TypeNum)
		if strings.ContainsAny(val, `.eE`) {
			sch.Format = FormatFloat64
		}
		return true
	}

	return false
}

func (self *Doc) schemaIfaceText(sch *Schema, typ r.Type) bool {
	// See the comment on `(*Doc).schemaIfaceJson` for the why.
	for typ.Kind() == r.Ptr {
		sch.TypeAdd(TypeNull)
		typ = typ.Elem()
	}

	val := r.New(typ)
	if self.schemaTextVal(sch, val) || typ.Size() == 0 {
		return true
	}

	return nonZero(val.Elem()) && self.schemaTextVal(sch, val)
}

func (self *Doc) schemaTextVal(sch *Schema, val r.Value) bool {
	chunk, err := toText(val.Convert(ifaceTextMarshaler).Interface().(encoding.TextMarshaler))
	if err != nil {
		return false
	}
	self.schemaTextInspect(sch, bytesString(chunk))
	return true
}

func (self *Doc) schemaTextInspect(sch *Schema, val string) {
	sch.TypeAdd(TypeStr)
	self.schemaTextInspectFormat(sch, val)
}

func (*Doc) schemaTextInspectFormat(sch *Schema, val string) {
	val = strings.TrimSpace(val)

	if isDateTimeRfc3339(val) {
		sch.Format = FormatDateTime
		return
	}

	if isDateIso8601(val) {
		sch.Format = FormatDate
		return
	}

	if isTimeIso8601ExtendedT(val) || isTimeIso8601Extended(val) {
		sch.Format = FormatTime
		return
	}

	if isUuid(val) {
		sch.Format = FormatUuid
		return
	}

	if isDurationIso8601(val) {
		sch.Format = FormatDuration
		return
	}
}

func (self *Doc) setSchema(name string, sch Schema) *Doc {
	if name == `` {
		panic(errMissingTitle)
	}

	_, ok := self.Comps.Schemas[name]
	if ok {
		panic(errSchemaRedundant(name))
	}

	self.Comps.Schemas.Init()[name] = sch
	return self
}

// Opposite of "inline". Term borrowed from compiler lingo.
func (self *Doc) outlineSchema(sch *Schema) {
	self.addSchema(*sch)
	sch.setRef(sch.ValidTitle())
}

func (self *Doc) addSchema(sch Schema) {
	self.Comps.Schemas.Init()[sch.ValidTitle()] = sch
}
