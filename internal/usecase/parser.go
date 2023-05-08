package usecase

import (
	"biocad/internal/entity"
	"container/list"
	"fmt"
	"github.com/dogenzaka/tsv"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type ParserUseCase struct {
	repo ParserRp
}

var _ ParserContract = (*ParserUseCase)(nil)

func NewParserUseCase(repo ParserRp) *ParserUseCase {
	return &ParserUseCase{
		repo: repo,
	}
}

func (p ParserUseCase) FindNewFiles(root string) ([]string, error) {
	log.Println("finding new files started")

	filesInDB, err := p.repo.GetAllFileNames()
	if err != nil {
		return nil, err
	}

	var allFiles []string

	//Получение списка файлов в директории
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ".tsv" {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	//Получение списка новых необработанных файлов в директории
	hashTable := make(map[string]bool)
	var newFiles []string

	for _, item := range filesInDB {
		hashTable[item] = true
	}

	for _, item := range allFiles {
		if _, exists := hashTable[item]; !exists {
			newFiles = append(newFiles, item)
		}
	}

	return newFiles, nil
}

func (p ParserUseCase) QueueManager(root string) error {
	log.Println("queue manager started")

	queue := list.New()

	newFiles, err := p.FindNewFiles(root)
	if err != nil {
		return err
	}
	for _, item := range newFiles {
		queue.PushBack(item)
	}

	for queue.Len() > 0 {
		item := queue.Front()
		str := fmt.Sprintf("%v", item.Value)
		log.Println(str)
		err := p.ParseData(str)
		if err != nil {
			return err
		}
		queue.Remove(item)
	}

	return nil
}

func (p ParserUseCase) ParseData(pathToFile string) error {
	log.Println("parsing data started")

	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data := entity.Data{}
	parser, err := tsv.NewParser(file, &data)
	if err != nil {
		return err
	}

	for {
		eof, err := parser.Next()
		if eof {
			log.Println("End of file")
			err = p.repo.ChangeProcessedStatus(pathToFile)
			if err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			log.Println("Error in parsing")
			return err
		}
		fmt.Println(data)
		err = p.repo.InsertData(data)
		if err != nil {
			return err
		}
	}
}
