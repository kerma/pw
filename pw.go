package pw

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/kerma/pw/store"
)

const (
	aliasConfigPath  = ".config/pw/alias.csv"
	clipboardTimeout = 30 * time.Second
	kindOTP          = "pw-otp-secret"
	kindSecret       = "pw-secret"
)

var (
	kinds = map[string]struct{}{
		"otp":      {},
		"password": {},
		"secret":   {},
	}
)

func load(alias string, name, service *string) error {
	if alias != "" {
		user, err := user.Current()
		if err != nil {
			return err
		}
		fp := path.Join(user.HomeDir, aliasConfigPath)
		f, err := os.Open(fp)
		if err != nil {
			return err
		}
		defer f.Close()

		found := false
		r := csv.NewReader(f)
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if rec[0] == alias && len(rec) == 3 {
				*name = rec[1]
				*service = rec[2]
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("alias not found")
		}
	}

	if *name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if *service == "" {
		return fmt.Errorf("service cannot be empty")
	}
	return nil
}

func parse(args []string) (string, string, error) {
	if len(args) < 2 {
		return "", "", fmt.Errorf("name and service arguments required")
	}

	var (
		name    string
		service string
	)
	if name = args[0]; name == "" {
		return "", "", fmt.Errorf("name argument required")
	}
	if service = args[1]; service == "" {
		return "", "", fmt.Errorf("service argument required")
	}
	return name, service, nil
}

func loadStore(args []any) store.SecretStore {
	if len(args) < 1 {
		log.Panic("store.SecretStore argument not provied")
	}
	s, ok := args[0].(store.SecretStore)
	if !ok {
		log.Panic("store.SecretStore argument not provied")
	}
	return s
}
