package meta

const (
	CategoryNone uint = iota
	CategoryBios
	CategoryFun
	CategoryMisc
	CategoryMeta
	CategoryAdminTools
	CategoryRoles
	CategoryVerification
)

var Descriptions = [CategoryVerification+1]string{
	CategoryNone: "Uncategorised",
	CategoryAdminTools: "Admin tools",
	CategoryBios: "Bios",
	CategoryFun: "Fun",
	CategoryMeta: "Meta",
	CategoryMisc: "Miscellaneous",
	CategoryRoles: "Reaction roles",
	CategoryVerification: "Verification",
}

func IterateCategories() *chan uint {
	c := make(chan uint)

	go func() {
		for i := CategoryNone; i <= CategoryVerification; i += 1 {
			c <- i
		}
		close(c)
	}()

	return &c
}