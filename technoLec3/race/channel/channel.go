package main //4
import (
	"errors"
	"fmt"
)

type AccountAsync struct {
	balance     float64
	balanceChan chan float64
	deltaChan   chan float64
	errChan     chan error
}

func NewAccount(balance float64) (a *AccountAsync) {
	a = &AccountAsync{
		balance:     balance,
		balanceChan: make(chan float64),
		deltaChan:   make(chan float64),
		errChan:     make(chan error, 1),
	}
	go a.run()
	return a
}

func (a *AccountAsync) Balance() float64 {
	return <-a.balanceChan
}

func (a *AccountAsync) Deposit(amount float64) error {
	a.deltaChan <- amount
	return <-a.errChan
}

func (a *AccountAsync) Withdraw(amount float64) error {
	a.deltaChan <- -amount
	return <-a.errChan
}

func (a *AccountAsync) applyDelta(amount float64) error {
	str := "Кладем"
	if amount < 0 {
		str = "Снимаем"
	}
	fmt.Println(str, amount)
	newBalance := a.balance + amount
	if newBalance < 0 {
		return errors.New("Недостаточно средств.")
	}
	a.balance = newBalance
	return nil
}

func (a *AccountAsync) run() {
	var delta float64
	for {
		select {
		case delta = <-a.deltaChan:
			a.errChan <- a.applyDelta(delta)
		case a.balanceChan <- a.balance:
		}
	}
}

func main()  {
	a := NewAccount(100)
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
