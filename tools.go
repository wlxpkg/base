/*
 * @Author: qiuling
 * @Date: 2019-04-29 19:32:36
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-04-29 19:45:07
 */
package pkg

import (
	"fmt"
	"runtime/debug"
)

func R(data interface{}, name string) {
	fmt.Printf("%v\n", name)
	fmt.Printf("%v\n", data)
}

func D(data interface{}) {
	fmt.Printf("%s\n", debug.Stack())
	fmt.Printf("%v\n", data)
}
