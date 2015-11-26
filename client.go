package main 
import (
"fmt"
"net/http"
"hash/crc32"
"sort"
"strconv"
"sync"
"strings"
"log"
)
type Ring []uint32
type Data struct {
	Key int 
	Value string
}
type DataPair struct {
	DataPair[]Data
}

// Len returns the length of the uints array.
func (x Ring) Len() int { 
	return len(x)
	 }

// Less returns true if element i is less than element j.
func (x Ring) Less(i, j int) bool {
 return x[i] < x[j] 
}

// Swap exchanges elements i and j.
func (x Ring) Swap(i, j int) {
 x[i], x[j] = x[j], x[i]
  }

type Node struct {
	ID int
	IP string
	Weight int
}
 func CreateNode(id int, ip string, weight int) *Node {
 	return &Node {
 		ID: id,
 		IP: ip,
 		Weight: weight,
 	}
 }
 type Consistent struct {
 	Nodes map[uint32]Node
 	NumberOfReplicas int
 	Mem map[int]bool
 	sortedRing Ring
 	sync.RWMutex
 }

 func New() *Consistent {
 	c := new(Consistent)
 	c.NumberOfReplicas = 10
 	c.Nodes = make(map[uint32]Node)
 	c.Mem = make(map[int]bool)
 	c.sortedRing = Ring{}
 	return c
 }

 func (c *Consistent) Add(node *Node) bool{
 	c.Lock()
 	defer c.Unlock()
 	_,val:= c.Mem[node.ID]
 	if val {
 		return false
 	}
 	count:= c.NumberOfReplicas * node.Weight
 	for i := 0; i < count; i++ {
 		string := c.merge(i, node)
 		c.Nodes[c.hashstring(string)] = *(node)
 	}
 	c.Mem[node.ID] = true
 	c.Ringsort()
 	return true
 }
func (c *Consistent) Ringsort() {
	c.sortedRing = Ring{}
	for k:= range c.Nodes {
		c.sortedRing = append(c.sortedRing, k)
	}
	sort.Sort(c.sortedRing)
}

func (c *Consistent) merge(i int, node *Node) string{
	return node.IP+ "*" + strconv.Itoa(node.Weight) +
	             "-" + strconv.Itoa(i) + "-" + 
	             strconv.Itoa(node.ID)
}

func (c *Consistent) hashstring(input string) uint32{
	return crc32.ChecksumIEEE([]byte(input))
}

func (c *Consistent) search(hash uint32) int{
	i:= sort.Search(len(c.sortedRing), func (i int) bool{
		return c.sortedRing[i] >= hash 
	})
	if i < len(c.sortedRing) {
		if i== len(c.sortedRing)-1 {
              return 0
		} else{
			return i
		}
	} else {
		return len(c.sortedRing)-1
	}
}

func (c *Consistent) Get(key string) Node {
	c.RLock()
	defer c.RUnlock()

	hash := c.hashstring(key)
	i:= c.search(hash)
	return c.Nodes[c.sortedRing[i]]
}

func Put(data Data, trackserver string) bool{
	url:= "http://localhost:"
	url = url+ trackserver
	url = url + "/keys/" + strconv.Itoa(data.Key) + "/" + data.Value
	fmt.Println(url)
	client := &http.Client{}
	req, err:= http.NewRequest("PUT", url, nil)
	resp,err := client.Do(req)
	if err!=nil {
		log.Fatal(err)
		fmt.Println("wrong PUT", err)
		panic(err)
	}
	if resp !=nil {
		return true
	}else {
		return false
	}
}


func main() {
	consish := New()
    consish.Add(CreateNode(0, "http://localhost:3000", 1))
    consish.Add(CreateNode(1, "http://localhost:3001", 1))
    consish.Add(CreateNode(2, "http://localhost:3002", 1))

    newmap := make(map[Data]string)
    var sharddata []Data
    sharddata = append(sharddata, Data{1, "a"})
    sharddata = append(sharddata, Data{2, "b"})
    sharddata = append(sharddata, Data{3, "c"})
    sharddata = append(sharddata, Data{4, "d"})
    sharddata = append(sharddata, Data{5, "e"})
    sharddata = append(sharddata, Data{6, "f"})
    sharddata = append(sharddata, Data{7, "g"})
    sharddata = append(sharddata, Data{8, "h"})
    sharddata = append(sharddata, Data{9, "i"})
    sharddata = append(sharddata, Data{10, "j"})


    for i := 0; i < 10; i++ {
    	k := consish.Get(sharddata[i].Value)
    	newmap[sharddata[i]] = k.IP
    }
    
    for k, v:= range newmap{
    	if strings.Contains(v, "3000") {
    		Put(k, "3000")
    	}else if strings.Contains(v, "3001") {
    		Put(k, "3001")
    	}else if strings.Contains(v, "3002") {
    		Put(k, "3002")
    	}
    }

}





















