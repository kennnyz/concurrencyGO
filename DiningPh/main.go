package main

import (
	"fmt"
	"sync"
	"time"
)

// The classic dining philosophers problem
// https://en.wikipedia.org/wiki/Dining_philosophers_problem

type Philosopher struct {
	name      string
	leftFork  int
	rightFork int
}

var philosophers = []Philosopher{
	{"Aristotle", 4, 0},
	{"Plato", 0, 1},
	{"Socrates", 1, 2},
	{"Descartes", 2, 3},
	{"Kant", 3, 4},
}

var hunger = 3
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 1 * time.Second

func main() {
	fmt.Println("Dining Philosophers problem.")
	fmt.Println("----------------------------")

	dine()

	fmt.Println("The table is empty.")

}

func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()
}

func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("%s is sea\n", philosopher.name)

	seated.Done()

	seated.Wait()

	for h := 0; h < hunger; h++ {

		if philosopher.leftFork < philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("%s has picked up right fork\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s has picked up left fork\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s has picked up left fork\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("%s has picked up right fork\n", philosopher.name)
		}

		fmt.Printf("%s is take a both fork and his eating\n", philosopher.name)

		time.Sleep(eatTime)

		fmt.Printf("%s is thinking\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("%s has put down forks\n", philosopher.name)
	}

	fmt.Println(philosopher.name, "is satisified.")
	fmt.Println(philosopher.name, "left the table.")

}
