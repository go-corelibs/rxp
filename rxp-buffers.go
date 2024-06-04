// Copyright (c) 2024  The Go-CoreLibs Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rxp

import (
	"strings"
	"sync"
)

// appendSlice is a better option than stock append for high-volume operations
func appendSlice[V interface{}](slice []V, data ...V) []V {
	m := len(slice)     // current length
	n := m + len(data)  // needed length
	if n > cap(slice) { // current cap size
		grown := make([]V, (m+1)*2) // double the existing space
		copy(grown, slice)          // populate grown
		slice = grown               // grown becomes slice
	}
	slice = slice[0:n]     // truncate in case too many present
	copy(slice[m:n], data) // populate with data
	return slice
}

var _ SyncPool[bool] = (*cSyncPool[bool])(nil)

type SyncPool[V interface{}] interface {
	Scale() int
	Ready() int
	Seed(count int)
	Get() V
	Put(v V)
}

type cSyncPool[V interface{}] struct {
	ready int // ready is an estimate
	scale int
	maker func() V
	pool  sync.Pool // pool is the underlying storage
}

func newSyncPool[V interface{}](scale int, maker func() V) SyncPool[V] {
	if scale <= 0 {
		scale = 2
	}
	p := &cSyncPool[V]{
		ready: 0,
		scale: scale,
		maker: maker,
	}
	p.pool = sync.Pool{New: func() interface{} {
		p.ready = 0
		return maker()
	}}
	p.Seed(scale)
	return p
}

func (p *cSyncPool[V]) Scale() int {
	return p.scale
}

func (p *cSyncPool[V]) Seed(count int) {
	if count <= 0 {
		count = p.scale
	}
	var total int
	if delta := count - p.ready; count > p.ready {
		for i := 0; i < delta; i++ {
			p.pool.Put(p.maker())
		}
		total += delta
	}
	p.ready += total
}

func (p *cSyncPool[V]) Ready() int {
	return p.ready
}

func (p *cSyncPool[V]) Get() V {
	ready := p.Ready()
	if ready <= 0 {
		p.ready = 0
		p.Seed(p.scale)
	}
	p.ready -= 1
	return p.pool.Get().(V)
}

func (p *cSyncPool[V]) Put(v V) {
	p.pool.Put(v)
	p.ready += 1
}

// strings.Builder memory pool

var (
	spStringsBuilder = newSyncPool(1, func() *strings.Builder {
		return new(strings.Builder)
	})
)

func getStringsBuilder() *strings.Builder {
	buf := spStringsBuilder.Get()
	buf.Reset()
	return buf
}

func putStringsBuilder(buf *strings.Builder) {
	if buf.Len() < 64000 {
		spStringsBuilder.Put(buf)
	}
}
