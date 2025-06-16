package dal

import (
	"GonIO/internal/domain"
	csvparser "GonIO/pkg/myCSV"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ObjectCSV struct {
}

func NewObjectCSVRepo() *ObjectCSV {
	return &ObjectCSV{}
}

var _ domain.ObjectDal = (*ObjectCSV)(nil)

func (repo ObjectCSV) List_Object(bucketname string) (domain.ObjectsList, error) {
	empty := domain.ObjectsList{}
	metafile, err := os.OpenFile(domain.BucketsPath+"/"+bucketname+"/objects.csv", os.O_RDWR, 0o666)
	if err != nil {
		return empty, err
	}
	defer metafile.Close()

	reader := csv.NewReader(metafile)
	records, err := reader.ReadAll()
	if err != nil {
		return empty, err
	}

	var Objects []domain.Object
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			continue
		}
		Objects = append(Objects, domain.Object{
			ObjectKey:    record[0],
			Size:         record[1],
			ContentType:  record[2],
			LastModified: record[3],
		})
	}

	objectlist := domain.ObjectsList{Objects: Objects}
	return objectlist, nil
}

func (repo ObjectCSV) UploadObject(bucketname, objectname string, r *http.Request) error {
	file, err := os.OpenFile(domain.BucketsPath+"/"+bucketname+"/"+objectname, os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, r.Body); err != nil {
		return err
	}

	date := time.Now()
	Metafile, err := os.OpenFile(domain.BucketsPath+"/"+bucketname+"/objects.csv", os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer Metafile.Close()

	filetype := r.Header.Get("Content-Type")
	data := []string{objectname, strconv.Itoa(int(r.ContentLength)), filetype, date.Format(time.ANSIC)}
	defer file.Close()

	if err = csvparser.WriteCSV(Metafile, data); err != nil {
		return err
	}

	err = csvparser.ReWriteCSV(domain.BucketsMetaPath, []string{bucketname, "", date.Format(time.ANSIC), domain.ActiveMark})
	if err != nil {
		return err
	}
	return nil
}

func (repo ObjectCSV) RetrieveObject(bucketname, objectname string, w http.ResponseWriter) error {
	objectFile, err := os.Open(domain.BucketsPath + "/" + bucketname + "/" + objectname)
	if err != nil {
		return err
	}
	defer objectFile.Close()

	if _, err = io.Copy(w, objectFile); err != nil {
		return err
	}
	return nil
}

func (repo ObjectCSV) DeleteObject(bucketname, objectname string) error {
	date := time.Now()
	if err := os.Remove(domain.BucketsPath + "/" + bucketname + "/" + objectname); err != nil {
		return err
	}

	if err := csvparser.ReWriteCSV(domain.BucketsPath+"/buckets.csv", []string{bucketname, "", date.Format(time.ANSIC), domain.ActiveMark}); err != nil {
		return err
	}

	if err := csvparser.DeleteRecord(domain.BucketsPath+"/"+bucketname+"/objects.csv", objectname); err != nil {
		return err
	}
	return nil

}

func (repo ObjectCSV) IsObjectExist(path, name string) (bool, error) {
	metaFile, err := os.OpenFile(path, os.O_RDWR, 0o666)
	if err != nil {
		return false, err
	}
	defer metaFile.Close()

	reader := csv.NewReader(metaFile)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	for i := 0; i < len(records); i++ {
		if len(records[i]) < 4 {
			continue
		}
		if records[i][0] == name {
			return true, nil
		}
	}
	return false, nil
}
