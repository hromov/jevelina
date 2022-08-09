package finances

type CatTotal struct {
	Category string
	Total    int
}

type CategorisedCashflow struct {
	Incomes  []CatTotal
	Expenses []CatTotal
}
