package route

import (
	"context"
	pb "evenkey/mygrpc" // 导入生成的 protobuf 代码
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/google/martian/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Port     int
	Addr     string
	Username string
	Password string
	token    string
	Keydelay int
	state    int
}

const (
	Start int = 1
)

func (c *Client) Run() {
	robotgo.KeySleep = 10

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	address := c.Addr + ":" + strconv.Itoa(c.Port)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Errorf("fail to dial: %v", err)
	}

	defer conn.Close()
	client := pb.NewMyServiceClient(conn)
	ctx := context.Background()
	Authctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	ret, err := client.Authenticate(Authctx, &pb.AuthRequest{Username: c.Username, Password: c.Password})
	if err != nil {
		panic(err)
	}
	if !ret.Success {
		log.Errorf("用户名密码错误")
		return
	}

	stream, err := client.PressKey(ctx)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}
	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				close(waitc)
				log.Errorf("routeChat error %s", err.Error())
				return
			}
			if !in.Success {
				log.Errorf("通讯授权错误")
			}

			fmt.Println(in.Key)
			//执行代码
			// robotgo.KeyDown(in.Key)
			// robotgo.KeyUp(in.Key)

		}
	}()

	for {
		select {
		case <-waitc:
			log.Infof("结束通讯")
		default:
		}
		// robotgo.KeySleep = 1
		// robotgo.EventHook(hook.KeyDown, []string{"b"}, func(e hook.Event) {
		// 	robotgo.EventEnd()
		// })
		// s := robotgo.EventStart()
		// <-robotgo.EventProcess(s)
		// ok := robotgo.AddEvents("b")
		// if ok {
		// 	c.state = Start
		// }
	}

}
