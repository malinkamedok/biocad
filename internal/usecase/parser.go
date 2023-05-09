package usecase

import (
	"biocad/internal/entity"
	"container/list"
	"fmt"
	"github.com/dogenzaka/tsv"
	"github.com/signintech/gopdf"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

// FindNewFiles Поиск новых файлов в директории
func (p ParserUseCase) FindNewFiles(root string) ([]string, error) {
	log.Println("finding new files started")

	filesInDB, err := p.repo.GetAllFileNames()
	if err != nil {
		return nil, err
	}

	log.Println("got all filenames already existing in DB")
	var allFiles []string

	//Получение списка файлов в директории
	log.Println(root)
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

	log.Println("all files in directory: ", allFiles)

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

	//Внесение записи о новых файлах в бд
	for _, item := range newFiles {
		err := p.repo.InsertFileName(item)
		if err != nil {
			return nil, err
		}
	}

	log.Println("inserted new file names in DB")
	return newFiles, nil
}

// QueueManager Организация очереди
func (p ParserUseCase) QueueManager(root string, pdfRoot string, fontRoot string) error {
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
		log.Println("Parsing file ", str)
		err := p.ParseData(str, pdfRoot, fontRoot)
		if err != nil {
			return err
		}
		err = p.repo.ChangeProcessedStatus(str)
		if err != nil {
			return err
		}
		queue.Remove(item)
	}

	return nil
}

// ParseData Парсинг файлов
func (p ParserUseCase) ParseData(pathToFile string, pdfRoot string, fontRoot string) error {
	log.Println("parsing data started")

	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	data := entity.Data{}
	parser, err := tsv.NewParser(file, &data)
	if err != nil {
		return err
	}

	for {
		eof, err := parser.Next()
		if eof {
			log.Println("End of file")
			if err != nil {
				return err
			}
			err = p.GetDataToGenerate(pdfRoot, fontRoot)
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
		// Если возникла ошибка - вносится запись о ней
		if err != nil {
			err0 := p.repo.FailureReport(pathToFile)
			if err0 != nil {
				return err0
			}
			return err
		}
	}
}

func (p ParserUseCase) GetDataToGenerate(pdfRoot string, fontRoot string) error {
	log.Println("started collecting data to generate PDF")

	guidList, err := p.repo.GetAllGUIDList()
	if err != nil {
		return err
	}
	log.Println("guid in DB list: ", guidList)

	for _, guid := range guidList {
		log.Println("started generating PDF by GUID ", guid)
		allDataByGUID, err := p.repo.GetAllDataByGUID(guid)
		if err != nil {
			return err
		}
		err = p.GeneratePDF(allDataByGUID, len(allDataByGUID), pdfRoot, fontRoot)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p ParserUseCase) GeneratePDF(file []entity.Data, pages int, pdfRoot string, fontRoot string) error {
	log.Println("generating PDF started with pages count: ", pages)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4, Unit: gopdf.UnitCM})

	err := pdf.AddTTFFont("courier-new", fontRoot)
	if err != nil {
		log.Println("failure in adding font")
		return err
	}
	err = pdf.SetFont("courier-new", "", 11)
	if err != nil {
		log.Println("failure in setting font")
		return err
	}

	for i := 0; i < pages; i++ {
		pdf.AddPage()

		pdf.SetXY(5, 1)
		err := pdf.Cell(nil, file[0].UnitGuid)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 2)
		err = pdf.Cell(nil, "n: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 2)
		err = pdf.Cell(nil, strconv.Itoa(file[i].Id))
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 3)
		err = pdf.Cell(nil, "mqtt: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 3)
		err = pdf.Cell(nil, file[i].Mqtt)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 4)
		err = pdf.Cell(nil, "invid: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 4)
		err = pdf.Cell(nil, file[i].Invid)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 5)
		err = pdf.Cell(nil, "unit_guid: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 5)
		err = pdf.Cell(nil, file[i].UnitGuid)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 6)
		err = pdf.Cell(nil, "msg_id: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 6)
		err = pdf.Cell(nil, file[i].MsgId)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 7)
		err = pdf.Cell(nil, "msg_text: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 7)
		err = pdf.Cell(nil, file[i].MsgText)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 8)
		err = pdf.Cell(nil, "context: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 8)
		err = pdf.Cell(nil, file[i].Context)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 9)
		err = pdf.Cell(nil, "class: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 9)
		err = pdf.Cell(nil, file[i].Class)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 10)
		err = pdf.Cell(nil, "level: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 10)
		err = pdf.Cell(nil, file[i].Level)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 11)
		err = pdf.Cell(nil, "area: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 11)
		err = pdf.Cell(nil, file[i].Area)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 12)
		err = pdf.Cell(nil, "addr: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 12)
		err = pdf.Cell(nil, file[i].Addr)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 13)
		err = pdf.Cell(nil, "block: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 13)
		err = pdf.Cell(nil, file[i].Block)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 14)
		err = pdf.Cell(nil, "type: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 14)
		err = pdf.Cell(nil, file[i].DataType)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 15)
		err = pdf.Cell(nil, "bit: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 15)
		err = pdf.Cell(nil, file[i].Bit)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}

		pdf.SetXY(1, 16)
		err = pdf.Cell(nil, "invert_bit: ")
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
		pdf.SetXY(5, 16)
		err = pdf.Cell(nil, file[i].InvertBit)
		if err != nil {
			log.Println("Failure in writing in pdf")
			return err
		}
	}

	saveFolder := path.Join(pdfRoot, file[0].UnitGuid)
	err = pdf.WritePdf(saveFolder)
	if err != nil {
		log.Println("fail in saving pdf")
		return err
	}
	return nil
}
