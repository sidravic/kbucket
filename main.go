package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/sidravic/kbucket/kbucket"
)

func main() {
	k := kbucket.NewKBucket("127.0.0.1", 3000)
	contact1 := kbucket.NewContact("127.0.0.2", 3000)
	k.Add(contact1)
	contact2 := kbucket.NewContact("127.0.0.3", 3000)
	k.Add(contact2)
	contact3 := kbucket.NewContact("127.0.0.4", 3000)
	k.Add(contact3)
	contact4 := kbucket.NewContact("127.0.0.5", 3000)
	k.Add(contact4)
	contact5 := kbucket.NewContact("127.0.0.6", 3000)
	k.Add(contact5)
	fmt.Println(spew.Sdump(k))
	contact6 := kbucket.NewContact("127.0.0.4", 3000)
	k.Add(contact6)

	spew.Dump(k.FindClosest(contact6, 3))
}
