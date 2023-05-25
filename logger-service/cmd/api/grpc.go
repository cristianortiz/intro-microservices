package main

import (
	"context"
	"logger-service/data"
	"logger-service/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer // a must have for backguards compatibility
	Models                             data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

}
