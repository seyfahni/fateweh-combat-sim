package cmd

import "fmt"

type ConsolePrinter struct{}

func (p ConsolePrinter) Print(line string) error {
	_, err := fmt.Println(line)
	return err
}
