# Cloud Generator

Cloud Generator creates the relational data-access and HTTP foundation used by
the Golang Acexy cloud modules. It reads MySQL or PostgreSQL table metadata and
generates code compatible with:

- `starter-gorm`
- `cloud-database`
- `cloud-web`
- `starter-gin`

The generated layers follow the project conventions:

```text
database table
    -> model and DTOs
    -> GORM Mapper
    -> cloud-database Repository
    -> cloud-web business service
    -> Gin base router
```

## Ecosystem Role

The generator is an application-development tool rather than a runtime starter. It converts database schema metadata into code that already follows the current `starter-gorm`, `cloud-database`, `cloud-web`, and `starter-gin` contracts, leaving business-specific extensions to the application.

## Requirements

- Go 1.26.7 or later
- MySQL or PostgreSQL
- An existing Go module for the generated application
- Database access to inspect the target tables

## Installation

```bash
go get github.com/golang-acexy/cloud-generator
```

## Basic Usage

Create the target tables before running the generator, then configure their
model and router names:

```go
db, err := gorm.Open(mysql.Open(
    "YOUR_USERNAME:YOUR_PASSWORD@tcp(DB_HOST:3306)/app?charset=utf8mb4&parseTime=True&loc=Local",
))
if err != nil {
    return err
}

generator := generatorcloud.NewGen(
    db,
    "/workspace/cloud-app/internal/model",
    []generatorcloud.TableConfig{
        {
            TableName: "demo_department",
            ModelName: "Department",
            Router: &generatorcloud.RouterConfig{
                BaseRouter: &generatorcloud.BaseRouter{
                    RelativeModelPath: []string{"..", "handler", "rest", "adm"},
                    GroupPath:         "adm/department",
                },
            },
        },
    },
)

generator.SetIncludeModelPkgPath(
    "github.com/example/cloud-app/internal/model",
)
generator.SetModelBase(&generatorcloud.ModelBase{
    DefaultTimeRangeField:  "created_at",
    AllowedTimeRangeFields: []string{"created_at", "updated_at"},
})
generator.SetRepoRelativeModelPath(
    []string{"..", "service", "repo"},
)
generator.SetServiceRelativeModelPath(
    []string{"..", "service", "biz"},
)

if err = generator.Create(); err != nil {
    return err
}
```

`SetIncludeModelPkgPath` must contain the Go import path of the generated model
package. Repository, service, and router imports are derived from this path and
their configured relative paths; the target packages and their dependencies must already be available. DTO conversion imports such as `copier` are resolved by import processing; it does not install missing dependencies or type-check the application.

`SetModelBase` must configure a non-empty default time-range field and allowed
field list. The default must be present in the allowed list. Names are normalized
to snake case and apply to every model produced by the generator run, so this
configuration is best suited to shared audit columns. Column existence is not checked. These settings are required even when service generation is disabled. Settings and path slices are borrowed; finish configuration before generation and do not run a Generator concurrently.

## Generated Output

With the configuration above, the default application structure is:

```text
internal/
├── model/
│   └── department_gen.go
├── service/
│   ├── repo/
│   │   └── department_repo.go
│   └── biz/
│       └── department_biz.go
└── handler/
    └── rest/
        └── adm/
            └── department_router.go
```

### Model and DTOs

Each table generates:

- the GORM model
- `TableName` and `DBType`
- save DTO (`SDTO`)
- modify DTO (`MDTO`)
- query DTO (`QDTO`)
- response DTO (`DTO`)
- model/DTO conversion methods

Database-generated fields are excluded from save and modify DTOs:

- `ID`
- creation timestamps
- update timestamps
- deletion timestamps

The concrete model primary-key type is propagated to the generated business
service and router instead of being fixed to `int64`.

### Mapper and Repository

The generated Mapper embeds `gormstarter.BaseMapper[T]` and implements
`WithTxMapper`:

```go
type DepartmentMapper struct {
    gormstarter.BaseMapper[model.Department]
}

func (m DepartmentMapper) WithTxMapper(tx *gorm.DB) DepartmentMapper {
    return DepartmentMapper{
        BaseMapper: m.GetBaseMapperWithTx(tx),
    }
}
```

The generated Repository preserves its concrete business type:

```go
type DepartmentRepo struct {
    rds.Repository[DepartmentRepo, DepartmentMapper, model.Department]
}

func (r DepartmentRepo) Columns() *model.DepartmentColumns {
    return r.RawMapper().Columns()
}
```

Business services can use `repo.Columns()` together with `repo.Wrapper()`,
`repo.PageWrapper(number, size)`, or `repo.ModifyWrapper()` without constructing a Mapper
or importing `starter-gorm`.

A package-level base Repository is initialized once. `NewDepartmentRepo`
returns a value copy, so transaction-bound repositories cannot mutate the
shared base Repository.

### Business Service

The generated service implements `webcloud.BaseBizService` and uses business
verbs. Its common operations include:

- `Save`, `SaveWithoutZeroFields`, and `SaveBatch`
- `QueryByID`, `QueryByIDs`, `QueryByCond`, and `QueryPage`
- `ExistsByID`, `CountByCond`, and `CountByMap`
- `ModifyByID`, `ModifyByCond`, and their zero-field or map variants
- `RemoveByID`, `RemoveByIDs`, and condition variants

`DefaultOrderBy` and `MaxQuerySize` can be configured globally:

