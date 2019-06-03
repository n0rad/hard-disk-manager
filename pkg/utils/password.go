package utils

import (
	"fmt"
	"github.com/n0rad/go-erlog/errs"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func AskPasswordWithConfirmation(confirmation bool) (string, error) {
	for {
		print("Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", errs.WithE(err, "Cannot read password")
		}

		print("\n")
		if !confirmation {
			return string(bytePassword), nil
		}

		print("Confirm: ")
		bytePassword2, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", errs.WithE(err, "Cannot read password")
		}
		print("\n")

		if string(bytePassword) == string(bytePassword2) && string(bytePassword) != "" {
			return string(bytePassword), nil
		} else {
			fmt.Println("\nEmpty password or do not match...\n")
		}
	}
}