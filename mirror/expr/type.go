package expr

import (
	"encoding/json"
	"fmt"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
)

type Type string

const (
	TypeString = "string"
	TypeNumber = "number"
	TypeBool   = "bool"
)

type AnyExpr struct {
	e Expr
}

func (e *AnyExpr) Eval(r mirror.Request) (interface{}, error) {
	return e.e.Eval(r)
}

func (e *AnyExpr) UnmarshalJSON(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	str, ok := v.(string)
	if !ok {
		e.e = identityExpr{v: v}
		return nil
	}

	expr, err := ParseTmpl(str)
	if err != nil {
		return err
	}

	e.e = expr
	return nil
}

func (e *AnyExpr) Static() bool {
	return e.e.Static()
}

type StringExpr struct {
	AnyExpr
}

func (e *StringExpr) UnmarshalJSON(b []byte) error {
	err := e.AnyExpr.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	if e.e.Type() != TypeString {
		return fmt.Errorf("unexpected type %s, expected %s", e.e.Type(), TypeString)
	}

	return nil
}

func (e *StringExpr) Eval(r mirror.Request) (string, error) {
	v, err := e.e.Eval(r)
	if err != nil {
		return "", err
	}

	return v.(string), nil
}

type NumberExpr struct {
	AnyExpr
}

func (e *NumberExpr) UnmarshalJSON(b []byte) error {
	err := e.AnyExpr.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	if e.e.Type() != TypeNumber {
		return fmt.Errorf("unexpected type %s, expected %s", e.e.Type(), TypeNumber)
	}

	return nil
}

func (e *NumberExpr) EvalInt(r mirror.Request) (int, error) {
	v, err := e.e.Eval(r)
	if err != nil {
		return 0, err
	}

	switch n := v.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case float64:
		return int(n), nil
	default:
		return 0, fmt.Errorf("unhandled type %T", v)
	}
}

func (e *NumberExpr) EvalFloat(r mirror.Request) (float64, error) {
	v, err := e.e.Eval(r)
	if err != nil {
		return 0, err
	}

	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, fmt.Errorf("unhandled type %T", v)
	}
}

type BoolExpr struct {
	AnyExpr
}

func (e *BoolExpr) UnmarshalJSON(b []byte) error {
	err := e.AnyExpr.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	if e.e.Type() != TypeBool {
		return fmt.Errorf("unexpected type %s, expected %s", e.e.Type(), TypeBool)
	}

	return nil
}

func (e *BoolExpr) Eval(r mirror.Request) (bool, error) {
	v, err := e.e.Eval(r)
	if err != nil {
		return false, err
	}

	return v.(bool), nil
}
