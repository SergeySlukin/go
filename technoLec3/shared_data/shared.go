package main //2
import "fmt"

type Account struct {
	balance float64
}

func (a *Account) Balance() float64 {
	return a.balance
}

func (a *Account) Deposit(amount float64) {
	fmt.Printf("depositing %f\n", amount)
	a.balance += amount
}

func (a *Account) Withdraw(amount float64) {
	fmt.Printf("withdrawing %f\n", amount)
	if amount > a.balance {
		return
	}
	a.balance -= amount
}

func main() {
	a := Account{}

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				if j % 2 == 1 {
					a.Withdraw(50)
					continue
				}
				a.Deposit(50)
			}
		}()
	}
	fmt.Scanln()
	fmt.Println(a.Balance())
}
