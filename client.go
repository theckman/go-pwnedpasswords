// Copyright (c) 2018 Tim Heckman
//
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file at the root of this repository.

// Package pwnedpasswords implements a client for checking passwords against the
// "Have I Been Pwned", Pwned Passwords API. The Pwned Passwords API implements
// a k-Anonymity model that allows you to check your password against the
// database without providing the API the full password or full SHA-1 password
// hash.
//
// This works by creating a SHA-1 hash of the password locally, hex-encodes the
// SHA-1 checksum, and then sends the first five bytes (prefix) to the Pwned
// Passwords API. The API then returns the suffix of hashes it has that start
// with that prefix. The client then compares the returned hashes locally to
// look for a match. This prevents the password, hashed or not, from leaving the
// local system.
package pwnedpasswords

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type hashCount struct {
	suffix string
	count  int
}

// DefaultURL is the default URL to the Pwned Passwords API.
const DefaultURL = "https://api.pwnedpasswords.com/range/"

// Version is the package version.
const Version = "1.0.0"

const userAgent = "go-pwnedpasswords/" + Version + " (https://github.com/theckman/go-pwnedpasswords) Go-http-client/1.1"

// Client is for checking passwords against the Pwned Passwords API without
// leaking the password.
type Client struct {
	h string
	c *http.Client
}

// New returns a new Client for checking passwords against the API. The urlStr
// argument should be the full path to the API endpoint, including the trailing
// slash. A good default is pwnedpasswords.DefaultURL.
func New(urlStr string) (Client, error) {
	_, err := url.Parse(urlStr)
	if err != nil {
		return Client{}, fmt.Errorf("failed to parse %q: %s", urlStr, err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          20,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}

	return Client{
		h: urlStr,
		c: httpClient,
	}, nil
}

// Check returns the number of times the password appears in PwnedPassowrds, and
// any errors that occur. If the value of the int is 0, your password is clean.
// If the value is greater than 0, change your password!
func (c Client) Check(password []byte) (int, error) {
	prefix, suffix := HashPassword(password)

	hashes, err := getHashes(c, prefix)
	if err != nil {
		return -1, fmt.Errorf("failed to get hashes: %s", err)
	}

	for _, hash := range hashes {
		if hash.suffix == suffix {
			return hash.count, nil
		}
	}

	return 0, nil
}

func getHashes(c Client, prefix string) ([]hashCount, error) {
	req, err := http.NewRequest(http.MethodGet, c.h+prefix, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request to %q failed: %s", req.URL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected http status code: %s", resp.Status)
	}

	return readHashes(resp.Body)
}

func readHashes(r io.Reader) ([]hashCount, error) {
	var hashes []hashCount

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")

		if len(parts) == 2 {
			i, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			hashes = append(hashes, hashCount{suffix: parts[0], count: i})
		}
	}

	if len(hashes) == 0 {
		return nil, errors.New("no hashes in response")
	}

	return hashes, nil
}

// HashPassword takes a password, returns the SHA-1 hash split in to the prefix
// and suffix. The prefix is what's used by the API, and the suffix should then
// be used to match returned results.
//
// Note: the full hash should *NEVER* be written to disk or sent across the
// network. If the value makes its way somewhere, it could be used to crack the
// password. You should only transmit the prefix to the PwnedPasswords API.
func HashPassword(password []byte) (prefix, suffix string) {
	h := sha1.New()

	h.Write(password)

	hash := fmt.Sprintf("%X", h.Sum(nil))

	return hash[0:5], hash[5:]
}
