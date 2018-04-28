# go-pwnedpasswords
[![License](https://img.shields.io/github/license/theckman/go-pwnedpasswords.svg)](https://github.com/theckman/go-pwnedpasswords/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/theckman/go-pwnedpasswords)
[![Latest Git Tag](https://img.shields.io/github/tag/theckman/go-pwnedpasswords.svg)](https://github.com/theckman/go-pwnedpasswords/releases)
[![Travis master Build Status](https://img.shields.io/travis/theckman/go-pwnedpasswords/master.svg?label=TravisCI)](https://travis-ci.org/theckman/go-pwnedpasswords/branches)
[![Go Cover Test Coverage](https://gocover.io/_badge/github.com/theckman/go-pwnedpasswords?v0)](https://gocover.io/github.com/theckman/go-pwnedpasswords)
[![Go Report Card](https://goreportcard.com/badge/github.com/theckman/go-pwnedpasswords)](https://goreportcard.com/report/github.com/theckman/go-pwnedpasswords)

 Package pwnedpasswords implements a client for checking passwords against the
 "Have I Been Pwned", Pwned Passwords API. The Pwned Passwords API implements a
 k-Anonymity model that allows you to check your password against the database
 without providing the API the full password or full SHA-1 password hash.

 This works by creating a SHA-1 hash of the password locally, hex-encodes the
 SHA-1 checksum, and then sends the first five bytes (prefix) to the Pwned
 Passwords API. The API then returns the suffix of hashes it has that start with
 that prefix. The client then compares the returned hashes locally to look for a
 match. This prevents the password, hashed or not, from leaving the local
 system.
 
 In effect, this allows you to confirm if your password is seen in one of the
 many database dumps where passwords were obtained. If your password is not on
 the list, it does not mean that it is safe or hasn't been compromised. Always
 remember to never share passwords between different sites or services, as the
 compromise of one can lead to the compromise of all of your accounts.
 
## License
This code is released under the MIT License. Please see the
[LICENSE](https://github.com/theckman/go-pwnedpasswords/blob/master/LICENSE) for the
full content of the license.

## Building the Binary
If you have the Go toolchain installed, you can use the following command to
install the pwnedpasswords command line client (`pp`):

```Shell
go get github.com/theckman/go-pwnedpasswords/cmd/pp
```

## Usage
If you plan to use this package as a client library in Go, here is a quick
example of how to use it:

```Go
client, err := pwnedpasswords.New(pwnedpasswords.DefaultURL)
// handle error

compromiseCount, err := client.Check([]byte("password"))
// handle error

// password was compromised on at least compromiseCount sites
if compromiseCount > 0 {
	// handle situation where password is compromised
	// in other words, never using it ever again...
}

// password may not be compromised
```
