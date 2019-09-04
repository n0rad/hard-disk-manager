package password

import (
	"fmt"
	"github.com/awnumar/memguard"
	"github.com/n0rad/go-erlog/errs"
	"github.com/n0rad/go-erlog/logs"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	password *memguard.Enclave
	stop     chan struct{}
}

func (s *Service) Stop(e error) {
	close(s.stop)
}

func (s *Service) Start() error {
	s.stop = make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-term:
	case <-s.stop:
	}

	logs.Info("Purge memguard")
	memguard.Purge()
	return nil
}

func (s *Service) FromConnection(conn net.Conn) error {
	buf := memguard.NewBufferFromReaderUntil(conn, '\n')
	s.password = buf.Seal()
	return nil
}

func (s *Service) FromStdin(confirmation bool) error {
	var password, passwordConfirm []byte
	var err error

	defer memguard.WipeBytes(password)
	defer memguard.WipeBytes(passwordConfirm)

	for {
		print("Password: ")
		password, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return errs.WithE(err, "Cannot read password")
		}

		print("\n")
		if !confirmation {
			s.password = memguard.NewEnclave(password)
			return nil
		}

		print("Confirm: ")
		passwordConfirm, err = terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return errs.WithE(err, "Cannot read password")
		}
		print("\n")

		if string(password) == string(passwordConfirm) && string(password) != "" {
			s.password = memguard.NewEnclave(password)
			return nil
		} else {
			fmt.Println("\nEmpty password or do not match...\n")
		}
	}
}

func (s Service) Write(writer io.Writer) error {
	var total, written int
	var err error

	lockedBuffer, err := s.password.Open()
	if err != nil {
		return errs.WithE(err, "Failed to open password enclave")
	}
	defer lockedBuffer.Destroy()

	bytes := lockedBuffer.Bytes()

	for total = 0; total < len(bytes); total += written {
		written, err = writer.Write(bytes[total:])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) Get() (*memguard.LockedBuffer, error) {
	return s.password.Open()
}

/////
