package renderer

import "text/template"

func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func ne(a, b interface{}) bool {
	return a != b
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"add": add,
		"sub": sub,
		"len": func(x any) int {
			switch v := x.(type) {
			case []map[string]string:
				return len(v)
			case *[]map[string]string:
				return len(*v)
			case []string:
				return len(v)
			case *[]string:
				return len(*v)
			case []any:
				return len(v)
			case []Row:
				return len(v)
			case *[]Row:
				return len(*v)
			default:
				return 0
			}
		},
		"ne": ne,
	}
}
