package renderer

type Row map[string]string

type QueryData struct {
	Rows []Row
}

func FromMapArr(mapArr []map[string]string) *QueryData {
	rows := make([]Row, len(mapArr))
	for i, m := range mapArr {
		rows[i] = Row(m)
	}
	return &QueryData{
		Rows: rows,
	}
}

type Renderer interface {
	Render(templateContent string, data QueryData) (string, error)
}
