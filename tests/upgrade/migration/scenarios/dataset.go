package scenarios


type IDataSet interface {
	Prepare() error
	Validate() error
}