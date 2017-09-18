package twitter

type TweetImageMetaData struct {
	Id   int64
	Url  string
	Text string
}

func (t *TweetImageMetaData) GetId() int64 {
	return t.Id
}

func (t *TweetImageMetaData) GetUrl() string {
	return t.Url
}

func (t *TweetImageMetaData) GetText() string {
	return t.Text
}

type MetaData interface {
	GetId() int64
	GetUrl() string
	GetText() string
}

type MetaDataSet struct {
	MediaType string
	Source    string
	List      []MetaData
}
