package markdown

import (
	"strings"
)

type Block struct {
	Title string
	Content string
}

func SplitByHeader(input string) []*Block {
	x := strings.Split(input, "# ") // split by header

	var nx []*Block
	for _, item := range x {
		if item == "" {
			// ignore any empty items
			continue
		}

		y := strings.Split(item, "\n")

		//{
		//	// filter out blank segments
		//	var n int
		//	for _, val := range y {
		//		if val != "" {
		//			y[n] = val
		//			n += 1
		//		}
		//	}
		//	y = y[:n]
		//}

		var title string
		var content []string
		for i, subitem := range y {
			if i == 0 {
				// assume the first line is the header
				title = strings.TrimSpace(subitem)
			} else {
				content = append(content, subitem)
			}
		}

		nx = append(nx, &Block{
			Title: title,
			Content: strings.TrimSpace(strings.Join(content, "\n")),
		})
	}

	return nx
}
