package dal

import (
	"GonIO/internal/domain"
	csvparser "GonIO/pkg/myCSV"
	"encoding/csv"
	"os"
	"time"
)

type BucketCSV struct {
	BucketMetaPath string
}

var _ domain.BucketDal = (*BucketCSV)(nil)

func NewBucketXMLRepo() *BucketCSV {
	return &BucketCSV{BucketMetaPath: domain.BucketsMetaPath}
}

func (repo BucketCSV) GetBucketList() ([]domain.Bucket, error) {
	file, err := os.Open(repo.BucketMetaPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var buckets []domain.Bucket
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 4 {
			continue
		}
		if record[3] == domain.DeletionMark {
			continue
		}
		buckets = append(buckets, domain.Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		})
	}
	return buckets, nil
}

func (repo BucketCSV) IsUniqueBucket(bucketName string) (bool, error) {
	metaFile, err := os.OpenFile(domain.BucketsMetaPath, os.O_RDWR, 0o666)
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
		if records[i][0] == bucketName && records[i][3] == domain.ActiveMark {
			return false, nil
		}
	}
	return true, nil
}

func (repo BucketCSV) CreateBucket(bucketname string) error {
	// Bucket directory creating
	err := os.Mkdir(domain.BucketsPath+"/"+bucketname, 5580)
	date := time.Now()
	if err != nil {
		return err
	}

	// Updating bucket metadata
	bucketMetafile, err := os.OpenFile(domain.BucketsMetaPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer bucketMetafile.Close()

	data := []string{bucketname, date.Format(time.ANSIC), date.Format(time.ANSIC), "Active"}
	err = csvparser.WriteCSV(bucketMetafile, data)
	if err != nil {
		return err
	}

	// Create object Metadata file
	file, err := os.Create(domain.BucketsPath + "/" + bucketname + "/objects.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	// Updating object metadata
	data = []string{"ObjectKey", "Size", "ContentType", "LastModified"}
	err = csvparser.WriteCSV(file, data)
	if err != nil {
		return err
	}
	return nil
}

func (repo BucketCSV) DeleteBucket(bucketName string) error {
	empty, err := csvparser.CheckEmpty(domain.BucketsPath + "/" + bucketName + "/objects.csv")
	if err != nil {
		return err
	}
	if !empty {
		return domain.ErrBucketIsNotEmpty
	}

	err = os.RemoveAll(domain.BucketsPath + "/" + bucketName)
	if err != nil {
		return err
	}

	date := time.Now()
	err = csvparser.ReWriteCSV(domain.BucketsPath+"/buckets.csv", []string{bucketName, "", date.Format(time.ANSIC), domain.DeletionMark})
	if err != nil {
		return err
	}
	return nil
}
