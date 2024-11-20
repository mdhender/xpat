// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package stdlib

import "os"

// IsDirExists returns true if the path exists and is a directory.
func IsDirExists(path string) (bool, error) {
	sb, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return sb.IsDir(), nil
}

// IsFileExists returns true if the path exists and is a regular file.
func IsFileExists(path string) (bool, error) {
	sb, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	} else if sb.IsDir() {
		return false, nil
	}
	return sb.Mode().IsRegular(), nil
}
