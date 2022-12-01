package keychain

import (
	"fmt"
	"os/exec"
)

const (
	addGenericPassword   = "add-generic-password"
	addInternetPassword  = "add-internet-password"
	findGenericPassword  = "find-generic-password"
	findInternetPassword = "find-internet-password"

	generic secretKind = iota + 1
	password
)

type secretKind int

type keychain struct {
	desc string
}

func New(desc string) *keychain {
	return &keychain{
		desc: desc,
	}
}

func (k *keychain) NewPassword(pw, name, service string) error {
	return add(password, pw, name, service, "", k.desc)
}

func (k *keychain) Password(name, service string) ([]byte, error) {
	return find(password, name, service, "")
}

func (k *keychain) NewSecret(kind, secret, name, service string) error {
	return add(generic, secret, name, service, kind, k.desc)
}

func (k *keychain) Secret(kind, name, service string) ([]byte, error) {
	return find(generic, name, service, kind)
}

func add(sk secretKind, secret, name, service, extKind, comment string) error {
	cargs := []string{}
	switch sk {
	case generic:
		cargs = append(cargs, addGenericPassword, "-D", extKind)
	case password:
		cargs = append(cargs, addInternetPassword, "-r", "htps")
	default:
		panic("invalid kind")
	}

	cargs = append(cargs, "-a", name, "-s", service, "-j", comment, "-w", secret)
	ec := exec.Command("security", cargs...)
	out, err := ec.CombinedOutput()
	if err != nil {
		msg := err.Error()
		if len(out) > 0 {
			msg = string(out)
		}
		return fmt.Errorf(msg)
	}
	return nil
}

func find(sk secretKind, name, service, extKind string) ([]byte, error) {
	cargs := []string{}
	switch sk {
	case generic:
		cargs = append(cargs, findGenericPassword, "-D", extKind)
	case password:
		cargs = append(cargs, findInternetPassword)
	default:
		panic("invalid kind")
	}

	cargs = append(cargs, "-a", name, "-s", service, "-w")
	ec := exec.Command("security", cargs...)
	out, err := ec.CombinedOutput()
	if err != nil {
		msg := err.Error()
		if len(out) > 0 {
			msg = string(out)
		}
		return nil, fmt.Errorf(msg)
	}

	return out[:len(out)-1], nil // cut newline
}
