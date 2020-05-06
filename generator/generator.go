package generator

import (
	"github.com/gobuffalo/packr"
	"google.golang.org/protobuf/types/pluginpb"
)

var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

var TemplateBox packr.Box

func init() {
	TemplateBox = packr.NewBox("./generator/templates")
}
