# go-onedrive [![Build Status](https://travis-ci.org/ggordan/go-onedrive.svg?branch=master)](https://travis-ci.org/ggordan/go-onedrive) [![Coverage Status](https://coveralls.io/repos/ggordan/go-onedrive/badge.svg?branch=master)](https://coveralls.io/r/ggordan/go-onedrive?branch=master)


go-onedrive is a Go client library for accessing the Microsoft OneDrive API.

# Documentation

https://godoc.org/github.com/ggordan/go-onedrive

# Example

Get an access token via the [token flow](http://onedrive.github.io/auth/msa_oauth.htm#token-flow) or the [code flow](http://onedrive.github.io/auth/msa_oauth.htm#code-flow)...

# TODO

- [x] Drives
 - [x] Get Default Drive
 - [x] Get Drive
 - [x] List all available drives
- [ ] Items
 - [ ] Create
 	- [x] Create folder
 - [ ] Copy
 	- [x] Copy file/folder
 	- [ ] Async job to track progress
 - [x] Delete
 - [ ] Download
 - [x] List children
 - [ ] Search
 - [x] Move
 - [ ] Upload
 	- [x] Simple item upload <100MB
 	- [ ] Resumable item upload
 	- [x] Upload from URL

# License

MIT
