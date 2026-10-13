package generatorcloud

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/acexy/gen/core/generate"
	"github.com/acexy/gen/core/model"
)

func TestTemplateSyntaxAndPublicComments(t *testing.T) {
	keep := func(string) bool { return false }
	meta := &generate.QueryStructMeta{ModelStructName: "Teacher", Fields: []*model.Field{{Name: "ID", Type: "int64"}, {Name: "Name", Type: "string"}}}
	cases := []struct {
		name, source, prefix string
		data                 any
	}{
		{"dto", dtoTmpl, "package model\n", DtoData{QueryStructMeta: meta, IsSExcluded: keep, IsMExcluded: keep, IsQExcluded: keep, IsDExcluded: keep}},
		{"columns", columnsTmpl, "package model\n", buildColumnsData("Teacher", meta)},
		{"method", methodTmpl, "package model\n", ModelData{StructName: "Teacher", DBType: "gormstarter.DBTypeMySQL"}},
		{"repo", repoTmpl, "", RepoData{StructName: "Teacher", ParamName: "teacher", Pkg: "example.com/app/model", PKG: "repo"}},
		{"service", serviceTmpl, "", ServiceData{StructName: "Teacher", ParamName: "teacher", ModelPKG: "example.com/app/model", RepoPKG: "example.com/app/repo", PKG: "biz", IDType: "int64", MaxQuerySize: 500}},
		{"router", routerImpl, "", RouterData{ModelStructName: "Teacher", ParamName: "teacher", PkgName: "api", ModelPKG: "example.com/app/model", BizPKG: "example.com/app/biz", IDType: "int64"}},
		{"authority", routerImpl, "", RouterData{ModelStructName: "Teacher", ParamName: "teacher", PkgName: "api", ModelPKG: "example.com/app/model", BizPKG: "example.com/app/biz", IDType: "int64", WithAuthority: true, AuthorityFetchCode: "biz.FetchAuthority", AuthorityStructField: "UserID", AuthorityColumn: "user_id"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			output.WriteString(tc.prefix)
			if err := (&Generator{}).render(tc.source, &output, tc.data); err != nil {
				t.Fatal(err)
			}
			file, err := parser.ParseFile(token.NewFileSet(), tc.name+".go", output.Bytes(), parser.ParseComments|parser.AllErrors)
			if err != nil {
				t.Fatalf("invalid rendered Go: %v\n%s", err, output.String())
			}
			ast.Inspect(file, func(node ast.Node) bool {
				switch decl := node.(type) {
				case *ast.FuncDecl:
					if decl.Name.IsExported() && decl.Doc == nil {
						t.Errorf("missing method comment: %s", decl.Name)
					}
				case *ast.GenDecl:
					for _, spec := range decl.Specs {
						if typ, ok := spec.(*ast.TypeSpec); ok && typ.Name.IsExported() && decl.Doc == nil && typ.Doc == nil {
							t.Errorf("missing type comment: %s", typ.Name)
						}
					}
				case *ast.StructType:
					for _, field := range decl.Fields.List {
						for _, name := range field.Names {
							if name.IsExported() && field.Doc == nil {
								t.Errorf("missing field comment: %s", name)
							}
						}
					}
				}
				return true
			})
		})
	}
}
