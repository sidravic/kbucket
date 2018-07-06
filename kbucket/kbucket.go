package kbucket

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/sidravic/kbucket/util"
)

type KBucket struct {
	Root                  *Node
	NumberOfNodesToPing   int
	NumberOfNodePerBucket int
	LocalNodeID           []byte
}

type Contact struct {
	ContactID []byte
	IP        string
	Port      int
}

type Node struct {
	ID        []byte
	Contacts  []*Contact
	DontSplit bool
	Left      *Node
	Right     *Node
	IP        string
	Port      int
	Leaf      bool
}

func NewNode() *Node {
	return &Node{
		ID:        util.GenerateRandomID(),
		Contacts:  []*Contact{},
		DontSplit: false,
		Left:      nil,
		Right:     nil,
		Leaf:      true,
	}
}

func NewContact(ip string, port int) *Contact {
	return &Contact{
		ContactID: util.GenerateRandomID(),
		IP:        ip,
		Port:      port,
	}
}

func GenerateRandomID() []byte {
	s := sha1.New()
	time := time.Now().String()
	s.Write([]byte(time))
	return s.Sum(nil)
}

func NewKBucket(localNodeIP string, localNodePort int) *KBucket {
	localNode := NewNode()
	localNode.IP = localNodeIP
	localNode.Port = localNodePort

	return &KBucket{
		Root: localNode,
		NumberOfNodePerBucket: 4,
		NumberOfNodesToPing:   3,
		LocalNodeID:           localNode.ID,
	}
}

func (k *KBucket) Add(c *Contact) {
	node := k.Root
	bitIndex := 0

	fmt.Println(reflect.TypeOf(node.Contacts))

	for node.Leaf == false {
		bitIndex++
		node = determineNode(node, c.ContactID, bitIndex)
	}

	contactIndexInNode := k.ContactExists(node, c)
	if contactIndexInNode != -1 {
		/*
			Call update
		*/
		k.Update(node, contactIndexInNode, c)
	}

	if len(node.Contacts) < k.NumberOfNodePerBucket {
		node.Contacts = append(node.Contacts, c)
	} else {
		k.split(node, bitIndex)
		fmt.Println("0---------------------------------------")
		fmt.Println(fmt.Sprintf("Node split at %s", c.IP))
		k.Add(c)
	}
}

func (k *KBucket) ContactExists(node *Node, contact *Contact) (indexOfContact int) {
	indexOfContact = -1

	for i := 0; i < len(node.Contacts); i++ {
		if bytes.Equal(node.Contacts[i].ContactID, contact.ContactID) {
			indexOfContact = i
			return
		}
	}
	return
}

func (k *KBucket) Update(node *Node, contactIndex int, c *Contact) {
	incumbentContact := node.Contacts[contactIndex]

	if bytes.Equal(incumbentContact.ContactID, c.ContactID) {
		node.Contacts = append(node.Contacts[:contactIndex], node.Contacts[contactIndex+1:]...)
		node.Contacts = append(node.Contacts, c)
		fmt.Println("0---------------------------------------")
		fmt.Println("Updated Node with IP", c.IP)
	}

}

func (k *KBucket) split(node *Node, bitIndex int) {
	node.Left = NewNode()
	node.Right = NewNode()

	for i := 0; i < len(node.Contacts); i++ {
		selectedNode := determineNode(node, node.Contacts[i].ContactID, bitIndex)
		selectedNode.Contacts = append(selectedNode.Contacts, node.Contacts[i])
	}

	node.Contacts = node.Contacts[:0]
	node.Leaf = false
	rootNodePresent := false
	leftNode := node.Left

	for i := 0; i < len(node.Left.Contacts); i++ {
		if string(leftNode.Contacts[i].ContactID[:]) == string(k.LocalNodeID[:]) {
			rootNodePresent = true
			break
		}
	}

	if rootNodePresent {
		node.Left.DontSplit = false
		node.Right.DontSplit = true
	}

	return
}

func determineNode(node *Node, contactID []byte, bitIndex int) *Node {
	bytesDescribedByBitIndex := bitIndex / 8
	bitIndexWithinTheByte := bitIndex % 8

	byteUnderConsideration := contactID[bytesDescribedByBitIndex]
	bitSequenceWithNthBitSet := math.Pow(2, 7-float64(bitIndexWithinTheByte))
	if (int64(byteUnderConsideration) & int64(bitSequenceWithNthBitSet)) != 0 {
		return node.Left
	} else {
		return node.Right
	}
}
