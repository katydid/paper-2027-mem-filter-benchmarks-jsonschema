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

import "testing"

func TestMedian4(t *testing.T) {
	nums := []int{10, 20, 30, 40}
	got := MedianInt(4, func(i int) int { return nums[i] })
	want := 25
	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}

func TestMedian3(t *testing.T) {
	nums := []int{10, 20, 30}
	got := MedianInt(3, func(i int) int { return nums[i] })
	want := 20
	if got != want {
		t.Fatalf("got %d want %d", got, want)
	}
}
