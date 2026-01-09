// 模拟抓取函数
package links

func Extract(url string) ([]string, error) {
	// 模拟返回几个假链接
	return []string{url + "/a", url + "/b"}, nil
}
