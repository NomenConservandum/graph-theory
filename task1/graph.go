package main

// граф: ориентированный / неориентированный, взвешанный / невзвешанный
// методы внутри и снаружи графа
// методы чтения из файла и в файл


type node[T any] struct {
	Id uint16
	Value T
	Connection *connection
}

type connection struct {
	List []*node[T]
}
