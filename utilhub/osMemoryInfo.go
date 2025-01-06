package utilhub

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// =====================================================================================================================
//                  🛠️ OS Memory Info (Tool)
// OS Memory Info is a collection of functions for retrieving available memory information from the operating system.
// =====================================================================================================================

// MemInfoFields defines the fields in /proc/meminfo
const (
	MemFree  = "MemFree"
	SwapFree = "SwapFree"
)

// GetLinuxMemoryValue ⛏️ retrieves a specific memory value from the Linux system's /proc/meminfo file.
func GetLinuxMemoryValue(field string) (uint64, error) {
	// Open the /proc/meminfo file for reading.
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		// If the file cannot be opened, return an error.
		return 0, fmt.Errorf("failed to open /proc/meminfo: %v", err)
	}
	defer file.Close() // Close the file when we're done with it.

	// Create a buffer to read the file contents into.
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		// If the file cannot be read, return an error.
		return 0, fmt.Errorf("failed to read /proc/meminfo: %v", err)
	}

	// Split the file contents into individual lines.
	lines := strings.Split(string(buf[:n]), "\n")

	// Iterate over each line in the file.
	for _, line := range lines {
		// Check if the line starts with the specified field name.
		if strings.HasPrefix(line, field+":") {
			// Split the line into individual fields.
			fields := strings.Fields(line)

			// Check if the line has at least two fields (the field name and the value).
			if len(fields) < 2 {
				// If the line is malformed, return an error.
				return 0, fmt.Errorf("error parsing /proc/meminfo")
			}

			// Parse the memory value from the second field.
			memoryValue, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				// If the value cannot be parsed, return an error.
				return 0, fmt.Errorf("failed to convert memory value: %v", err)
			}

			// Return the parsed memory value.
			return memoryValue, nil
		}
	}

	// If the specified field is not found, return an error.
	return 0, fmt.Errorf("%q field not found in /proc/meminfo", field)
}

// GetLinuxAvailableMemory ⛏️ retrieves the total available memory on a Linux system.
func GetLinuxAvailableMemory() (uint64, error) {
	// Check if the current operating system is Linux.
	if runtime.GOOS != "linux" {
		return 0, fmt.Errorf("GetLinuxAvailableMemory is only supported on Linux; Current OS is %s", runtime.GOOS)
	}

	// Retrieve the amount of free memory on the system.
	memFree, err := GetLinuxMemoryValue(MemFree)
	if err != nil {
		// If an error occurs while retrieving the free memory, return the error.
		return 0, err
	}

	// Retrieve the amount of free swap space on the system.
	swapFree, err := GetLinuxMemoryValue(SwapFree)
	if err != nil {
		// If an error occurs while retrieving the free swap space, return the error.
		return 0, err
	}

	// Calculate the total available memory by adding the free memory and swap space.
	return memFree + swapFree, nil
}

// List of memory fields for future use.
// const (
// MemFree  = "MemFree"
// SwapFree = "SwapFree"
/*
	MemTotal          = "MemTotal"
	MemAvailable      = "MemAvailable"
	Buffers           = "Buffers"
	Cached            = "Cached"
	SwapCached        = "SwapCached"
	Active            = "Active"
	Inactive          = "Inactive"
	ActiveAnon        = "Active(anon)"
	InactiveAnon      = "Inactive(anon)"
	ActiveFile        = "Active(file)"
	InactiveFile      = "Inactive(file)"
	Unevictable       = "Unevictable"
	Mlocked           = "Mlocked"
	SwapTotal         = "SwapTotal"
	Zswap             = "Zswap"
	Zswapped          = "Zswapped"
	Dirty             = "Dirty"
	Writeback         = "Writeback"
	AnonPages         = "AnonPages"
	Mapped            = "Mapped"
	Shmem             = "Shmem"
	KReclaimable      = "KReclaimable"
	Slab              = "Slab"
	SReclaimable      = "SReclaimable"
	SUnreclaim        = "SUnreclaim"
	KernelStack       = "KernelStack"
	PageTables        = "PageTables"
	SecPageTables     = "SecPageTables"
	NFSUnstable       = "NFS_Unstable"
	Bounce            = "Bounce"
	WritebackTmp      = "WritebackTmp"
	CommitLimit       = "CommitLimit"
	CommittedAS       = "Committed_AS"
	VmallocTotal      = "VmallocTotal"
	VmallocUsed       = "VmallocUsed"
	VmallocChunk      = "VmallocChunk"
	Percpu            = "Percpu"
	HardwareCorrupted = "HardwareCorrupted"
	AnonHugePages     = "AnonHugePages"
	ShmemHugePages    = "ShmemHugePages"
	ShmemPmdMapped    = "ShmemPmdMapped"
	FileHugePages     = "FileHugePages"
	FilePmdMapped     = "FilePmdMapped"
	Unaccepted        = "Unaccepted"
	HugePagesTotal    = "HugePages_Total"
	HugePagesFree     = "HugePages_Free"
	HugePagesRsvd     = "HugePages_Rsvd"
	HugePagesSurp     = "HugePages_Surp"
	Hugepagesize      = "Hugepagesize"
	Hugetlb           = "Hugetlb"
	DirectMap4k       = "DirectMap4k"
	DirectMap2M       = "DirectMap2M"
	DirectMap1G       = "DirectMap1G"
*/
// )
