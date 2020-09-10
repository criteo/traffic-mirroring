package expr

import (
	"fmt"
	"strings"
)

func ParseTmpl(str string) (Expr, error) {
	ostr := str
	parts := []Expr{}

	for len(str) > 0 {
		start := strings.IndexByte(str, '{')
		if start == -1 {
			parts = append(parts, identityExpr{v: str})
			break
		}

		if start > 0 {
			parts = append(parts, identityExpr{v: str[:start]})
		}
		str = str[start+1:]

		end := strings.IndexByte(str, '}')
		if end == -1 {
			return nil, fmt.Errorf("%q unterminated expression, missing '}'", ostr)
		}

		p, err := Parse(str[:end])
		if err != nil {
			return nil, err
		}

		parts = append(parts, p)
		str = str[end+1:]
	}

	if len(parts) == 1 {
		return parts[0], nil
	}

	return combineExpr{exprs: parts}, nil
}

func Parse(str string) (Expr, error) {
	ast, iss := env.Compile(str)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	ast, iss = env.Check(ast)
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	return progExpr{
		ast:  ast,
		prog: prg,
	}, nil
}
