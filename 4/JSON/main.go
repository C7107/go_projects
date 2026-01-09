package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/C7107/go_projects/4/JSON/github"
)

func plan1() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}

func plan2() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total items: %d\n", result.TotalCount)

	// 1. 定义三个切片用于存储分类后的 Issue
	var (
		lessThanMonth []*github.Issue
		lessThanYear  []*github.Issue
		moreThanYear  []*github.Issue
	)

	// 2. 计算时间阈值
	now := time.Now()
	monthAgo := now.AddDate(0, -1, 0) // 当前时间减去1个月
	yearAgo := now.AddDate(-1, 0, 0)  // 当前时间减去1年

	// 3. 遍历并分类
	for _, item := range result.Items {
		if item.CreatedAt.After(monthAgo) {
			lessThanMonth = append(lessThanMonth, item)
		} else if item.CreatedAt.After(yearAgo) {
			lessThanYear = append(lessThanYear, item)
		} else {
			moreThanYear = append(moreThanYear, item)
		}
	}

	// 4. 分别打印
	printIssues("Less than a month", lessThanMonth)
	printIssues("Less than a year", lessThanYear)
	printIssues("More than a year", moreThanYear)
}

func printIssues(category string, items []*github.Issue) {
	if len(items) == 0 {
		return
	}
	fmt.Printf("\n--- %s (%d issues) ---\n", category, len(items))
	for _, item := range items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}
func main() {
	plan1()
}
