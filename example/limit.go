package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// BufferSize 定义缓冲区大小
	BufferSize = 1024
	// KBps 限制下载速度（KB/s）
	KBps = 100
)

// ThrottledReader 包装原始`io.Reader`以限制速度
type ThrottledReader struct {
	reader         io.Reader
	bytesPerSecond int64
}

func (tr *ThrottledReader) Read(p []byte) (int, error) {
	start := time.Now()
	n, err := tr.reader.Read(p)

	// 计算实际所需的时间，以满足指定的速度限制
	duration := time.Duration(float64(n) / float64(tr.bytesPerSecond) * float64(time.Second))
	elapsed := time.Since(start)

	// 如果实际所需时间大于已用时间，那么暂停相应的时间
	if elapsed < duration {
		time.Sleep(duration - elapsed)
	}

	return n, err
}

func main() {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("/Users/qin/dev/go/src/q-meet/notes/package main.go")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		throttledReader := &ThrottledReader{
			reader:         file,
			bytesPerSecond: KBps * 1024,
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=your_file_name")
		io.CopyBuffer(w, throttledReader, make([]byte, BufferSize))
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
