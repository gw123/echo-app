package echoapp

//抽象的module 建议写在外面
type ExampleService interface {
	GetTime() string
}
