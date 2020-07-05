
# stdrouter (Standard Library Router Generator)

<p style="text-align: center">
<img src="https://user-images.githubusercontent.com/38237246/86530951-3515a800-bef8-11ea-8168-756b86a3b661.png" alt="Router">
</p>

<p style="text-align: center">
<a href="LICENSE"><img src="http://img.shields.io/badge/license-MIT-blue.svg?style=flat" alt="MIT License"></a
<img src="https://github.com/tetsuzawa/go-adflib/workflows/Test/badge.svg" alt="Test">
</p>

## Introduction

stdrouter generates http router that can handle __path parameter__ written only with Go standard library.

## Features

- No external library (only go standard library)
- Simple implementation
- Easy to use middleware


## Usage

1. Create routing configuration file as `router.go`
2. Run `stdrouter` in the same directory as `router.go`
3. `router_gen.go` will be created. This is the implementation of router.


See [example](_example) for detail.

## Installation

```shell script
$ go get -u github.com/tetsuzawa/stdrouter/cmd/stdrouter
```

## License

[MIT](LICENSE) License

Copyright (c) 2020-present, Tetsu Takizawa and Contributors.
