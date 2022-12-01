package pw

import (
	"context"
	"flag"
	"fmt"
	"log"

	cmd "github.com/google/subcommands"
	"golang.design/x/clipboard"
)

type CreateCmd struct {
	kind string
}

func (*CreateCmd) Name() string     { return "create" }
func (*CreateCmd) Synopsis() string { return "create a new entry" }
func (*CreateCmd) Usage() string {
	return "Usage: create [-k kind] <name> <service>\n\n"
}

func (c *CreateCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.kind, "k", "password", "password, otp or secret")
}

func (c *CreateCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...any) cmd.ExitStatus {

	name, service, err := parse(f.Args())
	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	if _, ok := kinds[c.kind]; !ok {
		log.Println("invalid kind")
		return cmd.ExitFailure
	}

	err = clipboard.Init()
	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	ctx, cancel := context.WithTimeout(ctx, clipboardTimeout)
	defer cancel()

	fmt.Println("waiting for clipboard...")
	ch := clipboard.Watch(ctx, clipboard.FmtText)

	var b []byte
	select {
	case <-ctx.Done():
		break
	case b = <-ch:
		break
	}

	if len(b) == 0 {
		fmt.Println("nothing found on clipboard")
		return cmd.ExitFailure
	}

	switch c.kind {
	case "otp":
		err = loadStore(args).NewSecret(kindOTP, string(b), name, service)
	case "password":
		err = loadStore(args).NewPassword(string(b), name, service)
	case "secret":
		err = loadStore(args).NewSecret(kindSecret, string(b), name, service)
	default:
		log.Panic("invalid kind")
	}

	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	log.Println("secret stored")
	return cmd.ExitSuccess
}
