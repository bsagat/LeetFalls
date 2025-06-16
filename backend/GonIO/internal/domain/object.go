package domain

type Object struct {
	ObjectKey    string `xml:"Objectkey"`
	Size         string `xml:"Size"`
	ContentType  string `xml:"Content-type"`
	LastModified string `xml:"Last_Modified_Time"`
}

type ObjectsList struct {
	Objects []Object `xml:"Object"`
}
