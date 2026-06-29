// Copyright 2026 The Sieve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestSieve(t *testing.T) {
	for i := range samples {
		result := TestMode(samples[i][:1024])
		if result > .1 {
			t.Fatal("result should be less than .1")
		}
	}
}
