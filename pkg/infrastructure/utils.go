package infrastructure

import "os"

func createFileStorageDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
