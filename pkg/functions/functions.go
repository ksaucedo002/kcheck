package functions

import (
	"fmt"

	"github.com/ksaucedo002/kcheck/pkg/models"
)

type ValidFunc func(atom models.Atom, args string) error
type MapFuncs map[string]ValidFunc

func Number(atom models.Atom, args string) error {
	fmt.Println("Numbre f:", atom, "Args:", args)
	return nil
}

func Lenght(atom models.Atom, args string) error {
	fmt.Println("Lenght f:", atom, "Args:", args)
	return nil
}
func MaxLenght(atom models.Atom, args string) error {
	fmt.Println("MaxLenght f:", atom, "Args:", args)
	return nil
}
func MinLenght(atom models.Atom, args string) error {
	fmt.Println("MinLenght f:", atom, "Args:", args)
	return nil
}
