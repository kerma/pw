package pw

import (
	"context"
	"flag"
	"log"

	cmd "github.com/google/subcommands"
	"golang.design/x/clipboard"
)

type CopyCmd struct {
	alias string
	kind  string
}

func (*CopyCmd) Name() string     { return "copy" }
func (*CopyCmd) Synopsis() string { return "copies a secret to clipboard" }
func (*CopyCmd) Usage() string {
	return "Usage: copy [-a alias] <name> <service>\n"
}

func (c *CopyCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.alias, "a", "", "alias for name/service pair")
	f.StringVar(&c.kind, "k", "password", "password or secret")
}

func (c *CopyCmd) Execute(_ context.Context, f *flag.FlagSet, args ...any) cmd.ExitStatus {
	var (
		err     error
		name    string
		service string
	)
	switch c.alias {
	case "":
		name, service, err = parse(f.Args())
		if err != nil {
			log.Println(err)
			return cmd.ExitFailure
		}
	default:
		if err := load(c.alias, &name, &service); err != nil {
			log.Println(err)
			return cmd.ExitFailure
		}
	}

	if _, ok := kinds[c.kind]; !ok {
		log.Println("invalid kind")
		return cmd.ExitFailure
	}

	var secret []byte
	switch c.kind {
	case "password":
		secret, err = loadStore(args).Password(name, service)
	default:
		secret, err = loadStore(args).Secret(c.kind, name, service)
	}

	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	err = clipboard.Init()
	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	clipboard.Write(clipboard.FmtText, secret)

	log.Printf("📋 %s copied to clipboard\n", c.kind)
	return cmd.ExitSuccess
}
