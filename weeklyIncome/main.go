package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	var bankBalance int
	var mutex sync.Mutex
	fmt.Printf("Initial account balance: $%d.00\n", bankBalance)
	incomes := []Income{
		{Source: "Main job", Amount: 500},
		{Source: "Part time", Amount: 50},
		{Source: "Gifts", Amount: 10},
		{Source: "Investments", Amount: 100},
	}

	for i, income := range incomes {
		wg.Add(1)
		go func(i int, income Income, wg *sync.WaitGroup, m *sync.Mutex) {
			defer wg.Done()
			for week := 1; week <= 52; week++ {
				mutex.Lock()
				temp := bankBalance
				temp += income.Amount

				bankBalance = temp
				mutex.Unlock()

				fmt.Printf("On week %d, you earned $%d.00 from %s\n", week, income.Amount, income.Source)
			}
		}(i, income, &wg, &mutex)
	}
	wg.Wait()

	fmt.Printf("Final balance in your bank account $%d.00\n", bankBalance)
}
