/*
References:

	https://oai.github.io/Documentation/

	https://spec.openapis.org/oas/v3.1.0

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00
*/
package oas

import (
	"fmt"
	"net/http"
	"path"
	r "reflect"
)

const (
	// OpenAPI version supported by this package.
	Ver = `3.1.0`

	TypeNull = `null`
	TypeInt  = `integer`
	TypeNum  = `number`
	TypeStr  = `string`
	TypeBool = `boolean`
	TypeObj  = `object`
	TypeArr  = `array`

	/**
	References:

		https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7.3

		https://swagger.io/docs/specification/data-models/data-types/

	The following are the formats that this package automatically detects.
	*/
	FormatInt32    = `int32`
	FormatInt64    = `int64`
	FormatFloat32  = `float`
	FormatFloat64  = `double`
	FormatDate     = `date`
	FormatTime     = `time`
	FormatDateTime = `date-time`
	FormatDuration = `duration`
	FormatUuid     = `uuid`

	// Well-known formats that this package doesn't automatically detect.
	FormatByte     = `byte`
	FormatBin      = `binary`
	FormatPassword = `password`
	FormatEmail    = `email`

	// Reference: https://spec.openapis.org/oas/v3.1.0#parameter-locations
	InPath   = `path`
	InQuery  = `query`
	InHeader = `header`
	InCookie = `cookie`

	ConTypeJson = `application/json`
)

// Shortcut for creating a ref with `#/components/schema/` for the given name.
func SchemaRef(name string) Ref {
	if name == `` {
		panic(errMissingTitle)
	}
	return RefFrom(path.Join(`#/components/schemas`, name))
}

// Shortcut for creating `oas.Ref` with the given ref, which is stored as-is.
func RefFrom(val string) Ref { return Ref{Ref: val} }

/*
Reference object. References:

	https://spec.openapis.org/oas/v3.1.0#reference-object

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.3.1
*/
type Ref struct {
	Ref  string `json:"$ref,omitempty"`
	Sum  string `json:"sum,omitempty"`
	Desc string `json:"description,omitempty"`
}

// True if the reference isn't empty.
func (self Ref) HasRef() bool { return self.Ref != `` }

/*
Top-level OpenAPI document. Reference:

	https://spec.openapis.org/oas/v3.1.0#openapi-object
*/
type Doc struct {
	Ref
	Openapi    string   `json:"openapi,omitempty"`
	Info       *Info    `json:"info,omitempty"`
	JsonSchema string   `json:"jsonSchemaDialect,omitempty"`
	Servers    []Server `json:"servers,omitempty"`
	Paths      Paths    `json:"paths,omitempty"`
	Webhooks   Paths    `json:"webhooks,omitempty"`
	Comps      Comps    `json:"components,omitempty"`
	Security   []SecReq `json:"security,omitempty"`
	Tags       []Tag    `json:"tags,omitempty"`
	ExtDoc     *ExtDoc  `json:"externalDocs,omitempty"`
}

// Shortcut for registering a route via `oas.Doc.Paths.Route`.
func (self *Doc) Route(path, meth string, op Op) *Doc {
	self.Paths.Init().Route(path, meth, op)
	return self
}

/*
Looks up a schema by the given name among the doc's components. The name must be
the exact schema title, not a reference path. May panic if the schema
unexpectedly has double indirection.
*/
func (self *Doc) GotCompSchema(name string) (Schema, bool) {
	val, ok := self.Comps.Schemas[name]
	if val.HasRef() {
		panic(errSchemaUnexpectedRef(name, val.Ref.Ref))
	}
	return val, ok
}

/*
Looks up the schema by a full reference path. Currently supports only component
references starting with `#/components/schemas/`.
*/
func (self *Doc) GotSchema(refPath string) (Schema, bool) {
	name, ok := unprefix(refPath, `#/components/schemas/`)
	if ok {
		return self.GotCompSchema(name)
	}
	panic(fmt.Errorf(`[oas] unsupported schema reference %q`, refPath))
}

/*
Dereferences the given schema, returning a non-reference. The bool is true if
the target schema was found, false otherwise. May panic if the schema
unexpectedly has double indirection.
*/
func (self *Doc) DerefSchema(sch Schema) (Schema, bool) {
	if sch.HasRef() {
		return self.GotSchema(sch.Ref.Ref)
	}
	return sch, true
}

