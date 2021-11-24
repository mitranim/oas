package oas

import (
	"encoding/json"
	"io"
	"testing"
	u "unsafe"
)

/*
At the time of writing, the package is stupidly unoptimized. Generating the
schema for `Doc` may take hundreds of microseconds. For very bloated codebases,
the total time to generate their schemas could reach into milliseconds. We may
look into optimizing this later. However, this is a one-off cost paid at app
initialization time, and it's nothing compared to the slowness of all other
generators.

The cost of JSON-encoding a large doc seems comparable to the cost of its
generation (slightly higher). We should probably provide a shortcut for a
request handler that JSON-encodes a doc once, lazily, and serves the
pre-encoded response for each subsequent request.
*/

func Benchmark_doc_schema_big(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var doc Doc
		doc.Sch(&doc)
	}
}

func Benchmark_doc_small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tDoc()
	}
}

func Benchmark_doc_schema_repeat(b *testing.B) {
	var doc Doc
	for i := 0; i < b.N; i++ {
		doc.Sch((*Outer)(nil))
	}
}

func Benchmark_doc_big_json_encode(b *testing.B) {
	var doc Doc
	doc.Sch(&doc)
	enc := json.NewEncoder(io.Discard)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		try(enc.Encode(&doc))
	}
}

func Benchmark_memcpy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var tar [4]byte
		var src [4]byte = [4]byte{10, 20, 30, 40}

		memcpy(
			uintptr(u.Pointer(&tar)),
			uintptr(u.Pointer(&src)),
			2,
		)
	}
}
