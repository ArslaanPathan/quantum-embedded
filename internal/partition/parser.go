package partition

import (
	"fmt"
	"strings"
	"strconv"
	"os"
	"github.com/BurntSushi/toml"
)

func ParseDiskLayout(filepath string) (*DiskLayout, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %w", filepath, err)
	}

	var diskLayout DiskLayout
	_, err = toml.Decode(string(data), &diskLayout)
	if err != nil {
		return nil, fmt.Errorf("error while parsing TOML! filepath: %s, error: %w", filepath, err)
	}

	// parse sizes and start bytes at runtime 
	totalSize, err := ParseSize(diskLayout.TotalSize)
	if err != nil {
		return nil, fmt.Errorf("error while parsing total size! TOML filepath %s, error: %w", filepath, err)
	}
	diskLayout.TotalSizeBytes = totalSize

	// 1MiB overhead to account for GPT/MBR
	usedBytes := uint64(1024 * 1024)

	for i := range diskLayout.Partitions {
		p := &diskLayout.Partitions[i]
		if p.Size == "0" {
			p.SizeBytes = totalSize - usedBytes 
		} else {
			size, err := ParseSize(p.Size)
			if err != nil {
				return nil, fmt.Errorf("error parsing size for partition %s: %w", p.Name, err)
			}
			p.SizeBytes = size
		}

		p.StartBytes = usedBytes
		usedBytes += p.SizeBytes
	}

	return &diskLayout, nil 
}

func ParseSize(size string) (uint64, error) {
	if size == "0" {
		return 0, nil
	}

	if len(size) < 2 {
		return 0, fmt.Errorf("invalid size format: %s", size)
	}

	unit := strings.ToUpper(size[len(size)-1:])
	valueStr := size[:len(size)-1]

	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size value: %s", size)
	}

	switch unit {
		case "K":
			return value * 1024, nil
		case "M":
			return value * 1024 * 1024, nil
		case "G":
			return value * 1024 * 1024 * 1024, nil
		case "T":
			return value * 1024 * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unknown size unit: %s (use K, M, G, or T)", unit)
	}
}
