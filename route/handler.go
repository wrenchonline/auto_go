package route

import (
	"context"
	pb "evenkey/mygrpc" // 导入生成的 protobuf 代码
	"io"
	"strconv"
	"time"

	"github.com/google/martian/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	port     int
	addr     string
	Username string
	Password string
	Token    string
}

func (c *Client) Run() {
	//var WG sync.WaitGroup //当前与jackdaw等待同步计数
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	address := c.addr + ":" + strconv.Itoa(c.port)
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

		}
	}()

}