```go
generator.SetServiceBase(&generatorcloud.ServiceBase{
    DefaultOrderBy: "id desc",
    MaxQuerySize:   100,
})
```

Generated typed conditions use value semantics. Repository query results are
received directly as `*T`, `[]*T`, or `Pager[T]`. Generated page methods validate and convert
`webcloud.TimeRange` values to `rds.TimeRange`, then delegate count,
ordering, limit, offset, and data loading to the Repository pagination methods.
Generic REST conditions remain Map-based so explicit zero values are preserved.
Direct service calls accept empty base query conditions; default HTTP `/query` and `/query-one` handlers require nonempty maps. `MaxQuerySize` limits `BaseQuery` and `QueryByCond`, not `QueryByIDs`, counts, or pages. HTTP page binding enforces 1–2000; direct service calls do not inherit that cap. Obvious empty update and delete
conditions remain rejected by the Repository layer; this does not prove every accepted predicate is selective. Generated services use
Repository Query constructors instead of issuing raw GORM queries.

## Authority-Aware Routers

Use `BaseRouterWithAuthority` when every base operation must be restricted by
the current identity:

```go
BaseRouterWithAuthority: &generatorcloud.BaseRouterWithAuthority{
    BaseRouter: generatorcloud.BaseRouter{
        RelativeModelPath: []string{"..", "handler", "rest", "usr"},
        GroupPath:         "usr/employee",
    },
    AuthorityFetchCode:   "biz.UsrAuthorityFetch",
    AuthorityStructField: "UserID",
    AuthorityColumn:      "user_id",
}
```

The generated router passes:

```go
webcloud.AuthorityDataField{
    StructField: "UserID",
    Column:      "user_id",
}
```

Authority behavior is provided by `cloud-web`:

- save requests have the authority field forcibly overwritten
- update requests cannot change ownership
- query, update, and delete conditions automatically include the authority
  database column
- the client cannot override the system authority condition

This behavior applies only when requests go through an authority-aware base
router. Direct service and Repository calls must enforce their own business
authorization.

## Extending Generated Code

Keep custom code in separate files in the generated package:

```go
// employee_extension.go
package repo

func (r EmployeeRepo) CountByDepartmentID(id int64) (int64, error) {
    return r.CountByMap(rds.NewMapQuery(map[string]any{"department_id": id}))
}
```

The same pattern can extend generated business services and routers. Generated
router files pass the concrete router to `RegisterBaseHandlers`, so a custom
router method with the same signature overrides the embedded `BaseRouter`
handler while all other handlers retain their default implementations. They
also contain a registration point for additional custom handlers.

Repository files are overwritten on regeneration. Existing business service and router files are skipped and do not automatically receive template updates. Model files are owned by the model
generator and should not contain handwritten business logic.

## End-to-End Demo

The integration fixture under `test/simple_demo`:

1. drops and recreates `demo_department` and `demo_employee`
2. generates code into the sibling `cloud-simple-demo` module
3. generates normal and authority-aware routers

Compilation and HTTP verification are separate steps in the target application; the generator test itself does not perform them.

Run generation:

```bash
go test -count=1 -v ./test/simple_demo
```

Warning: this test deletes and recreates the following tables in the configured
`test` database:

```text
demo_department
demo_employee
```

Run the generated application verification from `cloud-simple-demo`:

```bash
./scripts/verify.sh
```

The script verifies:

- generated model, Mapper, Repository, service, and router integration
- save, query, list, page, update, and delete routes
- authority field overwrite
- cross-user query and delete isolation
- custom Repository, service, and router extensions

## Design Notes

- MySQL and PostgreSQL are supported.
- Unsupported database dialectors return `ErrUnsupportedDatabase`.
- Missing model package paths return `ErrModelPackageRequired`.
- Unsupported HTTP primary-key types return `ErrUnsupportedIDType`.
- Invalid authority and router configurations fail generation explicitly.
- Local rendering and file-operation errors are returned, but the metadata dependency can panic. Generation is incremental and does not roll back written files.
- Generated source files use mode `0644`; generated directories use `0755`.

## Generation Boundaries

- Use conventional `model`, `repo`, and `biz` package names: some templates reference these names directly, even when import paths are configurable.
- `FilePrefix` adds a directory component, not a filename prefix; package names still come from the last `RelativeModelPath` component. Paths and injected Go expressions are trusted configuration, not sandboxed inputs.
- Plain and authority routers must use distinct output directories. If both resolve to the same path, the plain router is written first and the authority variant is skipped. Separate files in the same package can also collide on generated identifiers.
- Authority configuration checks only that three strings are nonempty. Ensure the callback type matches the ID, the save DTO retains the ownership field, and the ownership column does not accidentally replace the resource `id`.
- Metadata accepts string IDs, but generated save error paths currently return numeric `0`, so string-ID services do not compile without manual changes. Missing `ID` fields fall back to `int64` metadata but still produce invalid service references.
- Disabling service generation does not remove router imports or constructor calls; supply the required service yourself.
- Generated services and routers are shared singletons. Direct service calls have no automatic HTTP authority, request context propagation, or multi-operation transaction.

## Safe Template Checks

Run `go test -race ./generatorcloud` for isolated model and template checks. These checks validate rendering and Go syntax, not full generated application type compatibility. Tests under `test/` can access databases and overwrite files; the simple demo also recreates tables and writes into another module.
