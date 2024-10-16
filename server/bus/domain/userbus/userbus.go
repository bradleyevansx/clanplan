package userbus

type Storer interface {

}

type Business struct {
	storer Storer
}