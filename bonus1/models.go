package analytics

type Sale struct {
	Product          string
	Category         string
	Quantity         int
	UnitPriceKopecks int64
	Cancelled        bool
}

type Report struct {
	RevenueKopecks    int64
	UnitsSold         int
	RevenueByCategory map[string]int64
	TopProduct        string
	InvalidSales      int
}

func NewEmptyReport() Report {
	return Report{
		RevenueByCategory: make(map[string]int64),
	}
}

func (rep *Report) AddSale(sale Sale) {
	profit := sale.UnitPriceKopecks * int64(sale.Quantity)

	rep.RevenueKopecks += profit
	rep.UnitsSold += sale.Quantity
	rep.RevenueByCategory[sale.Category] += profit
}

func isValid(sale Sale) bool {
	return len(sale.Product) > 0 && len(sale.Category) > 0 &&
		sale.Quantity > 0 && sale.UnitPriceKopecks > 0
}
