package renderer

import (
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/Phillezi/common/utils/or"
)

// ExtractFields parses a Go text/template string and returns a list of field names it requires.
func ExtractFields(tmplText string, templates ...*template.Template) ([]string, error) {
	tmpl, err := or.Call(
		func() *template.Template { return or.Or(templates...) },
		func() *template.Template { return template.New("tmpl") },
	).Parse(tmplText)
	if err != nil {
		return nil, err
	}

	fieldSet := map[string]struct{}{}
	vars := map[string]string{} // tracks $var -> field path

	for _, t := range tmpl.Templates() {
		if t.Tree != nil && t.Tree.Root != nil {
			walk(t.Tree.Root, fieldSet, vars)
		}
	}

	fields := make([]string, 0, len(fieldSet))
	for field := range fieldSet {
		fields = append(fields, field)
	}
	return fields, nil
}

func walk(node parse.Node, fields map[string]struct{}, vars map[string]string) {
	switch n := node.(type) {
	case *parse.ListNode:
		for _, item := range n.Nodes {
			walk(item, fields, vars)
		}
	case *parse.ActionNode:
		walk(n.Pipe, fields, vars)
	case *parse.PipeNode:
		if len(n.Decl) > 0 && len(n.Cmds) > 0 {
			firstCmd := n.Cmds[0]
			if len(firstCmd.Args) > 0 {
				switch arg := firstCmd.Args[0].(type) {
				case *parse.FieldNode:
					base := "." + joinIdent(arg.Ident)
					for _, decl := range n.Decl {
						vars[decl.Ident[0]] = base + "[]" // assume it's a list when ranged
					}
				}
			}
		}
		for _, cmd := range n.Cmds {
			walk(cmd, fields, vars)
		}
	case *parse.CommandNode:
		for _, arg := range n.Args {
			walk(arg, fields, vars)
		}
	case *parse.FieldNode:
		fields["."+joinIdent(n.Ident)] = struct{}{}
	case *parse.VariableNode:
		// $row.Name etc.
		if len(n.Ident) > 1 {
			if base, ok := vars[n.Ident[0]]; ok {
				fields[base+"."+joinIdent(n.Ident[1:])] = struct{}{}
			} else {
				fields["$"+joinIdent(n.Ident)] = struct{}{}
			}
		} else {
			// just $length, $index
			fields["$"+n.Ident[0]] = struct{}{}
		}
	case *parse.IfNode:
		walk(n.Pipe, fields, vars)
		walk(n.List, fields, vars)
		if n.ElseList != nil {
			walk(n.ElseList, fields, vars)
		}
	case *parse.RangeNode:
		// Handle range base (e.g., .Rows)
		if n.Pipe != nil && len(n.Pipe.Cmds) > 0 && len(n.Pipe.Cmds[0].Args) > 0 {
			switch arg := n.Pipe.Cmds[0].Args[0].(type) {
			case *parse.FieldNode:
				fields["."+joinIdent(arg.Ident)] = struct{}{}
			}
		}
		// $length is implicitly created
		fields["$length"] = struct{}{}

		walk(n.Pipe, fields, vars)
		walk(n.List, fields, vars)
		if n.ElseList != nil {
			walk(n.ElseList, fields, vars)
		}
	case *parse.WithNode:
		walk(n.Pipe, fields, vars)
		walk(n.List, fields, vars)
		if n.ElseList != nil {
			walk(n.ElseList, fields, vars)
		}
	}
}

func joinIdent(parts []string) string {
	return strings.Join(parts, ".")
}
