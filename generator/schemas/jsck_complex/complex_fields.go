// Copyright 2026 Walter Schulze
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsck_complex

import (
	"strconv"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjson"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

// "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$"
func Base58() randjsonschema.Rand {
	return &randBase58{}
}

type randBase58 struct{}

func (o *randBase58) Right(r rand.Rand) string {
	l := r.Intn(99) + 1
	ss := make([]rune, l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(26))
			if ss[i] == 'l' {
				// this letter is not allowed
				ss[i] = 'a'
			}
		case 1:
			ss[i] = 'A' + rune(r.Intn(26))
			if ss[i] == 'I' || ss[i] == 'O' {
				// these two letters are not allowed
				ss[i] = 'A'
			}
		case 2:
			ss[i] = '1' + rune(r.Intn(9))
		}
	}
	return strconv.Quote(string(ss))
}

func (o *randBase58) Wrong(r rand.Rand) string {
	l := r.Intn(99) + 1
	ss := make([]rune, l)
	wrongIndex := r.Intn(l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(26))
		case 1:
			ss[i] = 'A' + rune(r.Intn(26))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
		if wrongIndex == i {
			ss[i] = '"'
		}
	}
	return strconv.Quote(string(ss))
}

// "^[0123456789A-Fa-f]+$"
func Hex() randjsonschema.Rand {
	return &randHex{}
}

type randHex struct{}

func (o *randHex) Right(r rand.Rand) string {
	l := r.Intn(99) + 1
	ss := make([]rune, l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(6))
		case 1:
			ss[i] = 'A' + rune(r.Intn(6))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
	}
	return strconv.Quote(string(ss))
}

func (o *randHex) Wrong(r rand.Rand) string {
	l := r.Intn(99) + 1
	ss := make([]rune, l)
	wrongIndex := r.Intn(l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(6))
		case 1:
			ss[i] = 'A' + rune(r.Intn(6))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
		if wrongIndex == i {
			ss[i] = '"'
		}
	}
	return strconv.Quote(string(ss))
}

// "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]+$"
// length = 34
func Address() randjsonschema.Rand {
	return &randAddress{}
}

type randAddress struct{}

func (o *randAddress) Right(r rand.Rand) string {
	l := 34
	ss := make([]rune, l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(26))
			if ss[i] == 'l' {
				// this letter is not allowed
				ss[i] = 'a'
			}
		case 1:
			ss[i] = 'A' + rune(r.Intn(26))
			if ss[i] == 'I' || ss[i] == 'O' {
				// these two letters are not allowed
				ss[i] = 'A'
			}
		case 2:
			ss[i] = '1' + rune(r.Intn(9))
		}
	}
	return strconv.Quote(string(ss))
}

func (o *randAddress) Wrong(r rand.Rand) string {
	l := 34 + r.Intn(5) - 2 // maybe vary the length
	ss := make([]rune, l)
	wrongIndex := r.Intn(l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(26))
		case 1:
			ss[i] = 'A' + rune(r.Intn(26))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
		if wrongIndex == i {
			ss[i] = '"'
		}
	}
	return strconv.Quote(string(ss))
}

// "^[0123456789A-Fa-f]+$"
// length = 64
func TxID() randjsonschema.Rand {
	return &randTxID{}
}

type randTxID struct{}

func (o *randTxID) Right(r rand.Rand) string {
	l := 64
	ss := make([]rune, l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(6))
		case 1:
			ss[i] = 'A' + rune(r.Intn(6))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
	}
	return strconv.Quote(string(ss))
}

func (o *randTxID) Wrong(r rand.Rand) string {
	l := 64 + r.Intn(5) - 2 // maybe vary the length
	ss := make([]rune, l)
	wrongIndex := r.Intn(l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(6))
		case 1:
			ss[i] = 'A' + rune(r.Intn(6))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
		if wrongIndex == i {
			ss[i] = '"'
		}
	}
	return strconv.Quote(string(ss))
}

// "^[0123456789A-Fa-f]+$"
// length = 128
func Signature() randjsonschema.Rand {
	return &randSignature{}
}

type randSignature struct{}

func (o *randSignature) Right(r rand.Rand) string {
	l := 128
	ss := make([]rune, l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(6))
		case 1:
			ss[i] = 'A' + rune(r.Intn(6))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
	}
	return strconv.Quote(string(ss))
}

func (o *randSignature) Wrong(r rand.Rand) string {
	l := 128 + r.Intn(5) - 2 // maybe vary the length
	ss := make([]rune, l)
	wrongIndex := r.Intn(l)
	for i := range l {
		switch r.Intn(3) {
		case 0:
			ss[i] = 'a' + rune(r.Intn(6))
		case 1:
			ss[i] = 'A' + rune(r.Intn(6))
		case 2:
			ss[i] = '0' + rune(r.Intn(10))
		}
		if wrongIndex == i {
			ss[i] = '"'
		}
	}
	return strconv.Quote(string(ss))
}

func ScriptType() randjsonschema.Rand {
	return &randScriptType{}
}

type randScriptType struct{}

func (o *randScriptType) Right(r rand.Rand) string {
	switch r.Intn(2) {
	case 0:
		return `"standard"`
	case 1:
		return `"p2sh"`
	}
	panic("unreachable")
}

func (o *randScriptType) Wrong(r rand.Rand) string {
	return randjson.String(r)
}

func Status() randjsonschema.Rand {
	return &randStatus{}
}

type randStatus struct{}

func (o *randStatus) Right(r rand.Rand) string {
	switch r.Intn(4) {
	case 0:
		return `"unsigned"`
	case 1:
		return `"unconfirmed"`
	case 2:
		return `"confirmed"`
	case 3:
		return `"invalid"`
	}
	panic("unreachable")
}

func (o *randStatus) Wrong(r rand.Rand) string {
	return randjson.String(r)
}