// Same as `.JsonBody(typ).Opt()` but slightly clearer.
func (self *Doc) JsonBodyOpt(typ interface{}) *Body {
	return self.JsonBody(typ).Opt()
}

/*
Shortcut. Returns `oas.Body` describing a JSON response with the schema of the
given type, after registering its schema in the document. The input is used
only as a type carrier; its actual value is ignored.
*/
func (self *Doc) JsonBody(typ interface{}) Body {
	return Body{Cont: MediaTypes{ConTypeJson: self.SchemaMedia(typ)}}
}

/*
Shortcut. Returns `oas.Resps` with 200 JSON for the given type, after
registering its schema in the document. The input is used only as a type
carrier; its actual value is ignored.
*/
func (self *Doc) RespsOkJson(typ interface{}) Resps {
	return Resps{
		`200`: Resp{
			Cont: MediaTypes{
				ConTypeJson: {Schema: self.Sch(typ)},
			},
		},
	}
}

/*
Shortcut. Returns `oas.MediaType` with the schema of the given type, after
registering its schema in the document. The input is used only as a type
carrier; its actual value is ignored.
*/
func (self *Doc) SchemaMedia(typ interface{}) MediaType {
	return MediaType{Schema: self.Sch(typ)}
}

/*
Shortcut for returning `.TypeSchema` from the input's type. The input value is
used only as a type carrier.
*/
func (self *Doc) Sch(typ interface{}) Schema {
	return self.TypeSchema(r.TypeOf(typ))
}

/*
Returns an OAS schema for the given Go type. May register various associated
types in `.Comps.Schemas`, mutating the document. The returned schema may be a
reference.
*/
func (self *Doc) TypeSchema(typ r.Type) (sch Schema) {
	self.schemaAny(&sch, typ)
	return
}

