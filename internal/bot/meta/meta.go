package meta

const (
	CategoryNone uint = iota
	CategoryBios
	CategoryFun
	CategoryMisc
	CategoryMeta
	CategoryAdminTools
)

var Descriptions = [CategoryAdminTools+1]string{
	CategoryNone: "Uncategorised",
	CategoryAdminTools: "Admin tools",
	CategoryBios: "Bios",
	CategoryFun: "Fun",
	CategoryMeta: "Meta",
	CategoryMisc: "Miscellaneous",
}

func IterateCategories() *chan uint {
	c := make(chan uint)

	go func() {
		for i := CategoryNone; i <= CategoryAdminTools; i += 1 {
			c <- i
		}
		close(c)
	}()

	return &c
}