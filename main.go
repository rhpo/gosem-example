package main

import (
	"fmt"

	. "github.com/rhpo/gosem"
)

var (
	pont      = NomSemaphore()
	mutexNord = NomSemaphore()
	mutexSud  = NomSemaphore()
	nordBlock = NomSemaphore()
)

var (
	nbNord   = 0
	nbSud    = 0
	sudCount = 0
)

func main() {
	I(pont, 1)
	I(mutexNord, 1)
	I(mutexSud, 1)
	I(nordBlock, 0) // Nord completely blocked until Sud finishes

	// Processus Nord (completely blocked until Sud is done)
	Process(func() {
		P(nordBlock) // Each car waits for explicit signal from Sud
		Repeat(5, func() {

			P(mutexNord)
			nbNord++
			if nbNord == 1 {
				P(pont)
			}
			V(mutexNord)

			fmt.Println("Une voiture nord a passé")

			P(mutexNord)
			nbNord--
			if nbNord == 0 {
				V(pont)
			}
			V(mutexNord)
		})
	})

	// Processus Sud (goes first)
	Process(func() {
		Repeat(5, func() {
			P(mutexSud)
			nbSud++
			if nbSud == 1 {
				P(pont)
			}
			sudCount++
			V(mutexSud)

			fmt.Println("Une voiture sud a passé")

			P(mutexSud)
			nbSud--
			if nbSud == 0 {
				V(pont)

				// After the last Sud car, unblock all Nord cars
				if sudCount == 5 {
					V(nordBlock)
				}
			}
			V(mutexSud)
		})
	})

	Wait()
}
