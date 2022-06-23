package analyzer

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "reassign",
		Doc:  "Checks that package variables are not reassigned",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		state := &fileState{imports: make(map[string]struct{})}
		ast.Inspect(f, func(node ast.Node) bool {
			return inspect(pass, node, state)
		})
	}
	return nil, nil
}

type fileState struct {
	imports map[string]struct{}
}

func inspect(pass *analysis.Pass, node ast.Node, state *fileState) bool {
	if importSpec, ok := node.(*ast.ImportSpec); ok {
		if importSpec.Name != nil {
			state.imports[importSpec.Name.Name] = struct{}{}
		} else {
			n, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				return true
			}
			if idx := strings.IndexByte(n, '/'); idx != -1 {
				n = n[idx+1:]
			}
			state.imports[n] = struct{}{}
		}
		return true
	}

	assignStmt, ok := node.(*ast.AssignStmt)
	if !ok {
		return true
	}

	if len(assignStmt.Lhs) != 1 {
		// multiple assignment is not supported for variables in other packages
		return true
	}

	selector, ok := assignStmt.Lhs[0].(*ast.SelectorExpr)
	if !ok {
		return true
	}

	// TODO(anuraaga): Allow configuring what is matched, instead of only looking for Err
	if !strings.HasPrefix(selector.Sel.Name, "Err") {
		return true
	}

	selectIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return true
	}

	if _, ok := state.imports[selectIdent.Name]; ok {
		pass.Reportf(node.Pos(), "reassigning sentinel error %s in other package %s", selector.Sel.Name, selectIdent.Name)
	}

	fmt.Println(selectIdent.Name)

	return true
}
