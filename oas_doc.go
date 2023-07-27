package oas

import (
	"fmt"
	r "reflect"
)

/*
Top-level OpenAPI document. Reference:

	https://spec.openapis.org/oas/v3.1.0#openapi-object
*/
type Doc struct {
	Ref        string   `json:"$ref,omitempty"        yaml:"$ref,omitempty"        toml:"$ref,omitempty"`
	Sum        string   `json:"sum,omitempty"         yaml:"sum,omitempty"         toml:"sum,omitempty"`
	Desc       string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Openapi    string   `json:"openapi,omitempty"           yaml:"openapi,omitempty"           toml:"openapi,omitempty"`
	Info       *Info    `json:"info,omitempty"              yaml:"info,omitempty"              toml:"info,omitempty"`
	JsonSchema string   `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty" toml:"jsonSchemaDialect,omitempty"`
	Servers    []Server `json:"servers,omitempty"           yaml:"servers,omitempty"           toml:"servers,omitempty"`
	Paths      Paths    `json:"paths,omitempty"             yaml:"paths,omitempty"             toml:"paths,omitempty"`
	Webhooks   Paths    `json:"webhooks,omitempty"          yaml:"webhooks,omitempty"          toml:"webhooks,omitempty"`
	Comps      Comps    `json:"components,omitempty"        yaml:"components,omitempty"        toml:"components,omitempty"`
	Security   []SecReq `json:"security,omitempty"          yaml:"security,omitempty"          toml:"security,omitempty"`
	Tags       []Tag    `json:"tags,omitempty"              yaml:"tags,omitempty"              toml:"tags,omitempty"`
	ExtDoc     *ExtDoc  `json:"externalDocs,omitempty"      yaml:"externalDocs,omitempty"      toml:"externalDocs,omitempty"`
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

/*
Shortcut for returning `.TypeSchema` from the input's type. The input value is
used only as a type carrier.
*/
func (self *Doc) Sch(typ interface{}) Schema {
	return self.TypeSchema(r.TypeOf(typ))
}

/*
Shortcut. Returns `oas.MediaType` with the schema of the given type, after
registering its schema in the document. The input is used only as a type
carrier; its actual value is ignored.
*/
func (self *Doc) SchemaMedia(typ interface{}) MediaType {
	return MediaType{Schema: self.Sch(typ)}
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
	if val.Ref != `` {
		panic(errSchemaUnexpectedRef(name, val.Ref))
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
	if sch.Ref != `` {
		return self.GotSchema(sch.Ref)
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
