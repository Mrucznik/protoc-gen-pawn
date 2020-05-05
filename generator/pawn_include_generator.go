package generator

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

func GenerateIncludeFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".inc"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	genGeneratedHeader(gen, g)

	g.P("// ---------- Enums ----------")
	for _, enum := range file.Enums {
		genEnum(g, enum)
	}

	g.P("// ---------- Messages ----------")
	for _, message := range file.Messages {
		//TODO
		g.P(message.GoIdent.GoName)
	}

	g.P("// ---------- Natives ----------")
	for _, service := range file.Services {
		genNatives(g, service)
	}

	return g
}

func genNatives(g *protogen.GeneratedFile, service *protogen.Service) {
	g.P(service.Comments.Leading)
	g.P("// ----- ", service.GoName, " -----")
	for _, method := range service.Methods {
		g.P(method.Comments.Leading,
			"native bool:",
			strings.ToLower(extractCapitals(service.GoName)), "_", method.GoName,
			"(", genNativeParams(method), ");",
			method.Comments.Trailing)
		g.P()
	}
}

func genNativeParams(method *protogen.Method) string {
	var strBuilder strings.Builder

	//input params
	for _, param := range method.Input.Fields {
		genParam(&strBuilder, param, true)
		strBuilder.WriteString(", ")
	}

	//strBuilder params
	for _, param := range method.Output.Fields {
		genParam(&strBuilder, param, false)
		strBuilder.WriteString(", ")
	}

	out := strBuilder.String()
	return out[:len(out)-2]
}

func getParamsInfo(param *protogen.Field) (prefix string, array int, message bool) {
	array = 0
	switch param.Desc.Kind() {
	case protoreflect.EnumKind:
		prefix = param.Enum.GoIdent.GoName + ":"
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		prefix = "Float:"
	case protoreflect.StringKind:
		array += 1
	case protoreflect.BytesKind:
		array += 1
	case protoreflect.MessageKind:
		prefix = param.Message.GoIdent.GoName + ":"
		message = true
		//TODO: recursive message fields extraction
	case protoreflect.GroupKind:
		//deprecated
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
	}

	if param.Desc.IsList() {
		array += 1
	}

	//TODO: map support
	//if param.Desc.IsMap() {
	//
	//}

	return prefix, array, message
}

func genParam(builder *strings.Builder, param *protogen.Field, inputParam bool) {
	prefix, array, message := getParamsInfo(param)
	if array == 0 && !inputParam {
		builder.WriteRune('&')
	}
	if inputParam {
		builder.WriteString(fmt.Sprintf("const %si_%s", prefix, param.GoName))
	} else {
		builder.WriteString(fmt.Sprintf("%so_%s", prefix, param.GoName))
	}
	for i := 0; i < array; i++ {
		builder.WriteString("[]")
	}
	if message {
		builder.WriteRune('[')
		builder.WriteString(param.Message.GoIdent.GoName)
		builder.WriteRune(']')
	}
}

func genEnum(g *protogen.GeneratedFile, enum *protogen.Enum) {
	g.P(enum.Comments.Leading, "enum ", enum.Desc.Name())
	g.P("{")
	for idx, value := range enum.Values {
		genEnumValue(g, value, idx == len(enum.Values)-1)
	}
	g.P("};")
}

func genEnumValue(g *protogen.GeneratedFile, value *protogen.EnumValue, last bool) {
	if last {
		g.P(value.Comments.Leading, "\t", value.Desc.Name(), " = ", value.Desc.Number())
	} else {
		g.P(value.Comments.Leading,
			"\t", value.Desc.Name(), " = ", value.Desc.Number(), ", ")
	}
}
