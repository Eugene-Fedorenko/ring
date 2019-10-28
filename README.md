# ring - high performance bloom filter

Package ring provides a high performance and thread safe Go implementation of a
bloom filter.

## Usage
Please see the [godoc](https://godoc.org/github.com/TheTannerRyan/ring) for
usage.

## Accuracy

Difference

**Original library**

```
goos: linux
goarch: amd64
pkg: github.com/thetannerryan/ring
BenchmarkGenerateMultiHash              25803452                44.6 ns/op             5 B/op          1 allocs/op
BenchmarkGenerateMultiHash-2            27672391                43.3 ns/op             5 B/op          1 allocs/op
BenchmarkGenerateMultiHash-4            26718860                43.4 ns/op             5 B/op          1 allocs/op
BenchmarkGenerateMultiHash-8            26856033                43.7 ns/op             5 B/op          1 allocs/op
BenchmarkGenerateMultiHash-16           26483736                44.0 ns/op             5 B/op          1 allocs/op
BenchmarkAddConcurrent                   1715552               700 ns/op               5 B/op          1 allocs/op
BenchmarkAddConcurrent-2                 1586790               759 ns/op               5 B/op          1 allocs/op
BenchmarkAddConcurrent-4                 1492144               786 ns/op               5 B/op          1 allocs/op
BenchmarkAddConcurrent-8                 1540778               762 ns/op               5 B/op          1 allocs/op
BenchmarkAddConcurrent-16                1570785               746 ns/op               5 B/op          1 allocs/op
BenchmarkTestConcurrent                  1736050               691 ns/op               5 B/op          1 allocs/op
BenchmarkTestConcurrent-2                3113769               389 ns/op               5 B/op          1 allocs/op
BenchmarkTestConcurrent-4                3564541               321 ns/op               5 B/op          1 allocs/op
BenchmarkTestConcurrent-8                3748862               322 ns/op               5 B/op          1 allocs/op
BenchmarkTestConcurrent-16               3725818               327 ns/op               5 B/op          1 allocs/op

=== RUN   TestGenerateMultiHash
--- PASS: TestGenerateMultiHash (0.00s)
=== RUN   TestConcurrentErrors
--- PASS: TestConcurrentErrors (174.29s)
    ring_test.go:126: FNR: 0.000000, errors: 0
    ring_test.go:127: FPR: 0.010232, errors: 10232
```

**With my changes**
```
goos: linux
goarch: amd64
pkg: github.com/eugene-fedorenko/ring
BenchmarkGenerateMultiHash              21212788                55.2 ns/op             0 B/op          0 allocs/op
BenchmarkGenerateMultiHash-2            22217754                52.1 ns/op             0 B/op          0 allocs/op
BenchmarkGenerateMultiHash-4            22598011                52.0 ns/op             0 B/op          0 allocs/op
BenchmarkGenerateMultiHash-8            19988134                52.5 ns/op             0 B/op          0 allocs/op
BenchmarkGenerateMultiHash-16           21819634                52.0 ns/op             0 B/op          0 allocs/op
BenchmarkAddConcurrent                   1665565               699 ns/op               0 B/op          0 allocs/op
BenchmarkAddConcurrent-2                 2890574               400 ns/op               0 B/op          0 allocs/op
BenchmarkAddConcurrent-4                 3272278               361 ns/op               0 B/op          0 allocs/op
BenchmarkAddConcurrent-8                 3295758               360 ns/op               0 B/op          0 allocs/op
BenchmarkAddConcurrent-16                3335684               363 ns/op               0 B/op          0 allocs/op
BenchmarkTestConcurrent                  1697246               701 ns/op               0 B/op          0 allocs/op
BenchmarkTestConcurrent-2                3122142               391 ns/op               0 B/op          0 allocs/op
BenchmarkTestConcurrent-4                3544603               341 ns/op               0 B/op          0 allocs/op
BenchmarkTestConcurrent-8                3705194               325 ns/op               0 B/op          0 allocs/op
BenchmarkTestConcurrent-16               3586874               324 ns/op               0 B/op          0 allocs/op

=== RUN   TestGenerateMultiHash
--- PASS: TestGenerateMultiHash (0.00s)
=== RUN   TestConcurrentErrors
--- PASS: TestConcurrentErrors (185.35s)
    ring_test.go:126: FNR: 0.000000, errors: 0
    ring_test.go:127: FPR: 0.009848, errors: 9848
```

## License
Copyright (c) 2019 Tanner Ryan. All rights reserved. Use of this source code is
governed by a BSD-style license that can be found in the LICENSE file.
