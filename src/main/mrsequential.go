package main

// // go run mrsequential.go wc.so pg*.txt
import (
	"dist-systems/mr"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"sort"
)

type ByKey []mr.KeyValue

func (a ByKey) Len() int {return len(a)}
func (a ByKey) Swap(i, j int) {a[i], a[j] = a[j], a[i]}
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func loadPlugin(filename string)(func(string, string) []mr.KeyValue, func(string, []string) string){

    p, err := plugin.Open(filename)
    if err != nil {
        log.Fatalf("%v", err)
    }
    xmapf, err := p.Lookup("Map")
    if err!=nil {
        log.Fatalf("cannot find map in %v", filename)
    }
    mapf := xmapf.(func(string, string) []mr.KeyValue)

    xreducef, err := p.Lookup("Reduce")
    if err!=nil {
        log.Fatalf("cannot find reduce in %v", filename)
    }
    reducef:= xreducef.(func(string, []string) string)
    return mapf, reducef

}

func main(){

    if len(os.Args)<3{
        fmt.Fprintf(os.Stderr,"Usage go run mrsequential.go xx xxx.txt")
        os.Exit(1)
    }
    mapf, reducef := loadPlugin(os.Args[1])

    intermediate := []mr.KeyValue{}

    for _, filename := range os.Args[2:] {
        fmt.Println("Processing file : ", filename)
        file, err := os.Open(filename)
        if err != nil {
            log.Fatalf("cannot open filename %s", filename)
        }
        defer file.Close()
        content, err := ioutil.ReadAll(file)
        if err != nil {
            log.Fatalf("Error reading the filename %s", filename)
        }
        kv := mapf(filename, string(content))
        intermediate=append(intermediate, kv...)
    }
    sort.Sort(ByKey(intermediate))
    // fmt.Fprintf(os.Stderr, "%v", intermediate)

    ofile, _ := os.Create("mr-out-0")
    defer ofile.Close()

    i := 0
    for i < len(intermediate) {
        j := i+1

        for j <len(intermediate) && intermediate[i].Key == intermediate[j].Key{
            j++
        }
        values := []string{}
        for k:=i; k<j; k++ {
            values = append(values, string(intermediate[k].Value))
        }
        output := reducef(intermediate[i].Key,values)
        fmt.Fprintf(ofile, "Found %s key in %v time\n", intermediate[i].Key, output)
        i = j
    }
}
