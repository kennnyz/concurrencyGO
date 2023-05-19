package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizza = 10

var pizzaMade, pizzaFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

type PizzaOrder struct {
	pizzaOrder int
	message    string
	success    bool
}

func pizzeria(pizzaMaker *Producer) {
	// keep tracking which pizza we are making
	i := 0
	//run until we got a quit signal
	for {
		// try to make a pizza
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaOrder
			select {
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizza {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd > 5 {
			pizzaFailed++
		} else {
			pizzaMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaOrder: pizzaNumber,
			message:    msg,
			success:    success,
		}
	}
	return &PizzaOrder{
		pizzaOrder: pizzaNumber,
	}
}

func main() {
	color.Cyan("This pizzeria its open for business!")
	color.Cyan("-------------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}
	// run the producer
	go pizzeria(pizzaJob)

	// create a consumer

	for i := range pizzaJob.data {
		if i.pizzaOrder <= NumberOfPizza {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out of delivery! ", i.pizzaOrder)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad! ")

			}
		} else {
			color.Cyan("Done making pizzas ")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("Error closing chanel! ", err)
			}
		}
	}

	// print the result
}
