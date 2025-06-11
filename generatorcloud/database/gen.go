package database

import (
	"github.com/acexy/gen"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"text/template"
)

type DatabaseGen struct {
	tableModelName          map[string]string
	disabledDefaultSettings bool
	gen                     *gen.Generator
	db                      *gorm.DB
	modelPkg                string
}

func NewDatabaseGen(db *gorm.DB, outPath string, tableModelName map[string]string) *DatabaseGen {
	d := &DatabaseGen{
		tableModelName: tableModelName,
		db:             db,
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithoutContext,
	})
	g.UseDB(db)
	d.gen = g
	return d
}

func NewDatabaseGenWithConfig(db *gorm.DB, tableModelName map[string]string, config gen.Config) *DatabaseGen {
	d := &DatabaseGen{
		tableModelName: tableModelName,
		db:             db,
	}
	g := gen.NewGenerator(config)
	g.UseDB(db)
	d.gen = g
	return d
}

func (d *DatabaseGen) DisableDefaultSettings() {
	d.disabledDefaultSettings = true
}

func (d *DatabaseGen) SetModelPkg(modelPkg string) {
	d.modelPkg = modelPkg
}

func (d *DatabaseGen) rawGen() *gen.Generator {
	return d.gen
}

func (d *DatabaseGen) dBType() string {
	switch d.db.Dialector.(type) {
	case *mysql.Dialector:
		return "mysql"
	case *postgres.Dialector:
		return "postgres"
	default:
		return "unknown"
	}
}

func (d *DatabaseGen) render(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}

type DBTypeData struct {
	StructName string
	DBType     string
}

func (d *DatabaseGen) Create() {
	NewModelGen(d).Create()
}
