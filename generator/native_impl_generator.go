package generator

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"log"
	"strings"
	"text/template"
)

type cppNativeTemplateParams struct {
	ServiceName    string
	RPCName        string
	RequestType    string
	ResponseType   string
	NativeName     string
	NativeParams   string
	RequestParams  string
	ResponseParams string
}

const CppNativeTemplate = `
// native {{.NativeName}}({{.NativeParams}});
cell Natives::{{.NativeName}}(AMX *amx, cell *params) {
    {{.RequestType}} request;
    {{.ResponseType}} response;
    ClientContext context;
    {{if .RequestParams}}
	// construct request from params
{{.RequestParams}}{{end}}
    // RPC call.
    Status status = API::Get()->{{.ServiceName}}Stub()->{{.RPCName}}(&context, request, &response);
    API::Get().setLastStatus(status);
    {{if .ResponseParams}}
	// convert response to amx structure
	if(status.ok())
	{
{{.ResponseParams}}
	}{{end}}
    return status.ok();
}`

const CppAsyncNativeTemplate = `
// native {{.NativeName}}Async({{.NativeParams}});
cell Natives::{{.NativeName}}(AMX *amx, cell *params) {
	
}`

const CppNativeClientStreamingTemplate = `
// native {{.NativeName}}({{.NativeParams}});
cell Natives::{{.NativeName}}(AMX *amx, cell *params) {
	
}`

const CppNativeServerStreamingTemplate = `
// native {{.NativeName}}({{.NativeParams}});
cell Natives::{{.NativeName}}(AMX *amx, cell *params) {
	
}`

const CppNativeBidirectionalStreamingTemplate = `
// native {{.NativeName}}({{.NativeParams}});
cell Natives::{{.NativeName}}(AMX *amx, cell *params) {
	
}`

// GenerateNativeDefinitions generates the content of _native_definitions.cpp
// These files contains list of native definitions to paset in natives.hpp
func GenerateNativeDefinitions(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_native_definitions.cpp"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	for _, service := range file.Services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			g.P(fmt.Sprintf("AMX_DEFINE_NATIVE(%s)", getNativeName(service, method)))
		}
	}

	return g
}

// GenerateNativeFile generates the content of a _natives.cpp file.
// These files contains code that provides translation for pawn native call -> grpc call
func GenerateNativeFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_natives.cpp"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	//genGeneratedHeader(gen, g)

	tmpl, err := template.New("cppNative").Parse(CppNativeTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	for _, service := range file.Services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() && method.Desc.IsStreamingServer() {
				continue
			} else if method.Desc.IsStreamingClient() {
				continue
			} else if method.Desc.IsStreamingServer() {
				continue
			}
			err = tmpl.Execute(g, cppNativeTemplateParams{
				ServiceName:    service.GoName,
				RPCName:        method.GoName,
				RequestType:    method.Input.GoIdent.GoName,
				ResponseType:   method.Output.GoIdent.GoName,
				NativeName:     getNativeName(service, method),
				NativeParams:   getNativeParams(method),
				RequestParams:  getInputFieldsCode(method.Input.Fields),
				ResponseParams: getOutputFieldsCode(method.Output.Fields, len(method.Input.Fields)),
			})
			if err != nil {
				log.Println(err)
			}
			g.P()
		}
	}

	return g
}

func getInputFieldsCode(fields []*protogen.Field) string {
	var builder strings.Builder

	for idx, field := range fields {
		if field.Desc.IsList() {
			builder.WriteString("\t\t// todo: list\n")
			continue
		}

		if field.Desc.IsMap() {
			builder.WriteString("\t\t// todo: map\n")
			continue
		}

		switch field.Desc.Kind() {
		case protoreflect.EnumKind:
			builder.WriteString(fmt.Sprintf("\trequest.set_%s(static_cast<%s>(params[%d]));\n",
				field.Desc.Name(), field.Enum.Desc.Name(), idx+1))
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			builder.WriteString(fmt.Sprintf("\trequest.set_%s(amx_ctof(params[%d]));\n",
				field.Desc.Name(), idx+1))
		case protoreflect.StringKind:
			builder.WriteString(fmt.Sprintf("\trequest.set_%s(amx_GetCppString(amx, params[%d]));\n",
				field.Desc.Name(), idx+1))
		case protoreflect.BytesKind:
			builder.WriteString("\t// TODO: bytes\n")
		case protoreflect.MessageKind:
			builder.WriteString("\t// TODO: message\n")
		case protoreflect.GroupKind:
			//deprecated, do nothing
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
			protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
			protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
			builder.WriteString(fmt.Sprintf("\trequest.set_%s(params[%d]);\n",
				field.Desc.Name(), idx+1))
		}
	}

	return builder.String()
}

func getOutputFieldsCode(fields []*protogen.Field, beginIndex int) string {
	var builder strings.Builder

	builder.WriteString("\t\tcell* addr = nullptr;\n")
	for idx, field := range fields {
		if field.Desc.IsList() {
			builder.WriteString("\t\t// todo: list\n")
			continue
		}

		if field.Desc.IsMap() {
			builder.WriteString("\t\t// todo: map\n")
			continue
		}

		switch field.Desc.Kind() {
		case protoreflect.FloatKind, protoreflect.DoubleKind:
			builder.WriteString(fmt.Sprintf(
				"\t\tamx_GetAddr(amx, params[%d], &addr);\n"+
					"\t\tfloat %s = response.%s();\n"+
					"\t\t*addr = amx_ftoc(%s);\n",
				idx+1+beginIndex, field.Desc.Name(), field.Desc.Name(), field.Desc.Name()))
		case protoreflect.StringKind:
			builder.WriteString(fmt.Sprintf(
				"\t\tamx_SetCppString(amx, params[%d], response.%s(), 256);\n",
				idx+1+beginIndex, field.Desc.Name()))
		case protoreflect.BytesKind:
			builder.WriteString("\t\t// TODO: bytes\n")
		case protoreflect.MessageKind:
			builder.WriteString("\t\t// TODO: message\n")
		case protoreflect.GroupKind:
			//deprecated, do nothing
		case protoreflect.EnumKind,
			protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
			protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
			protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
			builder.WriteString(fmt.Sprintf(
				"\t\tamx_GetAddr(amx, params[%d], &addr);\n"+
					"\t\t*addr = response.%s();\n",
				idx+1+beginIndex, field.Desc.Name()))
		}
	}

	return builder.String()
}
