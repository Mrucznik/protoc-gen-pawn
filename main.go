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
	"bytes"
	"fmt"
	"os"
)

func main() {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("error reading from stdin: %v", err))
	}
	//out, err := codeGenerator(buf.Bytes())
	//if err != nil {
	//	panic(err)
	//}
	_, err = os.Stdout.Write(buf.Bytes())
	if err != nil {
		panic(fmt.Errorf("error writing to stdout: %v", err))
	}
}
