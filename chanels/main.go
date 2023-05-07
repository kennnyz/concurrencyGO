package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	clientChan := make(chan string, 10)
	doneChan := make(chan bool)
	timeOpen := 10 * time.Second

	shop := BarberShop{
		ShopCapacity:    10,
		HairCutDuration: time.Second,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	shop.addBarber("Frank")
	shop.addBarber("Mike")
	shop.addBarber("Johnatan")
	shop.addBarber("Mark")
	shop.addBarber("Lampor")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	i := 1

	go func() {
		for {
			// get random number with average rate
			randomMillSeconds := rand.Int() % (2 * 100)

			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillSeconds)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	<-closed
}
