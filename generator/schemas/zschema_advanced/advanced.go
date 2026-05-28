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
	. "github.com/katydid/validator-jsonschema-benchmarks/generator/rand/randjsonschema"
)

func Advanced() Rand {
	return Object(WithFields(
		Field("/", Entry(), IsRequired()),
		PatternField(FieldName(), Entry()),
	))
}

func Entry() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("storage", Storage(), IsRequired()),
		Field("fstype", Or(WithAnyOf(Const(`"ext3"`), Const(`"ext4"`), Const(`"btrfs"`)))),
		Field("options", ArrayOf(String(), WithMinItems(1), WithUniqueItems())),
		Field("readonly", Bool()),
	))
}

func Storage() Rand {
	return Or(WithAnyOf(
		DiskDevice(),
		DiskUUID(),
		NFS(),
		TmpFS(),
	))
}

func DiskDevice() Rand {
	return Object(WithFields(
		Field("type", Const(`"disk"`), IsRequired()),
		Field("device", DiskDevicePattern(), IsRequired()),
	))
}

func DiskUUID() Rand {
	return Object(WithFields(
		Field("type", Const(`"disk"`), IsRequired()),
		Field("label", DiskUUIDLabel(), IsRequired()),
	))
}

func NFS() Rand {
	return Object(WithFields(
		Field("type", Const(`"nfs"`), IsRequired()),
		Field("remotePath", NFSRemotePath(), IsRequired()),
		Field("server", NFSServer(), IsRequired()),
	))
}

func TmpFS() Rand {
	return Object(WithFields(
		Field("type", Const(`"tmpfs"`), IsRequired()),
		Field("sizeInMB", Integer(WithMinimum(16), WithMaximum(512)), IsRequired()),
	))
}

func AdvancedrmUniqueItems() Rand {
	return Object(WithFields(
		Field("/", EntryrmUniqueItems(), IsRequired()),
		PatternField(FieldName(), EntryrmUniqueItems()),
	))
}

func EntryrmUniqueItems() Rand {
	return Object(WithAdditionalFields(), WithFields(
		Field("storage", Storage(), IsRequired()),
		Field("fstype", Or(WithAnyOf(Const(`"ext3"`), Const(`"ext4"`), Const(`"btrfs"`)))),
		Field("options", ArrayOf(String(), WithMinItems(1))),
		Field("readonly", Bool()),
	))
}
