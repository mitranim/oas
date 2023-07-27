package oas

import (
	"fmt"
	"path"
)

/*
Shortcut for making a reference-only schema pointing at
`#/components/schema/<name>`.
*/
// func RefSchema(name string) Schema { return Schema{Ref: SchemaRef(name)} }

func RefSchema(name string) (out Schema) {
	if name == `` {
		panic(errMissingTitle)
	}
	out.Ref = path.Join(`#/components/schemas`, name)
	return
}

/*
Shortcut for making a schema that wraps another and uses `.OneOf` with
`oas.TypeNull` to indicate nullability. Mostly for internal use; you should
never have to annotate nullability manually, as this package detects it
automatically.
*/
func NullSchema(name string, inner Schema) Schema {
	/**
	Note: this puts "null" in second position because OAS visualization tools may
	display the first element first, without trying to be "smart" and
	de-prioritize the null.
	*/
	return Schema{
		Title: name,
		OneOf: []Schema{inner, {Type: []string{TypeNull}}},
	}
}

/*
References:

	https://spec.openapis.org/oas/v3.1.0#schema-object

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00

The properties are listed in the order of their definition in the respective
specifications, on a best-effort basis.
*/
type Schema struct {
	// Ref `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	// Ref
	Ref  string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum  string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`

	/**
	OAS properties.
	*/

	Discr   *Discr  `json:"discriminator,omitempty" yaml:"discriminator,omitempty" toml:"discriminator,omitempty"`
	Xml     *Xml    `json:"xml,omitempty"           yaml:"xml,omitempty"           toml:"xml,omitempty"`
	ExtDoc  *ExtDoc `json:"externalDocs,omitempty"  yaml:"externalDocs,omitempty"  toml:"externalDocs,omitempty"`
	Example any     `json:"example,omitempty"       yaml:"example,omitempty"       toml:"example,omitempty"` // Also see `.Examples`.

	/**
	JSON Schema core properties. Reference:

		https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00
	*/

	// Subschemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.1
	AllOf []Schema `json:"allOf,omitempty" yaml:"allOf,omitempty" toml:"allOf,omitempty"`
	AnyOf []Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty" toml:"anyOf,omitempty"`
	OneOf []Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty" toml:"oneOf,omitempty"`
	Not   *Schema  `json:"not,omitempty"   yaml:"not,omitempty"   toml:"not,omitempty"`

	// Conditional subschemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.2
	If         *Schema `json:"if,omitempty"               yaml:"if,omitempty"               toml:"if,omitempty"`
	Then       *Schema `json:"then,omitempty"             yaml:"then,omitempty"             toml:"then,omitempty"`
	Else       *Schema `json:"else,omitempty"             yaml:"else,omitempty"             toml:"else,omitempty"`
	DepSchemas Schemas `json:"dependentSchemas,omitempty" yaml:"dependentSchemas,omitempty" toml:"dependentSchemas,omitempty"`

	// Array child schemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.1
	PrefixItems []Schema `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty" toml:"prefixItems,omitempty"`
	Items       *Schema  `json:"items,omitempty"       yaml:"items,omitempty"       toml:"items,omitempty"`
	Contains    *Schema  `json:"contains,omitempty"    yaml:"contains,omitempty"    toml:"contains,omitempty"`

	// Object subschemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.2
	Props     Schemas `json:"properties,omitempty"           yaml:"properties,omitempty"           toml:"properties,omitempty"`
	PatProps  Schemas `json:"patternProperties,omitempty"    yaml:"patternProperties,omitempty"    toml:"patternProperties,omitempty"`
	AddProps  *Schema `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty" toml:"additionalProperties,omitempty"`
	PropNames *Schema `json:"propertyNames,omitempty"        yaml:"propertyNames,omitempty"        toml:"propertyNames,omitempty"`

	// Unevaluated locations.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-11
	UnevalItems *Schema `json:"unevaluatedItems,omitempty"      yaml:"unevaluatedItems,omitempty"      toml:"unevaluatedItems,omitempty"`
	UnevalProps *Schema `json:"unevaluatedProperties,omitempty" yaml:"unevaluatedProperties,omitempty" toml:"unevaluatedProperties,omitempty"`

	/**
	JSON Schema validation properties. Reference:

		https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00
	*/

	// Validation for any instance.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.1
	Type  []string `json:"type,omitempty"  yaml:"type,omitempty"  toml:"type,omitempty"`
	Enum  []string `json:"enum,omitempty"  yaml:"enum,omitempty"  toml:"enum,omitempty"`
	Const any      `json:"const,omitempty" yaml:"const,omitempty"                       toml:"const,omitempty"`

	// Validation for numeric instances.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2
	MulOf   uint64 `json:"multipleOf,omitempty"       yaml:"multipleOf,omitempty"       toml:"multipleOf,omitempty"` // 0 represents "missing".
	Max     *int64 `json:"maximum,omitempty"          yaml:"maximum,omitempty"          toml:"maximum,omitempty"`
	ExlcMax *int64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty" toml:"exclusiveMaximum,omitempty"`
	Min     *int64 `json:"minimum,omitempty"          yaml:"minimum,omitempty"          toml:"minimum,omitempty"`
	ExclMin *int64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty" toml:"exclusiveMinimum,omitempty"`

	// Validation for strings.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.3
	MaxLen  uint64 `json:"maxLength,omitempty" yaml:"maxLength,omitempty" toml:"maxLength,omitempty"` // 0 represents "missing".
	MinLen  uint64 `json:"minLength,omitempty" yaml:"minLength,omitempty" toml:"minLength,omitempty"`
	Pattern string `json:"pattern,omitempty"   yaml:"pattern,omitempty"   toml:"pattern,omitempty"`

	// Validation for arrays.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4
	MaxItems  uint64 `json:"maxItems,omitempty"    yaml:"maxItems,omitempty"    toml:"maxItems,omitempty"`
	MinItems  uint64 `json:"minItems,omitempty"    yaml:"minItems,omitempty"    toml:"minItems,omitempty"`
	UniqItems bool   `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty" toml:"uniqueItems,omitempty"`
	MaxCont   uint64 `json:"maxContains,omitempty" yaml:"maxContains,omitempty" toml:"maxContains,omitempty"`
	MinCont   uint64 `json:"minContains,omitempty" yaml:"minContains,omitempty" toml:"minContains,omitempty"`

	// Validation for objects.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.5
	MaxProps uint64              `json:"maxProperties,omitempty"     yaml:"maxProperties,omitempty"     toml:"maxProperties,omitempty"`
	MinProps uint64              `json:"minProperties,omitempty"     yaml:"minProperties,omitempty"     toml:"minProperties,omitempty"`
	Requ     bool                `json:"required,omitempty"          yaml:"required,omitempty"          toml:"required,omitempty"`
	DepRequ  map[string][]string `json:"dependentRequired,omitempty" yaml:"dependentRequired,omitempty" toml:"dependentRequired,omitempty"`

	// Format.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7
	Format string `json:"format,omitempty" yaml:"format,omitempty" toml:"format,omitempty"`

	// Validation of string-encoded data.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-8
	ContEnc    string `json:"contentEncoding,omitempty"  yaml:"contentEncoding,omitempty"  toml:"contentEncoding,omitempty"`
	ContMedia  string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty" toml:"contentMediaType,omitempty"`
	ContSchema string `json:"contentSchema,omitempty"    yaml:"contentSchema,omitempty"    toml:"contentSchema,omitempty"`

	// Metadata annotations.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9
	Title string `json:"title,omitempty"       yaml:"title,omitempty"       toml:"title,omitempty"`
	// Desc     string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Default  any      `json:"default,omitempty"     yaml:"default,omitempty"     toml:"default,omitempty"`
	Depr     bool     `json:"deprecated,omitempty"  yaml:"deprecated,omitempty"  toml:"deprecated,omitempty"`
	Ronly    bool     `json:"readOnly,omitempty"    yaml:"readOnly,omitempty"    toml:"readOnly,omitempty"`
	Wonly    bool     `json:"writeOnly,omitempty"   yaml:"writeOnly,omitempty"   toml:"writeOnly,omitempty"`
	Examples Examples `json:"examples,omitempty"    yaml:"examples,omitempty"    toml:"examples,omitempty"`
}

