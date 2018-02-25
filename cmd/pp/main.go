// Copyright (c) 2018 Tim Heckman
//
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file at the root of this repository.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/howeyc/gopass"

	"github.com/theckman/go-pwnedpasswords"
	"github.com/theckman/humanize-go"
)

var (
	help bool
)

func initFlags() {
	flag.BoolVar(&help, "h", false, "print help and exit")
	flag.Parse()
}

func main() {
	initFlags()

	if help {
		printHelp()
	}

	fmt.Print("enter password\n> ")

	pass, err := gopass.GetPasswdMasked()

	// some extra space in case anything weird happened with password input
	fmt.Println()

	if err != nil {
		log.Fatalf("failed to read password: %s", err)
	}

	pp, err := pwnedpasswords.New(pwnedpasswords.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	n, err := pp.Check(pass)
	if err != nil {
		log.Fatalf("failed to check for password compromise: %s", err)
	}

	if n > 0 {
		fmt.Print("\t!!!!!!!!!!!!!!!\n")
		fmt.Print("\t!! ATTENTION !!\n")
		fmt.Print("\t!!!!!!!!!!!!!!!\n\n")
		fmt.Printf("your password has been compromised at least %s times\n", humanize.CommaInt(n))
		os.Exit(1)
	}

	fmt.Println("no compromises detected")
}

func printHelp() {
	fmt.Printf("Usage of %s\n\n", os.Args[0])
	fmt.Print("PwnedPasswords asks for your password over stdin, and checks it against the PwnedPasswords API\n\n")
	fmt.Print("This works by SHA-1 hashing your password, and sending the first five bytes of the hex-encoded hash.\n")
	fmt.Print("The API returns a list of hashes that start with the same prefix, and we then compare hashes locally.\n\n")
	fmt.Print("This utility does not transmit your password or the full SHA-1 hash, making it safe to use.\n\n")
	fmt.Printf("Version: %s, Copyright 2018 Tim Heckman\n", pwnedpasswords.Version)
	os.Exit(0)
}