// https://spec.openapis.org/oas/v3.1.0#info-object
type Info struct {
	Ref
	Title   string   `json:"title,omitempty"`
	Sum     string   `json:"sum,omitempty"`
	Desc    string   `json:"description,omitempty"`
	Terms   string   `json:"termsOfService,omitempty"`
	Contact *Contact `json:"contact,omitempty"`
	License *License `json:"license,omitempty"`
	Ver     string   `json:"version,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#contact-object
type Contact struct {
	Ref
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#license-object
type License struct {
	Ref
	Name  string `json:"name,omitempty"`
	Ident string `json:"identifier,omitempty"`
	Url   string `json:"url,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#server-object
type Server struct {
	Ref
	Url  string `json:"url,omitempty"`
	Desc string `json:"description,omitempty"`
	Vars Vars   `json:"variables,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#server-variable-object
type Vars map[string]Var

// https://spec.openapis.org/oas/v3.1.0#server-variable-object
type Var struct {
	Ref
	Enum    []string `json:"enum,omitempty"`
	Default string   `json:"default,omitempty"`
	Desc    string   `json:"description,omitempty"`
}

// Short for "components":
// https://spec.openapis.org/oas/v3.1.0#components-object
type Comps struct {
	Ref
	Schemas    Schemas    `json:"schemas,omitempty"`
	Resps      Resps      `json:"responses,omitempty"`
	Params     Params     `json:"parameters,omitempty"`
	Examples   Examples   `json:"examples,omitempty"`
	Reqs       Bodies     `json:"requestBodies,omitempty"`
	Heads      Heads      `json:"headers,omitempty"`
	SecSchemes SecSchemes `json:"securitySchemes,omitempty"`
	Links      Links      `json:"links,omitempty"`
	Callbacks  Callbacks  `json:"callbacks,omitempty"`
	Paths      Paths      `json:"pathItems,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#paths-object
type Paths map[string]Path

/*
Inits the receiving variable or property to non-nil, returning the resulting
mutable map. Handy for chaining.
*/
func (self *Paths) Init() Paths {
	if *self == nil {
		*self = Paths{}
	}
	return *self
}

/*
Shortcut for registering an "op" at the given path and method, via
`(*oas.Path).Method`.
*/
func (self Paths) Route(path, meth string, op Op) Paths {
	val := self[path]
	val.Method(meth, op)
	self[path] = val
	return self
}

// Called "path item" in the spec:
// https://spec.openapis.org/oas/v3.1.0#path-item-object
type Path struct {
	Ref
	Sum     string   `json:"sum,omitempty"`
	Desc    string   `json:"description,omitempty"`
	Get     *Op      `json:"get,omitempty"`
	Put     *Op      `json:"put,omitempty"`
	Post    *Op      `json:"post,omitempty"`
	Delete  *Op      `json:"delete,omitempty"`
	Options *Op      `json:"options,omitempty"`
	Head    *Op      `json:"head,omitempty"`
	Patch   *Op      `json:"patch,omitempty"`
	Trace   *Op      `json:"trace,omitempty"`
	Servers []Server `json:"servers,omitempty"`
	Params  []Param  `json:"parameters,omitempty"`
}

/*
Sets the "op" at the given method. The method must be well-known, otherwise this
will panic.
*/
func (self *Path) Method(meth string, op Op) *Path {
	switch meth {
	case http.MethodGet:
		self.Get = &op
	case http.MethodPut:
		self.Put = &op
	case http.MethodPost:
		self.Post = &op
	case http.MethodDelete:
		self.Delete = &op
	case http.MethodOptions:
		self.Options = &op
	case http.MethodHead:
		self.Head = &op
	case http.MethodPatch:
		self.Patch = &op
	case http.MethodTrace:
		self.Trace = &op
	default:
		panic(fmt.Errorf(`[oas] unrecognized method %q`, meth))
	}
	return self
}

// Short for "operation":
// https://spec.openapis.org/oas/v3.1.0#operation-object
type Op struct {
	Ref
	Tags      []Tag     `json:"tags,omitempty"`
	Sum       string    `json:"summary,omitempty"`
	Desc      string    `json:"description,omitempty"`
	ExtDoc    *ExtDoc   `json:"externalDocs,omitempty"`
	OpId      string    `json:"operationId,omitempty"`
	Params    []Param   `json:"parameters,omitempty"`
	ReqBody   *Body     `json:"requestBody,omitempty"`
	Resps     Resps     `json:"responses,omitempty"`
	Callbacks Callbacks `json:"callbacks,omitempty"`
	Depr      bool      `json:"deprecated,omitempty"`
	Sec       []SecReq  `json:"security,omitempty"`
	Servers   []Server  `json:"servers,omitempty"`
}

// Short for "external documentation":
// https://spec.openapis.org/oas/v3.1.0#external-documentation-object
type ExtDoc struct {
	Ref
	Desc string `json:"description,omitempty"`
	Url  string `json:"url,omitempty"`
}

// Short for "parameter":
// https://spec.openapis.org/oas/v3.1.0#parameter-object
type Param struct {
	Head
	Name string `json:"name,omitempty"`
	In   string `json:"in,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#parameter-object
type Params map[string]Param

// Short for "request body":
// https://spec.openapis.org/oas/v3.1.0#request-body-object
type Body struct {
	Ref
	Desc string     `json:"desc,omitempty"`
	Cont MediaTypes `json:"content,omitempty"`
	Requ bool       `json:"required,omitempty"`
}

// Value method that returns a pointer. Sometimes useful as a shortcut.
func (self Body) Opt() *Body { return &self }

// https://spec.openapis.org/oas/v3.1.0#request-body-object
type Bodies map[string]Body

// https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaType struct {
	Ref
	Schema   Schema    `json:"schema,omitempty"`
	Example  Any       `json:"example,omitempty"`
	Examples Examples  `json:"examples,omitempty"`
	Encoding Encodings `json:"encoding,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaTypes map[string]MediaType

// https://spec.openapis.org/oas/v3.1.0#encoding-object
type Encoding struct {
	Ref
	ConType  string `json:"contentType,omitempty"`
	Head     Heads  `json:"headers,omitempty"`
	Style    string `json:"style,omitempty"`
	Explode  bool   `json:"explode,omitempty"`
	Reserved bool   `json:"allowReserved,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#encoding-object
type Encodings map[string]Encoding

// Short for "response":
// https://spec.openapis.org/oas/v3.1.0#response-object
type Resp struct {
	Ref
	Desc  string     `json:"description,omitempty"`
	Head  Heads      `json:"headers,omitempty"`
	Cont  MediaTypes `json:"content,omitempty"`
	Links Links      `json:"links,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#responses-object
type Resps map[string]Resp

// https://spec.openapis.org/oas/v3.1.0#callback-object.
// May also be `{"$ref": "..."}`.
type Callback map[string]string

// https://spec.openapis.org/oas/v3.1.0#callback-object
type Callbacks map[string]Callback

// https://spec.openapis.org/oas/v3.1.0#example-object
type Example struct {
	Ref
	Sum   string `json:"summary,omitempty"`
	Desc  string `json:"description,omitempty"`
	Val   string `json:"value,omitempty"`
	ExVal string `json:"externalValue,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#example-object
type Examples map[string]Example

// https://spec.openapis.org/oas/v3.1.0#link-object
type Link struct {
	Ref
	OpRef   string  `json:"operationRef,omitempty"`
	OpId    string  `json:"operationId,omitempty"`
	Params  Anys    `json:"parameters,omitempty"`
	ReqBody Any     `json:"requestBody,omitempty"`
	Desc    string  `json:"description,omitempty"`
	Server  *Server `json:"server,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#link-object
type Links map[string]Link

// Short for "header":
// https://spec.openapis.org/oas/v3.1.0#header-object
type Head struct {
	Ref
	Desc     string     `json:"desc,omitempty"`
	Requ     bool       `json:"required,omitempty"`
	Depr     bool       `json:"deprecated,omitempty"`
	Empty    bool       `json:"allowEmptyValue,omitempty"`
	Style    string     `json:"style,omitempty"`
	Explode  bool       `json:"explode,omitempty"`
	Reserved bool       `json:"allowReserved,omitempty"`
	Schema   *Schema    `json:"schema,omitempty"`
	Example  Any        `json:"example,omitempty"`
	Examples Examples   `json:"examples,omitempty"`
	Cont     MediaTypes `json:"content,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#header-object
type Heads map[string]Head

// https://spec.openapis.org/oas/v3.1.0#tag-object
type Tag struct {
	Ref
	Name   string  `json:"name,omitempty"`
	Desc   string  `json:"desc,omitempty"`
	ExtDoc *ExtDoc `json:"externalDocs,omitempty"`
}

/*
Shortcut for making a reference-only schema pointing at
`#/components/schema/<name>`.
*/
func RefSchema(name string) Schema { return Schema{Ref: SchemaRef(name)} }

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
	Ref

	/**
	OAS properties.
	*/

	Discr   *Discr  `json:"discriminator,omitempty"`
	Xml     *Xml    `json:"xml,omitempty"`
	ExtDoc  *ExtDoc `json:"externalDocs,omitempty"`
	Example Any     `json:"example,omitempty"` // Also see `.Examples`.

	/**
	JSON Schema core properties. Reference:

		https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00
	*/

	// Subschemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.1
	AllOf []Schema `json:"allOf,omitempty"`
	AnyOf []Schema `json:"anyOf,omitempty"`
	OneOf []Schema `json:"oneOf,omitempty"`
	Not   *Schema  `json:"not,omitempty"`

	// Conditional subschemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.2.2
	If         *Schema `json:"if,omitempty"`
	Then       *Schema `json:"then,omitempty"`
	Else       *Schema `json:"else,omitempty"`
	DepSchemas Schemas `json:"dependentSchemas,omitempty"`

	// Array child schemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.1
	PrefixItems []Schema `json:"prefixItems,omitempty"`
	Items       *Schema  `json:"items,omitempty"`
	Contains    *Schema  `json:"contains,omitempty"`

	// Object subschemas.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-10.3.2
	Props     Schemas `json:"properties,omitempty"`
	PatProps  Schemas `json:"patternProperties,omitempty"`
	AddProps  *Schema `json:"additionalProperties,omitempty"`
	PropNames *Schema `json:"propertyNames,omitempty"`

	// Unevaluated locations.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-11
	UnevalItems *Schema `json:"unevaluatedItems,omitempty"`
	UnevalProps *Schema `json:"unevaluatedProperties,omitempty"`

	/**
	JSON Schema validation properties. Reference:

		https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00
	*/

	// Validation for any instance.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.1
	Type  []string `json:"type,omitempty"`
	Enum  []string `json:"enum,omitempty"`
	Const Any      `json:"const,omitempty"`

	// Validation for numeric instances.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.2
	MulOf   uint64 `json:"multipleOf,omitempty"` // 0 represents "missing".
	Max     *int64 `json:"maximum,omitempty"`
	ExlcMax *int64 `json:"exclusiveMaximum,omitempty"`
	Min     *int64 `json:"minimum,omitempty"`
	ExclMin *int64 `json:"exclusiveMinimum,omitempty"`

	// Validation for strings.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.3
	MaxLen  uint64 `json:"maxLength,omitempty"` // 0 represents "missing".
	MinLen  uint64 `json:"minLength,omitempty"`
	Pattern string `json:"pattern,omitempty"`

	// Validation for arrays.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.4
	MaxItems  uint64 `json:"maxItems,omitempty"`
	MinItems  uint64 `json:"minItems,omitempty"`
	UniqItems bool   `json:"uniqueItems,omitempty"`
	MaxCont   uint64 `json:"maxContains,omitempty"`
	MinCont   uint64 `json:"minContains,omitempty"`

	// Validation for objects.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-6.5
	MaxProps uint64              `json:"maxProperties,omitempty"`
	MinProps uint64              `json:"minProperties,omitempty"`
	Requ     bool                `json:"required,omitempty"`
	DepRequ  map[string][]string `json:"dependentRequired,omitempty"`

	// Format.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-7
	Format string `json:"format,omitempty"`

	// Validation of string-encoded data.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-8
	ContEnc    string `json:"contentEncoding,omitempty"`
	ContMedia  string `json:"contentMediaType,omitempty"`
	ContSchema string `json:"contentSchema,omitempty"`

	// Metadata annotations.
	// https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00#section-9
	Title    string   `json:"title,omitempty"`
	Desc     string   `json:"description,omitempty"`
	Default  Any      `json:"default,omitempty"`
	Depr     bool     `json:"deprecated,omitempty"`
	Ronly    bool     `json:"readOnly,omitempty"`
	Wonly    bool     `json:"writeOnly,omitempty"`
	Examples Examples `json:"examples,omitempty"`
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
	if self.HasRef() {
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

// Short for "discriminator":
// https://spec.openapis.org/oas/v3.1.0#discriminator-object
type Discr struct {
	Ref
	Prop string            `json:"propertyName,omitempty"`
	Map  map[string]string `json:"mapping,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#xml-object
type Xml struct {
	Ref
	Name   string `json:"name,omitempty"`
	Nspace string `json:"namespace,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Attr   bool   `json:"attribute,omitempty"`
	Wrap   bool   `json:"wrapped,omitempty"`
}

// Short for "security scheme".
// https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecScheme struct {
	Ref
	Type       string `json:"type,omitempty"`
	Desc       string `json:"description,omitempty"`
	Name       string `json:"name,omitempty"`
	In         string `json:"in,omitempty"`
	Scheme     string `json:"scheme,omitempty"`
	BearFormat string `json:"bearerFormat,omitempty"`
	Flows      *Flows `json:"flows,omitempty"`
	OidUrl     string `json:"openIdConnectUrl,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecSchemes map[string]SecScheme

// https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
type Flows struct {
	Ref
	Implicit   Flow `json:"implicit,omitempty"`
	Password   Flow `json:"password,omitempty"`
	ClientCred Flow `json:"clientCredentials,omitempty"`
	AuthCode   Flow `json:"authorizationCode,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#oauth-flow-object
type Flow struct {
	Ref
	AuthUrl    string            `json:"authorizationUrl,omitempty"`
	TokenUrl   string            `json:"tokenUrl,omitempty"`
	RefreshUrl string            `json:"refreshUrl,omitempty"`
	Scopes     map[string]string `json:"scopes,omitempty"`
}

// Short for "secutity requirement".
type SecReq map[string][]string

type Any interface{}

type Anys map[string]Any
