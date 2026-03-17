package chat

import (
	"go.uber.org/zap"
)

type indexJob struct {
	CompanyID int
	FilePath  string
	FileName  string
	ProcessID int
}

var indexQueue = make(chan indexJob, 50)

func init() {
	go indexWorker()
}

func indexWorker() {
	for job := range indexQueue {
		processJob(job)
	}
}

func processJob(job indexJob) {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorw("Panic durante indicizzazione, job saltato", "process_id", job.ProcessID, "error", r)
		}
	}()

	azureClient := GetAzureClient()
	if azureClient == nil || !azureClient.IsConfigured() {
		return
	}
	zap.S().Infow("Avvio indicizzazione documento", "process_id", job.ProcessID, "file", job.FileName)
	if err := IndexFile(azureClient, job.CompanyID, job.FilePath); err != nil {
		zap.S().Warnw("Errore indicizzazione documento", "error", err)
	} else {
		zap.S().Infow("Indicizzazione completata", "process_id", job.ProcessID)
	}
}

func EnqueueIndexing(companyID int, filePath, fileName string, processID int) {
	select {
	case indexQueue <- indexJob{CompanyID: companyID, FilePath: filePath, FileName: fileName, ProcessID: processID}:
	default:
		zap.S().Warnw("Coda indicizzazione piena, job scartato", "process_id", processID, "file", fileName)
	}
}
