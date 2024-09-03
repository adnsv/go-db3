package schema

type table_grid [][]string

func measure_cell(s string) int {
	return len([]rune(s))
}

func measure_cells(row []string, col_widths *[]int) {
	for i := range row {
		w := measure_cell(row[i])
		if i < len(*col_widths) {
			if w > (*col_widths)[i] {
				(*col_widths)[i] = w
			}
		} else {
			*col_widths = append(*col_widths, w)
		}
	}
}
