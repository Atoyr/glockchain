package core

func (cli *CLI) createTransaction(from, to string, amount int) {
	wallets := NewWallets()
	wallet := wallets.Wallets[from]
	NewTransaction(wallet, []byte(to), amount)
}
