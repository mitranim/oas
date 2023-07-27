package oas

import (
	"fmt"
	"net/http"
)

// Represents maps of "any type" in some OAS definitions.
type Anys map[string]any

/*
Reference object. References:

	https://spec.openapis.org/oas/v3.1.0#reference-object

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00#section-8.2.3.1

Most of our types have a `Ref string` field, and technically they should embed
this type to avoid unnecessary duplication. However, we copy the fields instead
of embedding the type to ensure better compatibility with 3rd party encoders,
some of which don't seem to support embedded structs, particularly for YAML.
*/
type Ref struct {
	Ref  string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum  string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#info-object
type Info struct {
	Ref     string   `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum     string   `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc    string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Title   string   `json:"title,omitempty"          yaml:"title,omitempty"          toml:"title,omitempty"`
	Terms   string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty" toml:"termsOfService,omitempty"`
	Contact *Contact `json:"contact,omitempty"        yaml:"contact,omitempty"        toml:"contact,omitempty"`
	License *License `json:"license,omitempty"        yaml:"license,omitempty"        toml:"license,omitempty"`
	Ver     string   `json:"version,omitempty"        yaml:"version,omitempty"        toml:"version,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#contact-object
type Contact struct {
	Ref   string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum   string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc  string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Name  string `json:"name,omitempty"  yaml:"name,omitempty"  toml:"name,omitempty"`
	Url   string `json:"url,omitempty"   yaml:"url,omitempty"   toml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty" toml:"email,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#license-object
type License struct {
	Ref   string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum   string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc  string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Name  string `json:"name,omitempty"       yaml:"name,omitempty"       toml:"name,omitempty"`
	Ident string `json:"identifier,omitempty" yaml:"identifier,omitempty" toml:"identifier,omitempty"`
	Url   string `json:"url,omitempty"        yaml:"url,omitempty"        toml:"url,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#server-object
type Server struct {
	Ref  string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum  string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Url  string `json:"url,omitempty"         yaml:"url,omitempty"         toml:"url,omitempty"`
	Vars Vars   `json:"variables,omitempty"   yaml:"variables,omitempty"   toml:"variables,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#server-variable-object
type Vars map[string]Var

// https://spec.openapis.org/oas/v3.1.0#server-variable-object
type Var struct {
	Ref     string   `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum     string   `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc    string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Enum    []string `json:"enum,omitempty"        yaml:"enum,omitempty"        toml:"enum,omitempty"`
	Default string   `json:"default,omitempty"     yaml:"default,omitempty"     toml:"default,omitempty"`
}

// Short for "components":
// https://spec.openapis.org/oas/v3.1.0#components-object
type Comps struct {
	Ref        string     `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum        string     `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc       string     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Schemas    Schemas    `json:"schemas,omitempty"         yaml:"schemas,omitempty"         toml:"schemas,omitempty"`
	Resps      Resps      `json:"responses,omitempty"       yaml:"responses,omitempty"       toml:"responses,omitempty"`
	Params     Params     `json:"parameters,omitempty"      yaml:"parameters,omitempty"      toml:"parameters,omitempty"`
	Examples   Examples   `json:"examples,omitempty"        yaml:"examples,omitempty"        toml:"examples,omitempty"`
	Reqs       Bodies     `json:"requestBodies,omitempty"   yaml:"requestBodies,omitempty"   toml:"requestBodies,omitempty"`
	Heads      Heads      `json:"headers,omitempty"         yaml:"headers,omitempty"         toml:"headers,omitempty"`
	SecSchemes SecSchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty" toml:"securitySchemes,omitempty"`
	Links      Links      `json:"links,omitempty"           yaml:"links,omitempty"           toml:"links,omitempty"`
	Callbacks  Callbacks  `json:"callbacks,omitempty"       yaml:"callbacks,omitempty"       toml:"callbacks,omitempty"`
	Paths      Paths      `json:"pathItems,omitempty"       yaml:"pathItems,omitempty"       toml:"pathItems,omitempty"`
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
	/**
	Tentative. This is useful for many UI visualizers, which would otherwise try
	to generate a summary from the description, which is annoying in practice.
	May revise.
	*/
	if op.Sum == `` {
		op.Sum = path
	}

	val := self[path]
	val.Method(meth, op)
	self[path] = val
	return self
}

// Called "path item" in the spec:
// https://spec.openapis.org/oas/v3.1.0#path-item-object
type Path struct {
	Ref     string   `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum     string   `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc    string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Get     *Op      `json:"get,omitempty"         yaml:"get,omitempty"         toml:"get,omitempty"`
	Put     *Op      `json:"put,omitempty"         yaml:"put,omitempty"         toml:"put,omitempty"`
	Post    *Op      `json:"post,omitempty"        yaml:"post,omitempty"        toml:"post,omitempty"`
	Delete  *Op      `json:"delete,omitempty"      yaml:"delete,omitempty"      toml:"delete,omitempty"`
	Options *Op      `json:"options,omitempty"     yaml:"options,omitempty"     toml:"options,omitempty"`
	Head    *Op      `json:"head,omitempty"        yaml:"head,omitempty"        toml:"head,omitempty"`
	Patch   *Op      `json:"patch,omitempty"       yaml:"patch,omitempty"       toml:"patch,omitempty"`
	Trace   *Op      `json:"trace,omitempty"       yaml:"trace,omitempty"       toml:"trace,omitempty"`
	Servers []Server `json:"servers,omitempty"     yaml:"servers,omitempty"     toml:"servers,omitempty"`
	Params  []Param  `json:"parameters,omitempty"  yaml:"parameters,omitempty"  toml:"parameters,omitempty"`
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
	Ref       string    `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum       string    `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc      string    `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Tags      []Tag     `json:"tags,omitempty"         yaml:"tags,omitempty"         toml:"tags,omitempty"`
	ExtDoc    *ExtDoc   `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty" toml:"externalDocs,omitempty"`
	OpId      string    `json:"operationId,omitempty"  yaml:"operationId,omitempty"  toml:"operationId,omitempty"`
	Params    []Param   `json:"parameters,omitempty"   yaml:"parameters,omitempty"   toml:"parameters,omitempty"`
	ReqBody   *Body     `json:"requestBody,omitempty"  yaml:"requestBody,omitempty"  toml:"requestBody,omitempty"`
	Resps     Resps     `json:"responses,omitempty"    yaml:"responses,omitempty"    toml:"responses,omitempty"`
	Callbacks Callbacks `json:"callbacks,omitempty"    yaml:"callbacks,omitempty"    toml:"callbacks,omitempty"`
	Depr      bool      `json:"deprecated,omitempty"   yaml:"deprecated,omitempty"   toml:"deprecated,omitempty"`
	Sec       []SecReq  `json:"security,omitempty"     yaml:"security,omitempty"     toml:"security,omitempty"`
	Servers   []Server  `json:"servers,omitempty"      yaml:"servers,omitempty"      toml:"servers,omitempty"`
}

// Short for "external documentation":
// https://spec.openapis.org/oas/v3.1.0#external-documentation-object
type ExtDoc struct {
	Ref  string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum  string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Url  string `json:"url,omitempty"         yaml:"url,omitempty"         toml:"url,omitempty"`
}

// Short for "parameter":
// https://spec.openapis.org/oas/v3.1.0#parameter-object
type Param struct {
	Head
	Name string `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	In   string `json:"in,omitempty"   yaml:"in,omitempty"   toml:"in,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#parameter-object
type Params map[string]Param

// Short for "request body":
// https://spec.openapis.org/oas/v3.1.0#request-body-object
type Body struct {
	Ref  string     `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum  string     `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc string     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Cont MediaTypes `json:"content,omitempty"  yaml:"content,omitempty"  toml:"content,omitempty"`
	Requ bool       `json:"required,omitempty" yaml:"required,omitempty" toml:"required,omitempty"`
}

// Value method that returns a pointer. Sometimes useful as a shortcut.
func (self Body) Opt() *Body { return &self }

// https://spec.openapis.org/oas/v3.1.0#request-body-object
type Bodies map[string]Body

// https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaType struct {
	Ref      string    `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum      string    `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc     string    `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Schema   Schema    `json:"schema,omitempty"   yaml:"schema,omitempty"   toml:"schema,omitempty"`
	Example  any       `json:"example,omitempty"  yaml:"example,omitempty"  toml:"example,omitempty"`
	Examples Examples  `json:"examples,omitempty" yaml:"examples,omitempty" toml:"examples,omitempty"`
	Encoding Encodings `json:"encoding,omitempty" yaml:"encoding,omitempty" toml:"encoding,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaTypes map[string]MediaType

// https://spec.openapis.org/oas/v3.1.0#encoding-object
type Encoding struct {
	Ref      string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum      string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc     string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	ConType  string `json:"contentType,omitempty"   yaml:"contentType,omitempty"   toml:"contentType,omitempty"`
	Head     Heads  `json:"headers,omitempty"       yaml:"headers,omitempty"       toml:"headers,omitempty"`
	Style    string `json:"style,omitempty"         yaml:"style,omitempty"         toml:"style,omitempty"`
	Explode  bool   `json:"explode,omitempty"       yaml:"explode,omitempty"       toml:"explode,omitempty"`
	Reserved bool   `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty" toml:"allowReserved,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#encoding-object
type Encodings map[string]Encoding

// Short for "response":
// https://spec.openapis.org/oas/v3.1.0#response-object
type Resp struct {
	Ref   string     `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum   string     `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc  string     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Head  Heads      `json:"headers,omitempty"     yaml:"headers,omitempty"     toml:"headers,omitempty"`
	Cont  MediaTypes `json:"content,omitempty"     yaml:"content,omitempty"     toml:"content,omitempty"`
	Links Links      `json:"links,omitempty"       yaml:"links,omitempty"       toml:"links,omitempty"`
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
	Ref   string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum   string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc  string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Val   string `json:"value,omitempty"         yaml:"value,omitempty"         toml:"value,omitempty"`
	ExVal string `json:"externalValue,omitempty" yaml:"externalValue,omitempty" toml:"externalValue,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#example-object
type Examples map[string]Example

// https://spec.openapis.org/oas/v3.1.0#link-object
type Link struct {
	Ref     string  `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum     string  `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc    string  `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	OpRef   string  `json:"operationRef,omitempty" yaml:"operationRef,omitempty" toml:"operationRef,omitempty"`
	OpId    string  `json:"operationId,omitempty"  yaml:"operationId,omitempty"  toml:"operationId,omitempty"`
	Params  Anys    `json:"parameters,omitempty"   yaml:"parameters,omitempty"   toml:"parameters,omitempty"`
	ReqBody any     `json:"requestBody,omitempty"  yaml:"requestBody,omitempty"  toml:"requestBody,omitempty"`
	Server  *Server `json:"server,omitempty"       yaml:"server,omitempty"       toml:"server,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#link-object
type Links map[string]Link

// Short for "header":
// https://spec.openapis.org/oas/v3.1.0#header-object
type Head struct {
	Ref      string     `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum      string     `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc     string     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Requ     bool       `json:"required,omitempty"        yaml:"required,omitempty"        toml:"required,omitempty"`
	Depr     bool       `json:"deprecated,omitempty"      yaml:"deprecated,omitempty"      toml:"deprecated,omitempty"`
	Empty    bool       `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty" toml:"allowEmptyValue,omitempty"`
	Style    string     `json:"style,omitempty"           yaml:"style,omitempty"           toml:"style,omitempty"`
	Explode  bool       `json:"explode,omitempty"         yaml:"explode,omitempty"         toml:"explode,omitempty"`
	Reserved bool       `json:"allowReserved,omitempty"   yaml:"allowReserved,omitempty"   toml:"allowReserved,omitempty"`
	Schema   *Schema    `json:"schema,omitempty"          yaml:"schema,omitempty"          toml:"schema,omitempty"`
	Example  any        `json:"example,omitempty"         yaml:"example,omitempty"         toml:"example,omitempty"`
	Examples Examples   `json:"examples,omitempty"        yaml:"examples,omitempty"        toml:"examples,omitempty"`
	Cont     MediaTypes `json:"content,omitempty"         yaml:"content,omitempty"         toml:"content,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#header-object
type Heads map[string]Head

// https://spec.openapis.org/oas/v3.1.0#tag-object
type Tag struct {
	Ref    string  `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum    string  `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc   string  `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Name   string  `json:"name,omitempty"         yaml:"name,omitempty"         toml:"name,omitempty"`
	ExtDoc *ExtDoc `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty" toml:"externalDocs,omitempty"`
}

// Short for "discriminator":
// https://spec.openapis.org/oas/v3.1.0#discriminator-object
type Discr struct {
	Ref  string            `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum  string            `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc string            `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Prop string            `json:"propertyName,omitempty" yaml:"propertyName,omitempty" toml:"propertyName,omitempty"`
	Map  map[string]string `json:"mapping,omitempty"      yaml:"mapping,omitempty"      toml:"mapping,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#xml-object
type Xml struct {
	Ref    string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum    string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc   string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Name   string `json:"name,omitempty"      yaml:"name,omitempty"      toml:"name,omitempty"`
	Nspace string `json:"namespace,omitempty" yaml:"namespace,omitempty" toml:"namespace,omitempty"`
	Prefix string `json:"prefix,omitempty"    yaml:"prefix,omitempty"    toml:"prefix,omitempty"`
	Attr   bool   `json:"attribute,omitempty" yaml:"attribute,omitempty" toml:"attribute,omitempty"`
	Wrap   bool   `json:"wrapped,omitempty"   yaml:"wrapped,omitempty"   toml:"wrapped,omitempty"`
}

// Short for "security scheme".
// https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecScheme struct {
	Ref        string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum        string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc       string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Type       string `json:"type,omitempty"             yaml:"type,omitempty"             toml:"type,omitempty"`
	Name       string `json:"name,omitempty"             yaml:"name,omitempty"             toml:"name,omitempty"`
	In         string `json:"in,omitempty"               yaml:"in,omitempty"               toml:"in,omitempty"`
	Scheme     string `json:"scheme,omitempty"           yaml:"scheme,omitempty"           toml:"scheme,omitempty"`
	BearFormat string `json:"bearerFormat,omitempty"     yaml:"bearerFormat,omitempty"     toml:"bearerFormat,omitempty"`
	Flows      *Flows `json:"flows,omitempty"            yaml:"flows,omitempty"            toml:"flows,omitempty"`
	OidUrl     string `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty" toml:"openIdConnectUrl,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecSchemes map[string]SecScheme

// https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
type Flows struct {
	Ref        string `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum        string `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc       string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Implicit   Flow   `json:"implicit,omitempty"          yaml:"implicit,omitempty"          toml:"implicit,omitempty"`
	Password   Flow   `json:"password,omitempty"          yaml:"password,omitempty"          toml:"password,omitempty"`
	ClientCred Flow   `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty" toml:"clientCredentials,omitempty"`
	AuthCode   Flow   `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty" toml:"authorizationCode,omitempty"`
}

// https://spec.openapis.org/oas/v3.1.0#oauth-flow-object
type Flow struct {
	Ref        string            `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum        string            `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc       string            `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	AuthUrl    string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"  toml:"authorizationUrl,omitempty"`
	TokenUrl   string            `json:"tokenUrl,omitempty"         yaml:"tokenUrl,omitempty"          toml:"tokenUrl,omitempty"`
	RefreshUrl string            `json:"refreshUrl,omitempty"       yaml:"refreshUrl,omitempty"        toml:"refreshUrl,omitempty"`
	Scopes     map[string]string `json:"scopes,omitempty"           yaml:"scopes,omitempty"            toml:"scopes,omitempty"`
}

// Short for "secutity requirement".
type SecReq map[string][]string
