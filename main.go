package main

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

func main() {
	protogen.Options{}.Run(generate)
}

func generate(gen *protogen.Plugin) error {
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		generateFile(gen, f)
	}
	return nil
}

func generateFile(gen *protogen.Plugin, f *protogen.File) {
	if len(f.Enums) == 0 {
		return
	}
	pkg := string(f.Desc.Package())
	g := gen.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.enums.go", f.GoImportPath)
	g.P("// Code generated by protoc-gen-go-enums. DO NOT EDIT.")
	g.P("// source: ", f.Desc.Path())
	g.P()
	g.P("package ", f.GoPackageName)
	g.P()

	for _, enum := range f.Enums {
		g.P("const (")
		for _, val := range enum.Values {
			if val.Desc.Options().(*descriptorpb.EnumValueOptions).GetDeprecated() {
				g.P("// Deprecated: Do not use.")
			}
			g.P(val.Desc.Name(), " = ", golangValue(pkg, val))
		}
		g.P(")")
		g.P()
	}
}

func golangValue(pkgName string, e *protogen.EnumValue) string {
	typeName := strings.TrimPrefix(string(e.Parent.Desc.FullName()), ".")
	typeName = strings.TrimPrefix(typeName, pkgName+".")

	parts := strings.Split(typeName, ".")
	if len(parts) > 1 {
		typeName = strings.Join(parts[0:len(parts)-1], "_")
	} else {
		typeName = parts[0]
	}

	return typeName + "_" + string(e.Desc.Name())
}