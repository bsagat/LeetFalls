package dal

import (
	"GonIO/internal/domain"
	csvparser "GonIO/pkg/myCSV"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ObjectCSV struct{}

func NewObjectCSVRepo() *ObjectCSV {
	return &ObjectCSV{}
}

var _ domain.ObjectDal = (*ObjectCSV)(nil)

func (repo ObjectCSV) ListObjects(bucketname string) ([]domain.Object, error) {
	metaPath := filepath.Join(domain.BucketsPath, bucketname, "objects.csv")
	metaFile, err := os.Open(metaPath)
	if err != nil {
		return nil, err
	}
	defer metaFile.Close()

	reader := csv.NewReader(metaFile)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var objects []domain.Object
	for i, record := range records {
		if i == 0 || len(record) < 4 {
			continue
		}
		objects = append(objects, domain.Object{
			ObjectKey:    record[0],
			Size:         record[1],
			ContentType:  record[2],
			LastModified: record[3],
		})
	}

	return objects, nil
}

func (repo ObjectCSV) UploadObject(bucketname, objectname string, image io.ReadCloser, contentLen int64, fileType string) error {
	defer image.Close()

	objectPath := filepath.Join(domain.BucketsPath, bucketname, objectname)
	file, err := os.OpenFile(objectPath, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, image); err != nil {
		return err
	}

	now := time.Now()
	metaPath := filepath.Join(domain.BucketsPath, bucketname, "objects.csv")
	metaFile, err := os.OpenFile(metaPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer metaFile.Close()

	metaData := []string{objectname, strconv.Itoa(int(contentLen)), fileType, now.Format(time.ANSIC)}
	if err = csvparser.WriteCSV(metaFile, metaData); err != nil {
		return err
	}

	return csvparser.ReWriteCSV(domain.BucketsMetaPath, []string{bucketname, "", now.Format(time.ANSIC), domain.ActiveMark})
}

func (repo ObjectCSV) RetrieveObject(bucketname, objectname string) (io.ReadCloser, error) {
	objectPath := filepath.Join(domain.BucketsPath, bucketname, objectname)
	return os.Open(objectPath)
}

func (repo ObjectCSV) DeleteObject(bucketname, objectname string) error {
	objectPath := filepath.Join(domain.BucketsPath, bucketname, objectname)
	if err := os.Remove(objectPath); err != nil {
		return err
	}

	now := time.Now()
	if err := csvparser.DeleteRecord(filepath.Join(domain.BucketsPath, bucketname, "objects.csv"), objectname); err != nil {
		return err
	}

	return csvparser.ReWriteCSV(filepath.Join(domain.BucketsPath, "buckets.csv"), []string{bucketname, "", now.Format(time.ANSIC), domain.ActiveMark})
}

func (repo ObjectCSV) IsObjectExist(csvPath, name string) (bool, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if len(record) >= 4 && record[0] == name {
			return true, nil
		}
	}
	return false, nil
}
