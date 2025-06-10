package template

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed file/model.gohtml
var modelTmpl string

type ModelData struct {
	PkgName    string
	PkType     string
	StructName string
}

func CreateModel() {

	data := ModelData{
		PkgName:    "model",
		StructName: "Student",
		PkType:     "uint64",
	}

	t := template.Must(template.New("model").Parse(modelTmpl))
	t.Execute(os.Stdout, data)
}
