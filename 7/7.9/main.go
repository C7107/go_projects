package main

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"
)

// --- 1. 数据结构 (保持不变) ---
type Track struct {
	Title  string
	Artist string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", 2012, length("3m38s")},
	{"Go", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

// --- 2. 复杂的排序逻辑 (保持不变) ---
// 这就是我们之前的“有状态”排序器
type LessFunc func(x, y *Track) bool

type customSort struct {
	t       []*Track
	history []LessFunc // 历史记录栈
}

func (x *customSort) Len() int      { return len(x.t) }
func (x *customSort) Swap(i, j int) { x.t[i], x.t[j] = x.t[j], x.t[i] }

func (x *customSort) Less(i, j int) bool {
	// 遍历历史记录，多级比较
	for _, less := range x.history {
		p, q := x.t[i], x.t[j]
		if less(p, q) {
			return true
		}
		if less(q, p) {
			return false
		}
	}
	return false
}

func (x *customSort) AddSortKey(f LessFunc) {
	// 把最新的规则放到最前面 (切片头部)
	x.history = append([]LessFunc{f}, x.history...)
}

// 全局的排序器实例 (为了保存点击状态)
var sorter = &customSort{t: tracks}

// --- 3. HTML 模板 ---
// 这是一个内嵌的 HTML 字符串。
// 注意 th 里面的 <a href="/?sort=...">，点击后会重新加载页面并带上参数
var trackTableTmpl = template.Must(template.New("trackList").Parse(`
<!DOCTYPE html>
<html>
<head>
<title>Music Tracks</title>
<style>
  table { border-collapse: collapse; width: 50%; }
  th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
  th { background-color: #f2f2f2; }
  a { text-decoration: none; color: blue; }
  a:hover { text-decoration: underline; }
</style>
</head>
<body>

<h2>Click headers to sort</h2>
<table>
  <tr>
    <th><a href="/?sort=Title">Title</a></th>
    <th><a href="/?sort=Artist">Artist</a></th>
    <th><a href="/?sort=Year">Year</a></th>
    <th><a href="/?sort=Length">Length</a></th>
  </tr>
  {{range .}}
  <tr>
    <td>{{.Title}}</td>
    <td>{{.Artist}}</td>
    <td>{{.Year}}</td>
    <td>{{.Length}}</td>
  </tr>
  {{end}}
</table>

</body>
</html>
`))

// --- 4. HTTP 处理逻辑 ---

func handler(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 中的 sort 参数 (例如 ?sort=Year)
	key := r.FormValue("sort")

	// 2. 根据参数更新排序器的“历史记录”
	// 我们需要根据字符串 "Year" 找到对应的比较函数
	switch key {
	case "Title":
		sorter.AddSortKey(func(x, y *Track) bool { return x.Title < y.Title })
	case "Artist":
		sorter.AddSortKey(func(x, y *Track) bool { return x.Artist < y.Artist })
	case "Year":
		sorter.AddSortKey(func(x, y *Track) bool { return x.Year < y.Year })
	case "Length":
		sorter.AddSortKey(func(x, y *Track) bool { return x.Length < y.Length })
	}

	// 3. 执行排序
	// 如果 key 是空的（刚打开页面），就不添加新规则，直接按现有规则排
	if key != "" || len(sorter.history) > 0 {
		sort.Sort(sorter)
	}

	// 4. 渲染 HTML
	if err := trackTableTmpl.Execute(w, sorter.t); err != nil {
		log.Printf("Template execution failed: %s", err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Listening on http://localhost:8000 ...")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
