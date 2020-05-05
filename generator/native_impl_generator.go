package generator

import "google.golang.org/protobuf/compiler/protogen"

// GenerateNativeFile generates the contents of a _natives.cpp file.
// This file contains code that provides translation for pawn native call -> grpc call
func GenerateNativeFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_natives.cpp"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	genGeneratedHeader(gen, g)

	return g
}
