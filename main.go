package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// set input data dir
	dirname, err := filepath.Abs(filepath.Dir(("C:/0x02LY/GitHub/wechat-image-dat-convert/dat/")))
	if err != nil {
		fmt.Println("error!!")
	}
	fmt.Println("dirname ", dirname)

	//set output dir
	outputDirname := filepath.Join(dirname, "output")
	fmt.Println("输出路径为： ", outputDirname)

	//get dat files
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println("errors!!")
	}
	var datFiles []string
	for key, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			fmt.Println("dir ", key, fileInfo.Name())
		} else {
			if !strings.HasSuffix(fileInfo.Name(), ".dat") {
				fmt.Println("suffix is not .dat ", fileInfo.Name())
				continue
			}
			inFilename := fmt.Sprintf("%s\\%s", dirname, fileInfo.Name())
			fmt.Println("inFilename is ", inFilename)
			datFiles = append(datFiles, inFilename)
		}
	}
	fmt.Println(datFiles)

	if len(datFiles) < 1 {
		fmt.Println("no dat file")
	}

	var wg sync.WaitGroup

	for key, datFile := range datFiles {
		wg.Add(1)
		go wechatDatToImagePip(datFile, outputDirname, &wg)
		fmt.Printf("completed: %d %s \n", key, datFile)
	}

	wg.Wait()
	// time.Sleep(3 * 1e9) // sleep for 5 seconds
}

func wechatDatToImagePip(datFile string, outputDirname string, wg *sync.WaitGroup) {
	defer wg.Done()
	datBytes, err := ioutil.ReadFile((datFile))
	if err != nil {
		fmt.Println(err)
	}

	imgBytes, err := WechatDatToImage(datBytes)
	if err != nil {
		fmt.Println(err)
	}

	datFilename := filepath.Base(datFile) + ".jpg"
	outFilename := filepath.Join(outputDirname, datFilename)
	err = ensureDirExist(outFilename)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(outFilename, imgBytes, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

// WechatDatToImage 可以将微信 dat 字节数组转图片字节数组
func WechatDatToImage(dat []byte) (img []byte, err error) {
	var buffer bytes.Buffer
	for _, value := range dat {
		err = buffer.WriteByte(value ^ 0xA7)
		if err != nil {
			return
		}
	}
	img = buffer.Bytes()
	return
}

// 写个函数，确保文件夹存在，省的重复写
func ensureDirExist(path string) error {
	dir := filepath.Dir(path)   //获取文件所在路径
	exists := isPathExists(dir) //判断路径是否存在
	if !exists {                //如果不存在
		err := os.MkdirAll(dir, os.ModePerm) //创建文件夹
		if err != nil {
			return err
		}
	}
	return nil
}

// 判断路径是否存在
func isPathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
