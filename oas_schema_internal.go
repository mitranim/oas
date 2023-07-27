package oas

import "fmt"

func (self *Schema) setRef(val string) {
	if val == `` {
		panic(errMissingTitle)
	}

	if self.Ref != `` {
		panic(fmt.Errorf(
			`[oas] attempted to componentize a schema that is already a reference: %#v`,
			*self,
		))
	}

	*self = RefSchema(val)
}

func (self *Schema) typeAdd(val string) {
	if self.TypeHas(val) {
		return
	}

	types := self.Type

	/**
	Motive: keep `TypeNull` at the end of the list, for better documentation
	display. Visualizers may use the first type in the list when rendering
	examples.
	*/
	if len(types) > 0 && types[len(types)-1] == TypeNull {
		self.Type = append(types[:len(types)-1], val, TypeNull)
	} else {
		self.Type = append(types, val)
	}
}
