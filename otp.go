package pw

import (
	"context"
	"flag"
	"log"

	cmd "github.com/google/subcommands"
	"github.com/xlzd/gotp"
	"golang.design/x/clipboard"
)

type OTPCmd struct {
	alias string
}

func (*OTPCmd) Name() string     { return "otp" }
func (*OTPCmd) Synopsis() string { return "copy otp code to clipboard" }
func (*OTPCmd) Usage() string {
	return "Usage: otp [-a alias] <name> <service>\n"
}

func (c *OTPCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.alias, "a", "", "alias for name/server pair")
}

func (c *OTPCmd) Execute(_ context.Context, f *flag.FlagSet, args ...any) cmd.ExitStatus {
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

	secret, err := loadStore(args).Secret(kindOTP, name, service)
	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	err = clipboard.Init()
	if err != nil {
		log.Println(err)
		return cmd.ExitFailure
	}

	code := gotp.NewDefaultTOTP(string(secret)).Now()
	clipboard.Write(clipboard.FmtText, []byte(code))

	log.Printf("📋 otp code %s copied to clipboard\n", code)
	return cmd.ExitSuccess
}
