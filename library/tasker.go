package library

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

type Tasker struct {
	sysInterval    time.Duration `任务执行间隔`
	sysTicker      *time.Ticker  `任务定时器`
	sysStatus      int           `任务启动状态，0-未启动，1-已启动(可监听控制命令，非程序是否在执行）, 2-已暂停`
	sysControlChan chan int      `控制任务的管道，0-暂停任务，1-启动任务，2-停止`

	runTimeout time.Duration             `程序执行超时时间`
	runFunc    func(time.Duration) error `程序执行函数`
	runMutex   *sync.Mutex               `程序执行状态锁,防止套圈；重启时等待协程退出`
}

func NewTasker(sysInterval, runTimeout time.Duration, runFunc func(time.Duration) error) (*Tasker, error) {
	t := new(Tasker)
	if err := t.init(sysInterval, runTimeout, runFunc); err == nil {
		return t, err
	}
	return nil, errors.New("Init error!")
}

//初始化tasker
func (this *Tasker) init(sysInterval, runTimeout time.Duration, runFunc func(time.Duration) error) error {
	this.sysInterval = sysInterval
	this.sysTicker = time.NewTicker(this.sysInterval)
	this.sysStatus = 0
	this.sysControlChan = make(chan int)

	this.runTimeout = runTimeout
	this.runFunc = runFunc
	this.runMutex = new(sync.Mutex)
	return nil
}

//run函数
func (this *Tasker) Run() (err error) {
	if this.sysStatus != 0 {
		err = errors.New("task is running!")
		return
	} else {
		this.sysStatus = 1
	}
	for {
		select {
		case sign := <-this.sysControlChan:
			switch sign {
			case 0:
				// pause
				//只有在启动状态的可以暂停
				if this.sysStatus == 1 {
					this.runMutex.Lock()
					this.sysTicker.Stop()
					this.sysStatus = 2
					this.runMutex.Unlock()
				}
			case 1:
				// start
				//只有在暂停状态的可以重新开始
				if this.sysStatus == 2 {
					this.sysTicker = time.NewTicker(this.sysInterval)
					this.sysStatus = 1
				}
			case 2:
				// stop
				this.runMutex.Lock()
				this.sysTicker.Stop()
				close(this.sysControlChan)
				this.sysStatus = 0
				this.runMutex.Unlock()
				runtime.Goexit()
				return
			}
		case <-this.sysTicker.C:
			go func() {
				//使用了互斥锁，注意加锁和解锁
				this.runMutex.Lock()
				err = this.runFunc(this.runTimeout)
				this.runMutex.Unlock()
			}()
		}
	}
	return
}

//pause函数
func (this *Tasker) Pause() (err error) {
	switch this.sysStatus {
	case 0:
		err = errors.New("sys is not running")
	case 1:
		this.sysControlChan <- 0
	case 2:
		err = errors.New("sys was Paused already")
	}
	return
}

//start函数
func (this *Tasker) Start() (err error) {
	switch this.sysStatus {
	case 0:
		err = errors.New("sys is not running")
	case 1:
		err = errors.New("sys was Started already")
	case 2:
		this.sysControlChan <- 1
	}
	return
}

//stop函数
func (this *Tasker) Stop() (err error) {
	switch this.sysStatus {
	case 0:
		err = errors.New("sys is not running")
	case 1, 2:
		this.sysControlChan <- 2
	}
	return
}
