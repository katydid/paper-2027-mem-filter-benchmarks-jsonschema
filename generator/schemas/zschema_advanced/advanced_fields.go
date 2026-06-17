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

package zschema_advanced

import (
	"fmt"
	"strings"

	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjson"
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func randPath(r rand.Rand, sep string) string {
	s := sep
	for {
		for range r.Intn(10) + 1 {
			s += string([]byte{byte('a' + r.Intn(25))})
		}
		if r.Intn(2) == 0 {
			return s
		}
		s += sep
	}
}

//	"patternProperties": {
//	    "^(/[^/]+)+$": {"$ref": "#/definitions/entry"}
//	},
//
// example: "/var", "/tmp", "/var/www"
func FieldName() Rand {
	return &randFieldName{}
}

type randFieldName struct{}

func (o *randFieldName) Right(r rand.Rand) string {
	switch r.Intn(3) {
	case 0:
		return `"/var"`
	case 1:
		return `"/tmp"`
	case 2:
		return `"/var/www"`
	case 3:
		return `"` + randPath(r, "/") + `"`
	}
	panic("unreachable")
}

func (o *randFieldName) Wrong(r rand.Rand) string {
	return strings.Replace(randjson.String(r), "/", "@", -1)
}

//	"device": {
//		"type": "string",
//		"pattern": "^/dev/[^/]+(/[^/]+)*$"
//	}
//
// example: "/dev/sda1"
func DiskDevicePattern() Rand {
	return &randDevicePattern{}
}

type randDevicePattern struct{}

func (o *randDevicePattern) Right(r rand.Rand) string {
	return `"/dev` + randPath(r, "/") + `"`
}

func (o *randDevicePattern) Wrong(r rand.Rand) string {
	return strings.Replace(randjson.String(r), "/", "@", -1)
}

//	"label": {
//		"type": "string",
//		"pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
//	}
//
// example: "8f3ba6f4-5c70-46ec-83af-0d5434953e5f"
func DiskUUIDLabel() Rand {
	return &randDiskUUIDLabel{}
}

type randDiskUUIDLabel struct{}

func randHex(r rand.Rand, n int) string {
	bs := make([]byte, n)
	for i := range n {
		switch r.Intn(3) {
		case 0:
			bs[i] = '0' + byte(r.Intn(10))
		case 1:
			bs[i] = 'a' + byte(r.Intn(6))
		case 2:
			bs[i] = 'A' + byte(r.Intn(6))
		}
	}
	return string(bs)
}

func (o *randDiskUUIDLabel) Right(r rand.Rand) string {
	return `"` + randHex(r, 8) + "-" + randHex(r, 4) + "-" + randHex(r, 4) + "-" + randHex(r, 4) + "-" + randHex(r, 12) + `"`
}

func (o *randDiskUUIDLabel) Wrong(r rand.Rand) string {
	return randjson.String(r)
}

//	"remotePath": {
//		"type": "string",
//		"pattern": "^(/[^/]+)+$"
//	},
//
// example: "/exports/mypath"
func NFSRemotePath() Rand {
	return &randNFSRemotePath{}
}

type randNFSRemotePath struct{}

func (o *randNFSRemotePath) Right(r rand.Rand) string {
	return `"` + randPath(r, "/") + `"`
}

func (o *randNFSRemotePath) Wrong(r rand.Rand) string {
	return strings.Replace(randjson.String(r), "/", "@", -1)
}

//	"server": {
//		"type": "string",
//		"anyOf": [{"format": "hostname"}, {"format": "ipv4"}, {"format": "ipv6"}]
//	}
//
// example: "my.nfs.server"
func NFSServer() Rand {
	return &randNFSServer{}
}

type randNFSServer struct{}

func (o *randNFSServer) Right(r rand.Rand) string {
	switch r.Intn(3) {
	case 0:
		// hostname
		return `"` + randPath(r, ".")[1:] + randPath(r, ".") + `"`
	case 1:
		// ipv4
		return IPv4().Right(r)
	case 2:
		// ipv6
		return randomIPv6(r)
	}
	panic("unreachable")
}

func (o *randNFSServer) Wrong(r rand.Rand) string {
	switch r.Intn(2) {
	case 0:
		// hostname / ipv4
		s := `"` + randPath(r, ".")[1:] + randPath(r, ".") + `"`
		wrong := strings.Replace(s, ".", "@", -1)
		return wrong
	case 1:
		// ipv6
		s := randomIPv6(r)
		return strings.Replace(s, ":", "@", 1)
	}
	panic("unreachable")
}

// https://generate-random.org/ip-addresses/go
func randomIPv6(r rand.Rand) string {
	parts := make([]string, 8)
	for i := 0; i < 8; i++ {
		parts[i] = fmt.Sprintf("%04x", r.Intn(0x10000))
	}
	return fmt.Sprintf("\"%s:%s:%s:%s:%s:%s:%s:%s\"",
		parts[0], parts[1], parts[2], parts[3],
		parts[4], parts[5], parts[6], parts[7])
}
