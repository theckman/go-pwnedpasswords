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
 
## License
This code is released under the MIT License. Please see the
[LICENSE](https://github.com/theckman/go-pwnedpasswords/blob/master/LICENSE) for the
full content of the license.

## Usage
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
