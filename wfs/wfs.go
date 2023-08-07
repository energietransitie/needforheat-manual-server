// Package wfs implements interfaces, types and methods for working with an fs.FS you can write to.
package wfs

import (
	"errors"
	"io/fs"
)

var (
	ErrInterfaceNotImplemented = errors.New("type does not implement interface")
)

// A WFS is a writable filesystem that implements fs.FS and other interfaces for writing.
type WFS interface {
	fs.FS

	MkdirFS
	MkdirAllFS
	MkdirTempFS
	CreateFileFS
	WriteFileFS
	RemoveFS
	RemoveAllFS
}
