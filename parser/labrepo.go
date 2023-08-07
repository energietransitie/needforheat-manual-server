package parser

import (
	"io/fs"

	"github.com/go-git/go-git/v5/plumbing/transport"
)

// LabRepoSource is a git repo that contains manuals made by a lab.
type LabRepoSource struct {
	fs.FS
}

// Create a new source filesystem from a directory at path.
func NewLabRepoSource(url string, auth transport.AuthMethod) (fs.FS, error) {
	gitFS, _, err := newGitFSWithAuth(url, auth)
	if err != nil {
		return nil, err
	}

	return LabRepoSource{gitFS}, nil
}

// Get the path to copy a file to at the destination filesystem.
// Use filePath at sourceFS to determine the path the file should be copied to at the destination filesystem.
func (repo LabRepoSource) GetDestinationFilePath(filePath string) string {
	return filePath
}

// Get the path to copy a directory to at the destination filesystem.
// Use dirPath at sourceFS to determine the path the directory should be copied to at the destination filesystem.
func (repo LabRepoSource) GetDestinationDirPath(dirPath string) string {
	return dirPath
}
