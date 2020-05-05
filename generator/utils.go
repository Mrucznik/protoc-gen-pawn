package generator

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"protoc-gen-pawn/version"
	"strings"
)

func genGeneratedHeader(gen *protogen.Plugin, g *protogen.GeneratedFile) {
	g.P("// Code generated by protoc-gen-pawn. DO NOT EDIT.")

	g.P("// versions:")
	protocGenPawnVersion := version.String()
	protocVersion := "(unknown)"
	if v := gen.Request.GetCompilerVersion(); v != nil {
		protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
	}
	g.P("// \tprotoc-gen-pawn ", protocGenPawnVersion)
	g.P("// \tprotoc          ", protocVersion)

	g.P()
}

func extractCapitals(input string) string {
	var output strings.Builder

	for _, c := range input {
		if c >= 'A' && c <= 'Z' {
			output.WriteByte(byte(c))
		}
	}

	return output.String()
}