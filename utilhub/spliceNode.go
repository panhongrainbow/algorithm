package utilhub

import (
	"os"
	"path"
)

// LinuxSpliceStreamWrite wraps the LinuxSpliceStreamWrite function to write a file stream to a file.
func (fn FileNode) LinuxSpliceStreamWrite(filename string, fileFlag int, filePerm os.FileMode) (dataChan chan [][]byte, finishChan chan struct{}, err error) {
	// Construct the absolute path of the file by joining the transfer directory and the filename.
	absPath := path.Join(fn.transfer, filename)

	// Call the LinuxSpliceStreamWrite function with the absolute path and other parameters.
	return LinuxSpliceStreamWrite(absPath, fileFlag, filePerm)
}
