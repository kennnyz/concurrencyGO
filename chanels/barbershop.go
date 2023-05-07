package main

import (
	"fmt"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	ClientsChan     chan string
	BarbersDoneChan chan bool
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		fmt.Println(barber, " goes to tha waiting room to check for clients ")

		for {
			if len(shop.ClientsChan) == 0 {
				fmt.Println("There is nothing to do, so ", barber, " takes a nap.")
				isSleeping = true
			}

			client, open := <-shop.ClientsChan

			if open {
				if isSleeping {
					fmt.Println(client, " wakes up ", barber)
					isSleeping = false
				}
				shop.cutHair(barber, client)
			} else {
				// shop closed so send barber to home
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	fmt.Println(barber, " start making hair cut to mister ", client)
	time.Sleep(shop.HairCutDuration)
	fmt.Println(barber, " finished to making hair cut to mister ", client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	fmt.Println("Send ", barber, " to home")
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	fmt.Println("Shop closing for the day ")
	close(shop.ClientsChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}

	close(shop.BarbersDoneChan)

	fmt.Println("----------------------------------------------------------------")
	fmt.Println("The barberShop is now closed for they day, everyone goes to home")
}

func (shop *BarberShop) addClient(client string) {
	fmt.Println("***Client ", client, " arrives!")
	if shop.Open {
		select {
		case shop.ClientsChan <- client:
			fmt.Println(client, " take a seat in the waiting room. ")
		default:
			fmt.Println("The waiting room is fuLl!, so ", client, " leaves :(")
		}
	} else {
		fmt.Println("The shop already closed so , ", client, " leaves ")
	}
}
