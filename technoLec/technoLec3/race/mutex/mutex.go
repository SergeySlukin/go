package main //1
import (
	"sync"
	"log"
	"fmt"
)

type AccountProtected struct {
	sync.Mutex
	balance float64
}

func (a *AccountProtected) Balance() float64 {
	a.Lock()
	defer a.Unlock()
	return a.balance
}

func (a *AccountProtected) Deposit(amount float64)  {
	a.Lock()
	defer a.Unlock()
	log.Printf("depositing %f", amount)
	a.balance += amount
}

func (a *AccountProtected) Withdraw(amount float64)  {
	a.Lock()
	defer a.Unlock()
	if a.balance >= amount {
		log.Printf("withdraw %f", amount)
		a.balance -= amount
	}
}

func main()  {
	ac := AccountProtected{}
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				if j % 2 == 1 {
					ac.Withdraw(50)
					continue
				}
				ac.Deposit(50)
			}
		}()
	}

	fmt.Scanln()
	fmt.Println(ac.Balance())
}