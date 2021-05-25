package main

const NULL = "\\0\\"

type Deleted int

const (
	NO   Deleted = 0
	SOFT Deleted = 1
	HARD Deleted = 2
)

type Node struct {
	nvolumes []string
	deleted  Deleted
	id       string
	prev     string
	next     string
	current  string
}

type Nodes struct {
	L []map[string]string
}
