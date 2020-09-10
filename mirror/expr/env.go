package expr

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	env *cel.Env
)

type customLib struct{}

func (customLib) CompileOptions() []cel.EnvOption {
	reqType := decls.NewObjectType("mirror.Request")
	dec := []*exprpb.Decl{
		decls.NewVar("req", reqType),
		decls.NewFunction("header",
			decls.NewInstanceOverload(
				"header",
				[]*exprpb.Type{reqType, decls.String},
				decls.String,
			),
		),
	}
	for name, v := range mirror.Method_value {
		c := &exprpb.Constant{ConstantKind: &exprpb.Constant_Int64Value{Int64Value: int64(v)}}
		dec = append(dec, decls.NewConst(name, decls.Int, c))
	}
	for name, v := range mirror.HTTPVersion_value {
		c := &exprpb.Constant{ConstantKind: &exprpb.Constant_Int64Value{Int64Value: int64(v)}}
		dec = append(dec, decls.NewConst(name, decls.Int, c))
	}
	return []cel.EnvOption{
		cel.Types(&mirror.Request{}),
		cel.Declarations(dec...),
	}
}

func (customLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "header",
				Binary: func(lhs, rhs ref.Val) ref.Val {
					req := lhs.Value().(*mirror.Request)
					name := rhs.Value().(string)

					v := req.Headers[name].Values
					if len(v) == 0 {
						return types.String("")
					}

					return types.String(v[0])
				},
			},
		),
	}
}

func init() {
	e, err := cel.NewEnv(cel.Lib(customLib{}))
	if err != nil {
		panic(err)
	}

	env = e
}
