package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)



type Task struct {
    id int
    state string
    done bool
    infile string
    outfile string

}


type WorkerMeta struct {
    id int
}

type Coordinator struct {
    maptasks []Task
    done bool
    maptasksdone bool
    workers []WorkerMeta
    reducetasks []Task
	// Your definitions here.
}
// Your code here -- RPC handlers for the worker to call.
func (c * Coordinator) GetTask(args *GetTaskArgs, reply *GetTaskReplyArgs) error {
    
    if !c.maptasksdone{
        task := c.maptasks[len(c.maptasks)-1]
        reply.Infile = task.infile
        reply.Operation = "map" 
        return nil
    }

    task := c.reducetasks[len(c.reducetasks)-1]
    reply.Infile = task.infile
    reply.Operation = "reduce"
    return nil
}
func (c *Coordinator) UpdateTaskState(args, reply){


}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//

func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil 
}
//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//


func (c *Coordinator) Done() bool {
	ret := false

    // Your code here
    // for _, task := range c.processedtasks {
    //     if task.done == false{
    //         return false 
    //     }
    // }
    
    return ret 
    }

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//

func remove[T any](s []T, i int)[]T{
    copy(s[i:], s[i+1:])
    return s[:len(s)-1]
}



func (c *Coordinator) popTaskAndMove(id int) {
    for i, task := range c.maptasks{
        if task.id == id {
            c.maptasks = remove(c.maptasks, i)
            c.reducetasks = append(c.reducetasks, task)
    }

}

func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
    for i, file := range files{
        c.maptasks = append(c.maptasks, Task{id: i, done: false, state: "stage",infile:file})
    }
	// Your code here.
	c.server()
	return &c
}
