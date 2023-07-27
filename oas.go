/*
References:

	https://oai.github.io/Documentation/

	https://spec.openapis.org/oas/v3.1.0

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-00

	https://datatracker.ietf.org/doc/html/draft-bhutton-json-schema-validation-00
*/
package oas

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
