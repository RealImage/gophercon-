package models

type QueueNode struct {
	Person Actor
	Role   string
	Path   []RelationNode
	IsDest bool
}

type RelationNode struct {
	Degree       int
	Movie        string
	FrstPerson   Actor
	SecondPerson Actor
}

type QueueReader struct {
	Queue      []QueueNode
	DestFound  bool
	ResultNode QueueNode
}
