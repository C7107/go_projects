// 文件路径: E:\go_projects\2\main.go
package main

import (
	"fmt"
	// 路径公式：模块名(go.mod里写的) + 文件夹路径
	"os"
	"strconv"

	"github.com/C7107/go_projects/2/baohewenj/tempconv"
)

func main() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cf: %v\n", err)
			os.Exit(1)
		}
		f := tempconv.Fahrenheit(t)
		c := tempconv.Celsius(t)
		fmt.Printf("%s = %s, %s = %s\n",
			f, tempconv.FToC(f), c, tempconv.CToF(c))
	}
}