// Returns `.Title` after validating that it's non-empty.
func (self Schema) ValidTitle() string {
	val := self.Title
	if val == `` {
		panic(fmt.Errorf(`[oas] missing title in schema %#v`, self))
	}
	return val
}

// Value method that returns a pointer. Sometimes useful as a shortcut.
func (self Schema) Opt() *Schema { return &self }

/*
Mostly for internal use. Mutates the receiver to indicate nullability by adding
`oas.TypeNull` to the type. For indicating nullability by wrapping, see
`NullSchema`.
*/
func (self *Schema) Nullable() {
	if self.Ref != `` {
		panic(fmt.Errorf(`[oas] attempted to nullarize schema reference %#v`, *self))
	}
	if !self.IsNullable() {
		self.TypeAdd(TypeNull)
	}
}

/*
True if either `.Type`, `.OneOf`, or `.AnyOf` indicates nullability. Note that
while `true` indicates nullability, `false` does NOT indicate non-nullability,
as the type may reference another, which in turn may be inherently nullable.
*/
func (self Schema) IsNullable() bool {
	return self.TypeHas(TypeNull) ||
		someSchema(self.OneOf, Schema.IsNullable) ||
		someSchema(self.AnyOf, Schema.IsNullable)
}

// Replaces `.Type` with the given vals.
func (self *Schema) TypeReplace(vals ...string) *Schema {
	self.Type = vals
	return self
}

// Adds `vals` to `.Type`, deduplicating them, like a set.
func (self *Schema) TypeAdd(vals ...string) *Schema {
	for _, val := range vals {
		self.typeAdd(val)
	}
	return self
}

// True if the given primitive type is among `.Type`.
func (self Schema) TypeHas(exp string) bool {
	return stringsContain(self.Type, exp)
}

// True if `.Type` exactly matches the given inputs.
func (self Schema) TypeIs(exp ...string) bool {
	return stringsEq(self.Type, exp)
}

// See the doc on the `oas.Schema` type.
type Schemas map[string]Schema

/*
Inits the receiving variable or property to non-nil, returning the resulting
mutable map. Handy for chaining.
*/
func (self *Schemas) Init() Schemas {
	if *self == nil {
		*self = Schemas{}
	}
	return *self
}
