package domain

var (
	DeletionMark = "Marked for deletion"
	ActiveMark   = "Active"
)

type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

// Buckets list for XML response
type BucketList struct {
	Buckets []Bucket `xml:"Buckets"`
}
