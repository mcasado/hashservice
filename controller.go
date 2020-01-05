package main

// Controller passes state to our handlers
type Controller struct {
	storage HashStorage
}

// New is a Controller 'constructor'
func NewController(storage HashStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}
