// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	computation := flag.String("time", "500ms", "Minimum computation time for the password hash. Based on this, the difficulty will be computed. Must be parseable by go time package")
	difficulty := flag.Int("difficulty", 0, "If set to non-zero, this difficulty will be used instead of calculating it based on execution time")
	password := flag.String("password", "", "Password. If empty, it will be read from stdin")

	flag.Parse()

	for *password == "" {
		fmt.Printf("Enter Password: ")
		pw1, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Print("\n")
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(pw1) == 0 {
			fmt.Println("No password, exiting")
			return
		}

		fmt.Printf("Repeat Password: ")
		pw2, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Print("\n")
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(pw2) == 0 {
			fmt.Println("No password, exiting")
			return
		}

		if bytes.Compare(pw1, pw2) != 0 {
			fmt.Println("Passwords do not match!")
			continue
		}

		*password = string(pw1)
	}

	var hash []byte

	if *difficulty != 0 {
		var err error
		hash, err = bcrypt.GenerateFromPassword([]byte(*password), *difficulty)
		if err != nil {
			log.Panicf("error computing bcrypt (difficulty %d): %s", *difficulty, err.Error())
		}
	} else {
		d, err := time.ParseDuration(*computation)
		found := false
		if err != nil {
			log.Panicf("Can not parse duration '%s': %s", *computation, err.Error())
		}
		for i := bcrypt.MinCost; i <= bcrypt.MaxCost; i++ {
			start := time.Now()
			hash, err = bcrypt.GenerateFromPassword([]byte(*password), i)
			end := time.Now()

			if err != nil {
				log.Panicf("error computing bcrypt (difficulty %d based on time): %s", *difficulty, err.Error())
			}

			computationTime := end.Sub(start)
			if computationTime > d {
				fmt.Printf("Using difficulty %d with time %s\n", i, computationTime.String())
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Could not find a difficulty enough for %s - using highest difficulty %d", d.String(), bcrypt.MaxCost)
		}

		if hash == nil {
			log.Panic("Can not computate hash for unknown reason")
		}
	}

	fmt.Printf("Hash: %s\n", base64.StdEncoding.EncodeToString(hash))

}
