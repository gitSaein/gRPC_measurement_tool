package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	c "gRPC_measurement_tool/cmd"
	errorModel "gRPC_measurement_tool/error"
	m "gRPC_measurement_tool/measure"
	u "gRPC_measurement_tool/util"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name    string
	option  m.Option
	startAt time.Time
)

// 프로그램 실행시 호출
func init() {
	startAt = time.Now()
	option = c.Basic()
}

func connectServer(wait *sync.WaitGroup, cmd m.Option, report *m.Report) {
	pid, opts, err, ctx, cancel := u.SetOption(cmd, startAt, report)
	errorModel.HandleReponse(err, pid, report, cmd, m.SetOption, startAt)
	defer wait.Done()
	defer cancel()

	// Set up a connection to the server.
	startAt = time.Now()
	conn, err := grpc.DialContext(ctx, cmd.Target, opts...)
	errorModel.HandleReponse(err, pid, report, cmd, m.Dial, startAt)

	// go u.CheckDialConnection(conn, ctx, pid, startAt, report)

	if conn != nil {
		defer func() {
			err = conn.Close()
			errorModel.HandleReponse(err, pid, report, cmd, m.Dial, startAt)
		}()

		startAt = time.Now()
		reply, err := u.CallMethod(cmd, conn, ctx)
		errorModel.HandleReponse(err, pid, report, cmd, m.CallMethod, startAt)
		log.Printf("message: %v", reply.GetMessage())
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	fmt.Println(runtime.GOMAXPROCS(0))   // 설정 값 출력

	report := &m.Report{}

	wg := new(sync.WaitGroup)
	wg.Add(option.Tr)

	for i := 0; i < option.Tr; i++ {
		go connectServer(wg, option, report)
	}
	wg.Wait() //Go루틴 모두 끝날 때까지 대기

	defer func() {
		report.Total = time.Since(startAt)
		m.PrintResult(report, option)
	}()

}
