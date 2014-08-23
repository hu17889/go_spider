// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package mcounter

type Mcounter interface {
    Incr()
    Decr()
    Count() int
}
