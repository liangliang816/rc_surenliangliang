package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// ApiRunRecord 结构体
type ApiRunRecord struct {
	ID             uint   `json:"id"`
	ApiCode        string `json:"api_code"`
	Status         string `json:"status"`
	ExecutionCount uint   `json:"execution_count"`
	ExecutionTime  string `json:"execution_time"`
	ResponseResult string `json:"response_result"`
	IsSuccess      bool   `json:"is_success"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// CreateApiJob 创建新的API任务
func CreateApiJob(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	var job ApiJob
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	result, err := db.Exec(
		"INSERT INTO api_job (name, code, url, request_header, request_body, is_executed) VALUES (?, ?, ?, ?, ?, ?)",
		job.Name, job.Code, job.URL, job.RequestHeader, job.RequestBody, job.IsExecuted,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "任务创建成功"})
}

// GetApiJobs 获取所有API任务
func GetApiJobs(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, code, url, request_header, request_body, is_executed, created_at, updated_at FROM api_job")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		return
	}
	defer rows.Close()

	var jobs []ApiJob
	for rows.Next() {
		var job ApiJob
		err := rows.Scan(&job.ID, &job.Name, &job.Code, &job.URL, &job.RequestHeader, &job.RequestBody, &job.IsExecuted, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析数据失败"})
			return
		}
		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, jobs)
}

// GetApiJob 根据ID获取单个API任务
func GetApiJob(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var job ApiJob
	err = db.QueryRow(
		"SELECT id, name, code, url, request_header, request_body, is_executed, created_at, updated_at FROM api_job WHERE id = ?",
		id,
	).Scan(&job.ID, &job.Name, &job.Code, &job.URL, &job.RequestHeader, &job.RequestBody, &job.IsExecuted, &job.CreatedAt, &job.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询任务失败"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// UpdateApiJob 更新API任务
func UpdateApiJob(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var job ApiJob
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	_, err = db.Exec(
		"UPDATE api_job SET name = ?, code = ?, url = ?, request_header = ?, request_body = ?, is_executed = ? WHERE id = ?",
		job.Name, job.Code, job.URL, job.RequestHeader, job.RequestBody, job.IsExecuted, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务更新成功"})
}

// DeleteApiJob 删除API任务
func DeleteApiJob(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	_, err = db.Exec("DELETE FROM api_job WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务删除成功"})
}

// GetApiRunRecords 获取所有执行记录
func GetApiRunRecords(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, api_code, status, execution_count, execution_time, response_result, is_success, created_at, updated_at FROM api_run_record")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询执行记录失败"})
		return
	}
	defer rows.Close()

	var records []ApiRunRecord
	for rows.Next() {
		var record ApiRunRecord
		err := rows.Scan(&record.ID, &record.ApiCode, &record.Status, &record.ExecutionCount, &record.ExecutionTime, &record.ResponseResult, &record.IsSuccess, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析数据失败"})
			return
		}
		records = append(records, record)
	}

	c.JSON(http.StatusOK, records)
}

// GetApiRunRecord 根据ID获取单个执行记录
func GetApiRunRecord(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var record ApiRunRecord
	err = db.QueryRow(
		"SELECT id, api_code, status, execution_count, execution_time, response_result, is_success, created_at, updated_at FROM api_run_record WHERE id = ?",
		id,
	).Scan(&record.ID, &record.ApiCode, &record.Status, &record.ExecutionCount, &record.ExecutionTime, &record.ResponseResult, &record.IsSuccess, &record.CreatedAt, &record.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "执行记录不存在"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询执行记录失败"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// GetApiRunRecordsByApiCode 根据API code获取执行记录
func GetApiRunRecordsByApiCode(c *gin.Context) {
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库连接失败"})
		return
	}
	defer db.Close()

	apiCode := c.Param("api_code")
	if apiCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的API code"})
		return
	}

	rows, err := db.Query(
		"SELECT id, api_code, status, execution_count, execution_time, response_result, is_success, created_at, updated_at FROM api_run_record WHERE api_code = ?",
		apiCode,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询执行记录失败"})
		return
	}
	defer rows.Close()

	var records []ApiRunRecord
	for rows.Next() {
		var record ApiRunRecord
		err := rows.Scan(&record.ID, &record.ApiCode, &record.Status, &record.ExecutionCount, &record.ExecutionTime, &record.ResponseResult, &record.IsSuccess, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析数据失败"})
			return
		}
		records = append(records, record)
	}

	c.JSON(http.StatusOK, records)
}
