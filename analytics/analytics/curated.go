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

package analytics

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetCuratedSchemas(rootFolder string) []string {
	curatedSchemasFilename := filepath.Join(rootFolder, "curated_schemas.txt")
	curatedSchemasBytes, err := os.ReadFile(curatedSchemasFilename)
	if err != nil {
		log.Fatal(err)
	}
	curated := strings.Split(string(curatedSchemasBytes), "\n")
	if len(curated) == 0 {
		log.Fatal("no curated schemas")
	}
	if curated[len(curated)-1] == "" {
		curated = curated[:len(curated)-1]
	}
	return curated
}
