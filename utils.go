package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

// GetTimeStamp ...
func GetTimeStamp() string {
	return time.Now().Format("2006-01-02 15-04-05")
}

// CopyFile ...
func CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CopyDir ...
func CopyDir(srcDir string, dstDir string) error {
	// fmt.Print([]byte(dstDir))
	err := os.Mkdir(dstDir, 0666)
	if err != nil {
		log.Println(err)
	}
	fileInfoList, _ := ioutil.ReadDir(srcDir)
	for i := 0; i < len(fileInfoList); i++ {
		// fmt.Println("Copying: ", fileInfoList[i].Name(), fileInfoList[i].IsDir(), "...")
		if fileInfoList[i].IsDir() {
			CopyDir(path.Join(srcDir, fileInfoList[i].Name()), path.Join(dstDir, fileInfoList[i].Name()))
		} else {
			CopyFile(path.Join(srcDir, fileInfoList[i].Name()), path.Join(dstDir, fileInfoList[i].Name()))
		}
	}
	return nil
}
