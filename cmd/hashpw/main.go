package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	var pw string
	if len(os.Args) > 1 {
		pw = os.Args[1]
	} else {
		fmt.Print("password: ")
		s := bufio.NewScanner(os.Stdin)
		if s.Scan() {
			pw = strings.TrimSpace(s.Text())
		}
	}
	if pw == "" {
		fmt.Fprintln(os.Stderr, "empty password")
		os.Exit(1)
	}
	h, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(h))
}
