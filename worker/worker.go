package worker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ApiJob 结构体
type ApiJob struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	URL           string `json:"url"`
	RequestHeader string `json:"request_header"`
	RequestBody   string `json:"request_body"`
	IsExecuted    bool   `json:"is_executed"`
}

// StartWorker 启动后台执行任务
func StartWorker() {
	// 初始化数据库连接池
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// 启动主工作协程
	go startMainWorker(db)
}

// startMainWorker 启动主工作协程
func startMainWorker(db *sql.DB) {
	fmt.Println("Main worker started")
	defer db.Close()

	for {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Worker recovered from panic: %v\n", r)
				// 重启工作协程
				go startMainWorker(db)
			}
		}()

		// 查找未执行的任务
		var job ApiJob
		err := db.QueryRow(
			"SELECT id, name, code, url, request_header, request_body, is_executed FROM api_job WHERE is_executed = 0 LIMIT 1",
		).Scan(&job.ID, &job.Name, &job.Code, &job.URL, &job.RequestHeader, &job.RequestBody, &job.IsExecuted)

		if err == sql.ErrNoRows {
			// 没有未执行的任务
			time.Sleep(10 * time.Second)
			continue
		} else if err != nil {
			fmt.Printf("查询任务失败: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// 开始执行任务
		success, response, status := executeRequest(job, 5)

		// 开始事务
		tx, err := db.Begin()
		if err != nil {
			fmt.Printf("开始事务失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 更新任务状态
		_, err = tx.Exec("UPDATE api_job SET is_executed = 1 WHERE id = ?", job.ID)
		if err != nil {
			tx.Rollback()
			fmt.Printf("更新任务状态失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 写入执行记录
		_, err = tx.Exec(
			"INSERT INTO api_run_record (api_code, status, execution_count, execution_time, response_result, is_success) VALUES (?, ?, ?, ?, ?, ?)",
			job.Code, status, 1, time.Now(), response, success,
		)
		if err != nil {
			tx.Rollback()
			fmt.Printf("写入执行记录失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 提交事务
		if err := tx.Commit(); err != nil {
			fmt.Printf("提交事务失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("任务执行完成: %s, 状态: %s, 成功: %v\n", job.Name, status, success)

		// 处理完一个任务后，短暂休眠再处理下一个
		time.Sleep(1 * time.Second)
	}
}

// processJob 处理单个任务（保留函数签名，方便后续扩展）
func processJob() {
	// 现在逻辑已移至 startMainWorker
}

// executeRequest 执行HTTP请求，支持重试
func executeRequest(job ApiJob, maxRetries int) (bool, string, string) {
	var lastError error

	for i := 0; i < maxRetries; i++ {
		// 解析请求头
		headers := make(http.Header)
		if job.RequestHeader != "" {
			var headerMap map[string]string
			if err := json.Unmarshal([]byte(job.RequestHeader), &headerMap); err == nil {
				for key, value := range headerMap {
					headers.Set(key, value)
				}
			}
		}

		// 创建请求
		req, err := http.NewRequest("POST", job.URL, strings.NewReader(job.RequestBody))
		if err != nil {
			lastError = err
			fmt.Printf("创建请求失败: %v, 重试 %d/%d\n", err, i+1, maxRetries)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// 设置请求头
		for key, values := range headers {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// 发送请求
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastError = err
			fmt.Printf("发送请求失败: %v, 重试 %d/%d\n", err, i+1, maxRetries)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		// 读取响应
		defer resp.Body.Close()
		response := fmt.Sprintf("HTTP %d", resp.StatusCode)

		// 检查是否成功
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return true, response, "success"
		}

		lastError = fmt.Errorf("HTTP error: %s", resp.Status)
		fmt.Printf("请求失败: %s, 重试 %d/%d\n", resp.Status, i+1, maxRetries)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return false, fmt.Sprintf("Error: %v", lastError), "failed"
}
