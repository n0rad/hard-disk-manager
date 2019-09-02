package utils

import (
	"fmt"
	"github.com/n0rad/go-erlog/errs"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func AskPasswordWithConfirmation(confirmation bool) ([]byte, error) {
	for {
		print("Password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, errs.WithE(err, "Cannot read password")
		}

		print("\n")
		if !confirmation {
			return password, nil
		}

		print("Confirm: ")
		passwordConfirm, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, errs.WithE(err, "Cannot read password")
		}
		print("\n")

		if string(password) == string(passwordConfirm) && string(password) != "" {
			return password, nil
		} else {
			fmt.Println("\nEmpty password or do not match...\n")
		}
	}
}
