package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type menuItem struct {
	label  string
	action func() error
}

func runMenu(title string, items []menuItem) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("==", title, "==")
		for i, item := range items {
			fmt.Printf("%d) %s\n", i+1, item.label)
		}
		fmt.Println("0) Exit")
		fmt.Print("Choose option: ")

		raw, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read input: %w", err)
		}

		choice := strings.TrimSpace(raw)
		if choice == "0" {
			fmt.Println("Bye.")
			return nil
		}

		handled := false
		for i, item := range items {
			if choice == fmt.Sprintf("%d", i+1) {
				handled = true
				if err := item.action(); err != nil {
					fmt.Println("Error:", err)
				}
				break
			}
		}

		if !handled {
			fmt.Println("Invalid option. Try again.")
		}
	}
}

func main() {
	items := []menuItem{
		{
			label: "Show welcome message",
			action: func() error {
				fmt.Println("Hello from Ariseth Fuko Sol Module")
				return nil
			},
		},
		{
			label: "Run sample task",
			action: func() error {
				fmt.Println("Sample task executed.")
				return nil
			},
		},
	}

	if err := runMenu("Ariseth CLI", items); err != nil {
		fmt.Println("Fatal:", err)
		os.Exit(1)
	}
}
