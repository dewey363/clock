package clock

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//ExampleClock_AddJobRepeat 基于函数回调，一个对重复任务的使用演示。
func ExampleClock_AddJobRepeat() {
	var (
		myClock   = NewClock()
		counter   = 0
		mut       sync.Mutex
		sigalChan = make(chan struct{}, 0)
	)
	fn := func() {
		fmt.Println("schedule repeat")
		mut.Lock()
		defer mut.Unlock()
		counter++
		if counter == 3 {
			sigalChan <- struct{}{}
		}

	}
	//create a task that executes three times,interval 50 millisecond
	event, inserted := myClock.AddJobRepeat(time.Duration(time.Millisecond*50), 0, fn)
	if !inserted {
		log.Println("failure")
	}

	//等待阻塞信号
	<-sigalChan
	myClock.DelJob(event)

	//wait a second,watching
	time.Sleep(time.Second)
	//Output:
	//
	//schedule repeat
	//schedule repeat
	//schedule repeat
}

//ExampleClock_AddJobRepeat2 ，基于函数回调，演示添加有次数限制的重复任务
//  执行3次之后，撤销定时事件
func ExampleClock_AddJobRepeat2() {
	var (
		myClock = NewClock()
	)
	//define a repeat task
	fn := func() {
		fmt.Println("schedule repeat")
	}
	//add in clock,execute three times,interval 200 millisecond
	_, inserted := myClock.AddJobRepeat(time.Duration(time.Millisecond*200), 3, fn)
	if !inserted {
		log.Println("新增事件失败")
	}
	//wait a second,watching
	time.Sleep(time.Second)
	//Output:
	//
	//schedule repeat
	//schedule repeat
	//schedule repeat
}

//ExampleClock_AddJobWithTimeout 基于函数回调，对一次性任务正常使用的演示。
func ExampleClock_AddJobWithTimeout() {
	var (
		jobClock = NewClock()
		jobFunc  = func() {
			fmt.Println("schedule once")
		}
	)
	//add a task that executes once,interval 100 millisecond
	jobClock.AddJobWithInterval(time.Duration(100*time.Millisecond), jobFunc)

	//wait a second,watching
	time.Sleep(1 * time.Second)

	//Output:
	//
	//schedule once
}

//ExampleClock_AddJobWithDeadtime 基于事件提醒，对一次性任务中途放弃的使用演示。
func ExampleClock_AddJobWithDeadtime() {
	var (
		myClock = NewClock()
		jobFunc = func() {
			fmt.Println("schedule once")
		}
		actionTime = time.Now().Add(time.Millisecond * 500)
	)
	//创建一次性任务，定时500ms
	job, _ := myClock.AddJobWithDeadtime(actionTime, jobFunc)

	//任务执行前，撤销任务
	time.Sleep(time.Millisecond * 300)
	myClock.DelJob(job)

	//等待2秒，正常情况下，事件不会再执行
	time.Sleep(2 * time.Second)

	//Output:
	//
	//
}
