package fileio

import (
	"fmt"
	"strings"
)

type ColumnDef struct {
	Header string
	Type   string // "int", "string", "float", "percent"
	Align  string // "left", "right", "center"
}

type TableData struct {
	Columns []ColumnDef
	Rows    [][]interface{}
}

type TeamTerritoryData struct {
	TeamID       int
	PlayerCount  int
	TileCount    int
	LandPercent  float64
	WorldPercent float64
	Players      []string
}

type TeamAnalysisData struct {
	TeamID       int
	PlayerCount  int
	CityCount    int
	UnitCount    int
	LandPercent  float64
	WorldPercent float64
	Players      []string
}

type TableFormatter struct{}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

func (tf *TableFormatter) PrintTable(data TableData) {
	if len(data.Columns) == 0 || len(data.Rows) == 0 {
		return
	}

	widths := tf.calculateColumnWidths(data.Columns, data.Rows)

	tf.printTopBorder(widths)
	tf.printHeaderRow(data.Columns, widths)
	tf.printSeparatorRow(widths)

	for _, row := range data.Rows {
		tf.printDataRow(data.Columns, row, widths)
	}

	tf.printBottomBorder(widths)
}

func (tf *TableFormatter) calculateColumnWidths(columns []ColumnDef, rows [][]interface{}) []int {
	widths := make([]int, len(columns))

	for i, col := range columns {
		widths[i] = len(col.Header)

		for _, row := range rows {
			if i < len(row) {
				valueStr := tf.formatValue(row[i], col.Type)
				if len(valueStr) > widths[i] {
					widths[i] = len(valueStr)
				}
			}
		}

		widths[i] += 2
	}

	return widths
}

func (tf *TableFormatter) formatValue(value interface{}, typeStr string) string {
	switch typeStr {
	case "int":
		return fmt.Sprintf("%d", value)
	case "string":
		return fmt.Sprintf("%s", value)
	case "float":
		return fmt.Sprintf("%.2f", value)
	case "percent":
		if f, ok := value.(float64); ok {
			return fmt.Sprintf("%.1f%%", f*100)
		}
		return fmt.Sprintf("%.1f%%", value)
	default:
		return fmt.Sprintf("%v", value)
	}
}

func (tf *TableFormatter) printTopBorder(widths []int) {
	fmt.Print("┌")
	for i, width := range widths {
		fmt.Print(strings.Repeat("─", width))
		if i < len(widths)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐")
}

func (tf *TableFormatter) printHeaderRow(columns []ColumnDef, widths []int) {
	fmt.Print("│")
	for i, col := range columns {
		fmt.Printf(" %-*s ", widths[i]-2, col.Header)
		if i < len(columns)-1 {
			fmt.Print("│")
		}
	}
	fmt.Println("│")
}

func (tf *TableFormatter) printSeparatorRow(widths []int) {
	fmt.Print("├")
	for i, width := range widths {
		fmt.Print(strings.Repeat("─", width))
		if i < len(widths)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

func (tf *TableFormatter) printDataRow(columns []ColumnDef, row []interface{}, widths []int) {
	fmt.Print("│")
	for i, col := range columns {
		value := ""
		if i < len(row) {
			value = tf.formatValue(row[i], col.Type)
		}

		switch col.Align {
		case "right":
			fmt.Printf(" %*s ", widths[i]-2, value)
		case "center":
			padding := (widths[i] - 2 - len(value)) / 2
			fmt.Printf(" %*s%s%*s ", padding, "", value, widths[i]-2-len(value)-padding, "")
		default: // left
			fmt.Printf(" %-*s ", widths[i]-2, value)
		}

		if i < len(columns)-1 {
			fmt.Print("│")
		}
	}
	fmt.Println("│")
}

func (tf *TableFormatter) printBottomBorder(widths []int) {
	fmt.Print("└")
	for i, width := range widths {
		fmt.Print(strings.Repeat("─", width))
		if i < len(widths)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("┘")
}
