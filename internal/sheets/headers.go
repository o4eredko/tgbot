package sheets

type Headers struct {
	row []string
}

func NewHeadersFromRow(row []interface{}) *Headers {
	rowStr := make([]string, len(row))
	for i := range row {
		rowStr[i] = row[i].(string)
	}

	return &Headers{
		row: rowStr,
	}
}

func (h *Headers) GetByIndex(i int) string {
	if i >= len(h.row) {
		return ""
	}

	return h.row[i]
}
