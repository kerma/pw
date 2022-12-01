package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"

	cmd "github.com/google/subcommands"
	"github.com/kerma/pw"
	"github.com/kerma/pw/store"
	"github.com/kerma/pw/store/keychain"
)

func NewStore(desc string) store.SecretStore {
	if runtime.GOOS == "darwin" {
		return keychain.New(desc)
	}
	panic("platform not supported")
}

func main() {
	log.SetFlags(0)

	debug := flag.Bool("debug", false, "turn on debug logging")

	cmd.Register(cmd.FlagsCommand(), "")
	cmd.Register(new(pw.CopyCmd), "")
	cmd.Register(new(pw.CreateCmd), "")
	cmd.Register(new(pw.OTPCmd), "")

	flag.Parse()

	if *debug {
		log.SetFlags(log.Lshortfile)
	}

	store := NewStore("Created by pw cli")
	os.Exit(int(cmd.Execute(context.Background(), store)))
}
