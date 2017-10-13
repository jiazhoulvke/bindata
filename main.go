package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jiazhoulvke/goutil"
)

var (
	pkgName  string
	output   string
	dataPath string
)

func init() {
	flag.StringVar(&pkgName, "package", "main", "package name")
	flag.StringVar(&output, "output", "bindata.go", "output filename")
	flag.StringVar(&dataPath, "data", "", "bin data path")
}

func main() {
	flag.Parse()
	var err error
	if dataPath == "" {
		fmt.Println("need data path")
		os.Exit(1)
	}
	gofmt, err := exec.LookPath("gofmt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data, err := fileMap(dataPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := bytes.NewBufferString("")
	for k, v := range data {
		buf.WriteString(fmt.Sprintf(`"%s": []byte("%s"),`, k, v))
	}
	dirPath := filepath.Dir(output)
	if err := goutil.CreateParentDir(dirPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bindata, err := os.Create(output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bindata.WriteString(fmt.Sprintf(`
//this file genetated by bindata (https://github.com/jiazhoulvke/bindata)
package %s

import (
	"bytes"
	"errors"
	"strings"
	"fmt"
)

//BinData bin data
var BinData  map[string][]byte

func init() {
	BinData = map[string][]byte{
		%s
	}
	var err error
	for k,v:= range BinData {
		BinData[k],err= hex2bytes(string(v))
		if err!=nil {
			fmt.Println("convert bin data error:",err)
		}
	}
}

func hex2bytes(s string) ([]byte, error) {
	bs := bytes.NewBufferString("")
	sBytes := []byte(strings.ToLower(s))
	l := len(sBytes)
	if l%%2 != 0 {
		return bs.Bytes(), errors.New("error data")
	}
	var b uint8
	for i := 0; i < l; i += 2 {
		b = hex2dec(sBytes[i])*16 + hex2dec(sBytes[i+1])
		bs.WriteByte(byte(b))
	}
	return bs.Bytes(), nil
}

func hex2dec(b byte) uint8 {
	switch b {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	case 'a':
		return 10
	case 'b':
		return 11
	case 'c':
		return 12
	case 'd':
		return 13
	case 'e':
		return 14
	case 'f':
		return 15
	}
	return 0
}`, pkgName, buf.String()))
	bindata.Close()

	cmd := exec.Command(gofmt, output)
	sourceCode, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(output, sourceCode, 0666)
	if err != nil {
		panic(err)
	}
}

func fileMap(dataPath string) (map[string][]byte, error) {
	data := make(map[string][]byte)
	err := filepath.Walk(dataPath, func(p string, info os.FileInfo, err1 error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(dataPath, p)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		data[relPath] = []byte(fmt.Sprintf("%x", content))
		return nil
	})
	return data, err
}
