package main

import "log"

func main() {
	queue := make(chan bool, 3)
	go func() {
		queue <- true
		log.Println("Here")
	}()
	go func() {
		queue <- true
		log.Println("Here")
	}()
	go func() {
		queue <- true
		log.Println("Here")
	}()
	close(queue)
	for range queue {
		log.Println("Bam")
	}
}
