package main

import "time"

// import ("fmt")

type File struct { //FileSystemItem, Readable, Writeable
	name string
	path string
	createdAt time.Time
	modifiedAt time.Time
	data []byte
}

type Katalog struct { //Directory
	name string
	path string
	createdAt time.Time
	modifiedAt time.Time
	items []FileSystemItem
}

type SymLink struct { // FileSystemItem
	name string
	path string
	createdAt time.Time
	modifiedAt time.Time
	refersTo FileSystemItem
}

type ReadOnlyFile struct { // FileSysItem, Readable
	name string
	path string
	createdAt time.Time
	modifiedAt time.Time
	data []byte
}