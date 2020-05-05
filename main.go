// protoc-gen-pawn is a plugin for the Google protocol buffer compiler to generate
// Pawn code. Install it by building this program and making it accessible within
// your PATH with the name:
//	protoc-gen-pawn
//
// The 'pawn' suffix becomes part of the argument for the protocol compiler,
// such that it can be invoked as:
//	protoc --pawn_out=paths=source_relative:. path/to/file.proto
//
// This generates Pawn bindings for the protocol buffer defined by file.proto.
// With that input, the output will be written to:
//	path/to/file.pb.go
//
// See the README and documentation for protocol buffers to learn more:
//	https://developers.google.com/protocol-buffers/

package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"protoc-gen-pawn/generator"
	"strings"
)

func main() {
	var (
		flags        flag.FlagSet
		plugins      = flags.String("plugins", "", "list of plugins to enable (supported values: natives)")
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}
	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		natives := false
		for _, plugin := range strings.Split(*plugins, ",") {
			switch plugin {
			case "natives":
				natives = true
			case "":
			default:
				return fmt.Errorf("protoc-gen-go: unknown plugin %q", plugin)
			}
		}
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generator.GenerateFile(gen, f)
			if natives {
				generator.GenerateNativeFiles(gen, f)
			}
		}
		gen.SupportedFeatures = generator.SupportedFeatures
		return nil
	})
}
