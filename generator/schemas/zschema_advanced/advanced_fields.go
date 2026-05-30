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
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand"
	"github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjson"
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

//	"patternProperties": {
//	    "^(/[^/]+)+$": {"$ref": "#/definitions/entry"}
//	},
//
// example: "/var", "/tmp", "/var/www"
func FieldName() Rand {
	return &randFieldName{}
}

type randFieldName struct{}

// TODO: generate more options
func (o *randFieldName) Right(r rand.Rand) string {
	switch r.Intn(3) {
	case 0:
		return `"/var"`
	case 1:
		return `"/tmp"`
	case 2:
		return `"/var/www"`
	}
	panic("unreachable")
}

func (o *randFieldName) Wrong(r rand.Rand) string {
	return randjson.String(r)
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

// TODO: generate more options
func (o *randDevicePattern) Right(r rand.Rand) string {
	return `"/dev/sda1"`
}

func (o *randDevicePattern) Wrong(r rand.Rand) string {
	return randjson.String(r)
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

// TODO: generate more options
func (o *randDiskUUIDLabel) Right(r rand.Rand) string {
	return `"8f3ba6f4-5c70-46ec-83af-0d5434953e5f"`
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

// TODO: generate more options
func (o *randNFSRemotePath) Right(r rand.Rand) string {
	return `"/exports/mypath"`
}

func (o *randNFSRemotePath) Wrong(r rand.Rand) string {
	return randjson.String(r)
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

// TODO: generate more options
func (o *randNFSServer) Right(r rand.Rand) string {
	return `"my.nfs.server"`
}

// TODO: generate more options
func (o *randNFSServer) Wrong(r rand.Rand) string {
	return `"123.123.123.45@"`
}
