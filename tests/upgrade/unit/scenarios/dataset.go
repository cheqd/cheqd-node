package scenarios

type IDataSet interface {
	WriteExisting() error
	CheckExpected() error
}
