package model

import (
	"fmt"
	"time"

	"github.com/RSOI/answer/utils"
	"github.com/jackc/pgx"
)

// RequestInfo interface. Provides one request data.
type RequestInfo struct {
	Request           string    `json:"request"`
	RequestTime       time.Time `json:"request_time"`
	ResponseStatus    int       `json:"response_status"`
	ResponseErrorText string    `json:"response_error_text"`
}

// ServiceStatus interface. Provides usage data.
type ServiceStatus struct {
	Address       string      `json:"address"`
	RequestsCount int         `json:"requests_count"`
	LastUsage     RequestInfo `json:"last_usage"`
}

// GetUsageStatistic provides access to logs
func (service *AService) GetUsageStatistic(host string) (ServiceStatus, error) {
	var err error

	utils.LOG("Accessing database...")
	row := service.Conn.QueryRow(`
		SELECT cnt.*, last_usage.* FROM
			(SELECT count(*) FROM answer.services) AS cnt,
			(SELECT request, request_time, response_status, response_error_text
				FROM answer.services ORDER BY id DESC LIMIT 1
			) AS last_usage
	`)
	var ServiceResponse ServiceStatus
	err = row.Scan(
		&ServiceResponse.RequestsCount,
		&ServiceResponse.LastUsage.Request,
		&ServiceResponse.LastUsage.RequestTime,
		&ServiceResponse.LastUsage.ResponseStatus,
		&ServiceResponse.LastUsage.ResponseErrorText)
	if err == pgx.ErrNoRows {
		ServiceResponse.Address = host
		ServiceResponse.RequestsCount = 0
		ServiceResponse.LastUsage.Request = "Service wasn't used yet"
		ServiceResponse.LastUsage.RequestTime = time.Time{}
		err = nil
	}

	return ServiceResponse, err
}

// LogStat Set request into log db table
func (service *AService) LogStat(request []byte, responseStatus int, responseError string) {
	var err error

	res, err := service.Conn.Exec(`
		INSERT INTO answer.services 
			(request, response_status, response_error_text) VALUES ($1, $2, $3)
	`, string(request), responseStatus, responseError)

	if err != nil {
		utils.LOG(fmt.Sprintf("Error while storing statistic: %s", err.Error()))
	} else {
		if res.RowsAffected() != 1 {
			utils.LOG("Unable store statistic")
		} else {
			utils.LOG("Statistic stored successfully")
		}
	}
}
