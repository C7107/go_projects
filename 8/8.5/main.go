package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
)

func main() {
	// 记录开始时间
	start := time.Now()

	// 创建内存中的图片对象
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 获取当前电脑的CPU逻辑核心数
	workers := runtime.NumCPU()
	fmt.Printf("检测到 %d 个CPU核心，即将启动对应数量的 Workers...\n", workers)

	// 创建任务通道（存放行号 y）
	// 使用 buffered channel 稍微提高一点投放速度，虽然不关键
	jobs := make(chan int, height)

	var wg sync.WaitGroup

	// 1. 启动工人 (Worker Pool)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 工人不停地从 channel 拿行号，直到 channel 关闭
			for y := range jobs {
				for x := 0; x < width; x++ {
					// 计算像素颜色并填充
					// 注意：并发写 img 是安全的，因为每个 goroutine 处理不同的 y 行，
					// 不会同时写入同一个内存地址。
					img.Set(x, y, mandelbrot(complex(float64(x)/width*(xmax-xmin)+xmin, float64(y)/height*(ymax-ymin)+ymin)))
				}
			}
		}()
	}

	// 2. 派发任务 (Producer)
	for y := 0; y < height; y++ {
		jobs <- y
	}
	close(jobs) // 任务发完了，关闭通道，工人读完后会自动退出循环

	// 3. 等待所有工人下班
	wg.Wait()

	fmt.Printf("耗时: %s\n", time.Since(start))

	// 输出文件（不计入并发测试时间，因为 IO 是串行的）
	f, _ := os.Create("mandelbrot.png")
	png.Encode(f, img)
	f.Close()
}

// 具体的数学计算逻辑（保持原样）
func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
