package reporting

import "fmt"

type ComplexityNote struct {
	Operation string
	Time      string
	Space     string
	Reason    string
}

func ComplexityNotes() []ComplexityNote {
	return []ComplexityNote{
		{Operation: "尾插", Time: "O(1)", Space: "O(1)", Reason: "tail指针直接定位"},
		{Operation: "按编号删除", Time: "O(n)", Space: "O(1)", Reason: "线性查找节点后摘链"},
		{Operation: "编号插入排序", Time: "O(n^2)", Space: "O(n)", Reason: "逐项移动到已排序前缀"},
		{Operation: "借用日期快速排序", Time: "平均O(n log n)", Space: "O(log n)", Reason: "原地分区递归"},
		{Operation: "按借用人查询", Time: "O(n)", Space: "O(n)", Reason: "遍历并复制匹配记录"},
	}
}

func ComplexityText() string {
	lines := make([]string, 0)
	for _, note := range ComplexityNotes() {
		lines = append(lines, fmt.Sprintf("%s: time %s, space %s (%s)", note.Operation, note.Time, note.Space, note.Reason))
	}
	return joinLines(lines)
}

func joinLines(lines []string) string {
	result := ""
	for index, line := range lines {
		if index > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}
