package main

import (
	"bytes"
	"encoding/gob"
	"strings"
)

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

func getNext(tokens *[]string, index int) string {
	if (index + 1) < len(*tokens) {
		return (*tokens)[index+1]
	}
	return NULL
}

func getPrev(tokens *[]string, index int) string {
	if (index - 1) > 0 {
		return (*tokens)[index-1]
	}
	return NULL
}

func toNodes(data []byte) Nodes {
	buffer := bytes.NewReader(data)
	dec := gob.NewDecoder(buffer)
	nodes := Nodes{L: make([]map[string]string, 0)}

	if err := dec.Decode(&nodes); err != nil {
		panic("error: cannot convert nodes into bytes")
	}
	return nodes
}

func fromNodes(nodes Nodes) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(nodes); err != nil {
		panic("error: cannot convert nodes into bytes")
	}
	return buf.Bytes()
}

func NSplit(seq []byte, size int) []string {
	sep := " "
	s := string(seq)
	tks := strings.Split(s, sep)
	tokenizedSequence := []string{}
	for i := 0; i < len(tks); i++ {
		tmp := ""
		for c := 0; c < size; c++ {
			if i+c < len(tks) {
				tmp += tks[i+c] + " "
			} else {
				break
			}

		}
		tokenizedSequence = append(tokenizedSequence, string(tmp))
	}
	return tokenizedSequence
}

// func main() {
// 	s := "In this chapter we shall discuss one of the most remarkable and amusing consequences of mechanics, the behavior of a rotating wheel. In order to do this we must first extend the mathematical formulation of rotational motion, the principles of angular momentum, torque, and so on, to three-dimensional space. We shall not use these equations in all their generality and study all their consequences, because this would take many years, and we must soon turn to other subjects."
// 	sarr := NSplit([]byte(s), 4)
// 	tokens := make(map[string][]Node)
// 	xstr := ""
// 	t, _ := uuid.NewUUID()
// 	// ar := []map[string]string{}
// 	id := t.String()
// 	marr := []map[string]string{}
// 	for index, value := range sarr {
// 		n := Node{nvolumes: []string{}, deleted: HARD, current: value, id: id, next: getNext(&sarr, index), prev: getPrev(&sarr, index)}
// 		tokens[id] = append(tokens[id], n)

// 		// create bytes
// 		var buf bytes.Buffer
// 		enc := gob.NewEncoder(&buf)

// 		m := make(map[string]string)

// 		m["nvolumes"] = strings.Join(n.nvolumes, ",")
// 		m["deleted"] = string(n.deleted)
// 		m["current"] = n.current
// 		m["id"] = n.id
// 		m["next"] = n.next
// 		m["prev"] = n.prev

// 		marr = append(marr, m)

// 		if err := enc.Encode(m); err != nil {
// 			panic("oops")
// 		}
// 		// fmt.Println(buf.Bytes())
// 		// fmt.Println()

// 		input := buf.Bytes()
// 		bff := bytes.NewBuffer(input)
// 		dec := gob.NewDecoder(bff)
// 		xx := make(map[string]string)

// 		if err := dec.Decode(&xx); err != nil {
// 			panic("boi")
// 		}
// 		// fmt.Println(xx["nvolumes"])
// 		// fmt.Println(xx["deleted"])
// 		// fmt.Println(xx["current"])
// 		// fmt.Println(xx["id"])
// 		// fmt.Println(xx["next"])
// 		// fmt.Println(xx["prev"])
// 		// fmt.Println()
// 		// fmt.Println()
// 		// fmt.Println()
// 		// fmt.Println()

// 	}

// 	var buf bytes.Buffer
// 	enc := gob.NewEncoder(&buf)
// 	w := Nodes{L: marr}
// 	if err := enc.Encode(w); err != nil {
// 		panic("oops")
// 	}
// 	fmt.Println(buf.Bytes())
// 	// input -> buffer -> decoder -> Decode(schema) -> output
// 	input := buf.Bytes()
// 	buffer := bytes.NewReader(input)
// 	dec := gob.NewDecoder(buffer)
// 	ww := Nodes{L: make([]map[string]string, 0)}

// 	if err := dec.Decode(&ww); err != nil {
// 		panic("boi1")
// 	}
// 	fmt.Println(ww)

// 	// for _, value := range tokens {
// 	// 	fmt.Println(value)
// 	// }
// 	fmt.Println(xstr)

// }
