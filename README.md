# pdfcpu: a Go PDF processor

[![Build Status](https://travis-ci.org/hhrutter/pdfcpu.svg?branch=master)](https://travis-ci.org/hhrutter/pdfcpu)
[![GoDoc](https://godoc.org/github.com/denisbetsi/pdfcpu?status.svg)](https://godoc.org/github.com/denisbetsi/pdfcpu)
[![Coverage Status](https://coveralls.io/repos/github/hhrutter/pdfcpu/badge.svg?branch=master)](https://coveralls.io/github/hhrutter/pdfcpu?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/denisbetsi/pdfcpu)](https://goreportcard.com/report/github.com/denisbetsi/pdfcpu)
[![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest release](https://img.shields.io/github/release/hhrutter/pdfcpu.svg)](https://github.com/denisbetsi/pdfcpu/releases)

<img src="resources/logoSmall.png" width="150">

pdfcpu is a simple PDF processing library written in [Go](http://golang.org) supporting encryption.
It provides both an API and a CLI. Supported are all versions up to PDF 1.7 (ISO-32000).

## Motivation

This is an effort to build a comprehensive PDF processing library from the ground up written in Go. Over time pdfcpu aims to support the standard range of PDF processing features and also any interesting use cases that may present themselves along the way.

<p align="center">
  <kbd><a href="https://pdfcpu.io/core/grid"><img src="resources/gridpdf.png" height="150"></a></kbd>&nbsp;
  <kbd><a href="https://pdfcpu.io/core/watermark"><img src="resources/wmi1abs.png" height="150"></a></kbd>&nbsp;
  <kbd><a href="https://pdfcpu.io/core/nup"><img src="resources/nup9pdf.png" height="150"></a></kbd>&nbsp;
  <kbd><a href="https://pdfcpu.io/core/stamp"><img src="resources/stt31.png" height="150"></a></kbd><br><br>
  <kbd><a href="https://pdfcpu.io/core/stamp"><img src="resources/sti.png" height="150"></a></kbd>&nbsp;
  <kbd><img src="resources/hold3.png" height="150"></kbd>&nbsp;
  <kbd><a href="https://pdfcpu.io/core/watermark"><img src="resources/wmi4.png" height="150"></a></kbd>&nbsp;<br><br>
  <kbd><a href="https://pdfcpu.io/core/stamp"><img src="resources/stp.png" height="150"></a></kbd>&nbsp;
  <kbd><a href="https://pdfcpu.io/core/grid"><img src="resources/gridimg.png" height="150"></a></kbd>
</p>

## Focus

The main focus lies on strong support for batch processing and scripting via a rich command line. At the same time pdfcpu wants to make it easy to integrate PDF processing into your Go based backend system by providing a robust command set.

## Command Set

* [attachments](https://pdfcpu.io/attach/attach)
* [change owner password](https://pdfcpu.io/encrypt/change_opw)
* [change user password](https://pdfcpu.io/encrypt/change_upw)
* [decrypt](https://pdfcpu.io/encrypt/decryptPDF)
* [encrypt](https://pdfcpu.io/encrypt/encryptPDF)
* [extract](https://pdfcpu.io/extract/extract)
* [grid](https://pdfcpu.io/core/grid)
* [import](https://pdfcpu.io/generate/import)
* [info](https://pdfcpu.io/info)
* [merge](https://pdfcpu.io/core/merge)
* [nup](https://pdfcpu.io/core/nup)
* [optimize](https://pdfcpu.io/core/optimize)
* [pages](https://pdfcpu.io/pages/pages)
* [permissions](https://pdfcpu.io/encrypt/perm_add)
* [split](https://pdfcpu.io/core/split)
* [rotate](https://pdfcpu.io/core/rotate)
* [stamp](https://pdfcpu.io/core/stamp)
* [trim](https://pdfcpu.io/core/trim)
* [validate](https://pdfcpu.io/core/validate)
* [watermark](https://pdfcpu.io/core/watermark)

## Documentation

* The main entry point is [pdfcpu.io](https://pdfcpu.io).
* For CLI examples also go to [pdfcpu.io](https://pdfcpu.io). There you will find explanations of all the commands and their parameters.
* For API examples of all pdfcpu operations please refer to [GoDoc](https://godoc.org/github.com/denisbetsi/pdfcpu/pkg/api).

### GoDoc

* [pdfcpu package](https://godoc.org/github.com/denisbetsi/pdfcpu)
* [pdfcpu API](https://godoc.org/github.com/denisbetsi/pdfcpu/pkg/api)
* [pdfcpu CLI](https://godoc.org/github.com/denisbetsi/pdfcpu/pkg/cli)

## Reminder

Always make sure your work is based on the latest commit!<br>
pdfcpu is still *Alpha* - bugfixes are committed on the fly and will be mentioned on the next release notes.<br>

## Demo Screencast

(using older version with a smaller command set)

[![asciicast](resources/demo.png)](https://asciinema.org/a/P5jaAo9kgZXKj2iSA1OqIdLAU)

## Installation

There are no dependencies outside the Go standard library other than `pkg/errors`.<br>
Required go version for building: go1.10 and up

### Using GOPATH

```
go get github.com/denisbetsi/pdfcpu/cmd/...
cd $GOPATH/src/github.com/denisbetsi/pdfcpu/cmd/pdfcpu
go install
pdfcpu version
```

### Using Go Modules (go1.11 and up)

```
git clone https://github.com/denisbetsi/pdfcpu
cd pdfcpu/cmd/pdfcpu
go install
pdfcpu ve
```

## Contributing

### What

* Please open an issue if you find a bug or want to propose a change.
* Feature requests - always welcome!
* Bug fixes - always welcome!
* PRs - also welcome, although I can't promise a merge-in right now.
* pdfcpu is stable but still *Alpha* and occasionally undergoing heavy changes.

### How

* If you want to report a bug please attach the *very verbose* (`pdfcpu cmd -vv ...`) output and ideally a test PDF that you can share.
* Always make sure your contribution is based on the latest commit.
* Please sign your commits.
* Please sign the [CLA](https://cla-assistant.io/hhrutter/pdfcpu) before you submit a PR.

## Contributors

Thanks goes to these wonderful people:
<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
| [<img src="https://avatars1.githubusercontent.com/u/11322155?v=4" width="100px;"/><br /><sub><b>Horst Rutter</b></sub>](https://github.com/denisbetsi)<br /> |[<img src="https://avatars0.githubusercontent.com/u/5140211?v=4" width="100px;"/><br /><sub><b>haldyr</b></sub>](https://github.com/haldyr)<br /> | [<img src="https://avatars3.githubusercontent.com/u/20608155?v=4" width="100px;"/><br /><sub><b>Vyacheslav</b></sub>](https://github.com/SimePel)<br /> | [<img src="https://avatars1.githubusercontent.com/u/617459?v=4" width="100px;"/><br /><sub><b>Erik Unger</b></sub>](https://github.com/ungerik)<br /> ||||
| :---: | :---: | :---: | :---: | :---: | :---: | :---: |
<!-- ALL-CONTRIBUTORS-LIST:END - Do not remove or modify this section -->

## Code of Conduct

Please note that this project is released with a Contributor [Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Disclaimer

Usage of pdfcpu assumes you know about and respect all copyrights of any PDF content you may be processing. This applies to the PDF files as such, their content and in particular all embedded resources like font files or images. Credit goes to [Renee French](https://instagram.com/reneefrench) for creating our beloved Gopher.

## License

Apache-2.0

## Powered By

<p align="center">
  <a href="https://golang.org"> <img src="resources/Go-Logo_Aqua.png" width="200"> </a>
</p>
