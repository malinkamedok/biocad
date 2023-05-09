package usecase

import (
	"biocad/internal/entity"
	"context"
)

type (
	UserRp interface {
		GetAllGuidList(ctx context.Context) ([]string, error)
		GetDataByUnitGuid(ctx context.Context, guid string, limit uint64, offset uint64) ([]entity.Data, error)
		CheckGuidExists(ctx context.Context, guid string) (bool, error)
	}

	ParserRp interface {
		InsertData(data entity.Data) error
		InsertFileName(fileName string) error
		GetAllFileNames() ([]string, error)
		ChangeProcessedStatus(fileName string) error
		FailureReport(fileName string) error
		GetAllGUIDList() ([]string, error)
		GetAllDataByGUID(guid string) ([]entity.Data, error)
	}

	UserContract interface {
		GetAllUniqueGuidList(ctx context.Context) ([]string, error)
		GetAllDataByGuid(ctx context.Context, guid string, limit uint64, page uint64) ([]entity.Data, error)
	}

	ParserContract interface {
		FindNewFiles(root string) ([]string, error)
		QueueManager(root string, pdfRoot string, fontRoot string) error
		ParseData(pathToFile string, pdfRoot string, fontRoot string) error
		GetDataToGenerate(pdfRoot string, fontRoot string) error
		GeneratePDF(file []entity.Data, pages int, pdfRoot string, fontRoot string) error
	}
)
