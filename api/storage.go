package api

type StorageStatusResponse interface  {
	error
	StatusCode() int

}

type Storage interface  {
	Create(string,interface{}) (string,StorageStatusResponse)
	Get(string,string) (interface{},StorageStatusResponse)
	GetAll(string) ([]interface{},StorageStatusResponse)
}
