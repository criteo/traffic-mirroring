package expr

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/shimmerglass/http-mirror-pipeline/mirror"
)

type Expr interface {
	Type() Type
	Eval(mirror.Request) (interface{}, error)
	Static() bool
}

type identityExpr struct {
	v interface{}
}

func (e identityExpr) Type() Type {
	switch e.v.(type) {
	case float64:
		return TypeNumber
	case int64:
		return TypeNumber
	case string:
		return TypeString
	case bool:
		return TypeBool
	default:
		return Type(fmt.Sprintf("%T", e.v))
	}
}

func (e identityExpr) Eval(r mirror.Request) (interface{}, error) {
	return e.v, nil
}

func (e identityExpr) Static() bool {
	return true
}

type progExpr struct {
	prog cel.Program
	ast  *cel.Ast
}

func (e progExpr) Type() Type {
	typ := e.ast.ResultType()

	switch typ {
	case decls.Double:
		return TypeNumber
	case decls.Int:
		return TypeNumber
	case decls.String:
		return TypeString
	case decls.Bool:
		return TypeBool
	default:
		return Type(typ.String())
	}
}

func (e progExpr) Eval(r mirror.Request) (interface{}, error) {
	v, _, err := e.prog.Eval(map[string]interface{}{
		"req": &r,
	})
	if err != nil {
		return nil, err
	}
	return v.Value(), nil
}

func (e progExpr) Static() bool {
	return false
}

type combineExpr struct {
	exprs []Expr
}

func (e combineExpr) Type() Type {
	return TypeString
}

func (e combineExpr) Eval(r mirror.Request) (interface{}, error) {
	s := strings.Builder{}
	for _, p := range e.exprs {
		v, err := p.Eval(r)
		if err != nil {
			return nil, err
		}

		if str, ok := v.(string); ok {
			s.WriteString(str)
		} else {
			s.WriteString(fmt.Sprint(v))
		}
	}

	return s.String(), nil
}

func (e combineExpr) Static() bool {
	for _, e := range e.exprs {
		if !e.Static() {
			return false
		}
	}

	return true
}
