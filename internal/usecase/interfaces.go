package usecase

import "biocad/internal/entity"

type (
	UserRp interface {
	}

	ParserRp interface {
		InsertData(data entity.Data) error
		InsertFileName(fileName string) error
		GetAllFileNames() ([]string, error)
		ChangeProcessedStatus(fileName string) error
	}

	UserContract interface {
	}

	ParserContract interface {
		FindNewFiles(root string) ([]string, error)
		ParseData(pathToFile string) error
	}
)
