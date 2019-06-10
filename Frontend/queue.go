package main

type ServerQueue struct{
	queue []BEServer
	size int
}

func (q *ServerQueue) Push(s BEServer) {
	q.queue = append(q.queue,s)
	q.size++
}

func (q *ServerQueue) Pop() {
	q.queue = q.queue[1:]
	q.size--
}

func (q *ServerQueue) Front() BEServer {
	return q.queue[0]
}

func (q *ServerQueue) Empty() bool {
	return q.size == 0
}