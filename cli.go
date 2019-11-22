package main

import (
	"fmt"
	pb "github.com/karlhjm/consignment-service/proto/consignment"
	"io/ioutil"
	"encoding/json"
	"errors"
	micro "github.com/micro/go-micro"
	//"github.com/micro/go-micro/config/cmd"
	"os"

	//microclient "github.com/micro/go-micro/client"
	"log"
	//"os"
	"context"

)

const (
	DEFAULT_INFO_FILE = "consignment.json"
)

// 读取 consignment.json 中记录的货物信息
func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error")
	}
	return consignment, nil
}

func main() {


	service := micro.NewService(micro.Name("go.micro.srv.consignment"))
	service.Init()

	client := pb.NewShippingService("go.micro.srv.consignment", service.Client())


	// 在命令行中指定新的货物信息 json 文件
	infoFile := DEFAULT_INFO_FILE
	if len(os.Args) > 1 {
		infoFile = os.Args[1]
	}

	// 解析货物信息
	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 调用 RPC
	// 将货物存储到我们自己的仓库里
	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	// 新货物是否托运成功
	fmt.Println("created: ", resp.Created)


	resp, err = client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}

	for _, c := range resp.Consignments {
		fmt.Println(c)
	}
}