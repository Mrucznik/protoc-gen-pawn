package generator

import (
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
	{{.RequestParams}}
    {{end}}

    // RPC call.
    Status status = API::Get().{{.ServiceName}}Stub()->{{.RPCName}}(&context, request, &response);
    API::Get().setLastStatus(status);

    {{if .ResponseParams}}
	// convert response to amx structure
	if(status.ok())
	{
		{{.ResponseParams}}
	}
    {{end}}

    return status.ok();
}`

const xd = `
            request.set_name(amx_GetCppString(amx, params[2]));
            request.set_description(amx_GetCppString(amx, params[3]));
            request.set_base_weight(params[4]);
            request.set_base_volume(params[5]);
            request.set_model_hash(params[6]);`
const xdd = `cell* addr = NULL;
                amx_GetAddr(amx, params[1], &addr);
                *addr = response.id();`

// GenerateNativeFile generates the contents of a _natives.cpp file.
// This file contains code that provides translation for pawn native call -> grpc call
func GenerateNativeFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_natives.cpp"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	genGeneratedHeader(gen, g)

	tmpl, err := template.New("cppNative").Parse(CppNativeTemplate)
	if err != nil {
		log.Fatalln(err)
	}

	for _, service := range file.Services {
		for _, method := range service.Methods {
			err = tmpl.Execute(g, cppNativeTemplateParams{
				ServiceName:    service.GoName,
				RPCName:        method.GoName,
				RequestType:    method.Input.GoIdent.GoName,
				ResponseType:   method.Output.GoIdent.GoName,
				NativeName:     getNativeName(service, method),
				NativeParams:   getNativeParams(method),
				RequestParams:  getInputFieldsCode(method.Input.Fields),
				ResponseParams: getOutputFieldsCode(method.Output.Fields),
			})
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return g
}

func getInputFieldsCode(fields []*protogen.Field) string {
	var builder strings.Builder

	for _, field := range fields {
		switch field.Desc.Kind() {
		case protoreflect.EnumKind:

		case protoreflect.FloatKind, protoreflect.DoubleKind:

		case protoreflect.StringKind:

		case protoreflect.BytesKind:

		case protoreflect.MessageKind:

		case protoreflect.GroupKind:
			//deprecated, do nothing
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
			protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
			protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		}
	}

	return builder.String()
}

func getOutputFieldsCode(fields []*protogen.Field) string {
	var builder strings.Builder

	for _, field := range fields {
		switch field.Desc.Kind() {
		case protoreflect.EnumKind:

		case protoreflect.FloatKind, protoreflect.DoubleKind:

		case protoreflect.StringKind:

		case protoreflect.BytesKind:

		case protoreflect.MessageKind:

		case protoreflect.GroupKind:
			//deprecated, do nothing
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
			protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
			protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind,
			protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		}
	}

	return builder.String()
}
