package finances

type Wallet struct {
	ID      uint16
	Name    string
	Balance int64
	Closed  bool
}
