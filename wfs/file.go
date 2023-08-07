package wfs

import (
	"io"
	"io/fs"
)

// A File provides read and write access to a single file.
// The File interface is the minimum implementation required of the file.
// Directory files should also implement ReadDirFile.
// A file may implement io.ReaderAt or io.Seeker as optimizations.
//
// Unlike an fs.File, a wfs.File also implements the io.Writer interface.
type File interface {
	fs.File
	io.Writer
}

// CreateFileFS is the interface implemented by a file system
// that provides an implementation of CreateFile.
type CreateFileFS interface {
	// Create creates or truncates the named file. If the file already exists,
	// it is truncated. If the file does not exist, it is created with mode 0666
	// (before umask). If successful, methods on the returned File can
	// be used for I/O; the associated file descriptor has mode O_RDWR.
	// If there is an error, it will be of type *PathError.
	CreateFile(name string) (File, error)
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
// If there is an error, it will be of type *PathError.
func CreateFile(fsys fs.FS, name string) (File, error) {
	if fsys, ok := fsys.(CreateFileFS); ok {
		return fsys.CreateFile(name)
	}
	return nil, &fs.PathError{Op: "createfile", Path: name, Err: ErrInterfaceNotImplemented}
}

// WriteFileFS is the interface implemented by a file system
// that provides an implementation of WriteFile.
type WriteFileFS interface {
	fs.FS

	// WriteFile writes data to the named file, creating it if necessary.
	// If the file does not exist, WriteFile creates it with permissions perm (before umask);
	// otherwise WriteFile truncates it before writing, without changing permissions.
	// Since Writefile requires multiple system calls to complete, a failure mid-operation
	// can leave the file in a partially written state.
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

func WriteFile(fsys fs.FS, name string, data []byte, perm fs.FileMode) error {
	if fsys, ok := fsys.(WriteFileFS); ok {
		return fsys.WriteFile(name, data, perm)
	}
	return &fs.PathError{Op: "writefile", Path: name, Err: ErrInterfaceNotImplemented}
}

// RemoveFS is the interface implemented by a file system
// that provides an implementation of Remove.
type RemoveFS interface {
	fs.FS

	// Remove removes the named file or (empty) directory.
	// If there is an error, it will be of type *PathError.
	Remove(name string) error
}

// Remove removes the named file or (empty) directory.
// If there is an error, it will be of type *PathError.
func Remove(fsys fs.FS, name string) error {
	if fsys, ok := fsys.(RemoveFS); ok {
		return fsys.Remove(name)
	}
	return &fs.PathError{Op: "remove", Path: name, Err: ErrInterfaceNotImplemented}
}
