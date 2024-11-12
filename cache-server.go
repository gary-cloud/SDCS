package main

import (
    "fmt"
    "log"
    "net"
    "flag"
    "time"
    "context"
    "net/http"
    "hash/fnv"
	"encoding/json"
	"strings"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
	pb "SDCS/kvrpc"
)

func keyHashFunc(key string) int {
    hashFunc := fnv.New32a()
    hashFunc.Write([]byte(key))
    return int(hashFunc.Sum32() % 3)
}


var (
    // id
    id = flag.Int("id", 0, "The server id")
    my_id int
    // http port
    port = flag.Int("port", 8080, "The server port")
    http_port int
    // rpc port
	rpc_port string
    // data to store json
    data map[string]json.RawMessage
    // rpc client
    client [3]pb.ServiceKVClient
)

type server struct {
	pb.UnimplementedServiceKVServer
}

func (s *server) PostKV(ctx context.Context, in *pb.PostRequest) (*pb.PostReply, error) {
    key := in.GetKey()
    json := in.GetJson()
    data[key] = json
	log.Printf("(%v) Post Received Key: %v Json: %v", rpc_port, key, string(json))
	return &pb.PostReply{Success: true}, nil
}

func (s *server) GetKV(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
    key := in.GetKey()
	log.Printf("(%v) Get Received Key: %v", rpc_port, key)
    json, exists := data[key]
    if !exists {
        return &pb.GetReply{Success: false, Json: nil}, nil
    }
	return &pb.GetReply{Success: true, Json: json}, nil
}

func (s *server) DeleteKV(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
    key := in.GetKey()
	log.Printf("(%v) Delete Received Key: %v", rpc_port, key)
    _, exists := data[key]
    if !exists {
        return &pb.DeleteReply{Success: false}, nil
    }
    delete(data, key)
	return &pb.DeleteReply{Success: true}, nil
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    var temp map[string]interface{}
    var unique_key string

    // 读取请求体并将其存储为 json.RawMessage
	var rawMessage json.RawMessage
    err := json.NewDecoder(r.Body).Decode(&rawMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // 将 json.RawMessage 解码为一个 map
    err = json.Unmarshal(rawMessage, &temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取第一个键
	for key := range temp {
        unique_key = key
        break
	}

    target_id := keyHashFunc(unique_key)

    fmt.Printf("POST target_id is %d\n", target_id)

    if target_id == my_id {
        data[unique_key] = rawMessage
    } else {
        // Contact the server and print out its response.
	    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
        _, err := client[target_id].PostKV(ctx, &pb.PostRequest{Key: unique_key, Json: rawMessage})
        if err != nil {
            log.Fatalf("could not post kv: %v", err)
        }
    }

    fmt.Fprintf(w, "Received POST with key:%v; json: %+v\n", unique_key, string(rawMessage))
}

func getHandler(w http.ResponseWriter, key string) {
    target_id := keyHashFunc(key)
    var json json.RawMessage

    fmt.Printf("GET target_id is %d\n", target_id)

    if target_id == my_id {
        var exists bool
        json, exists = data[key]
        if !exists {
            // 如果没有匹配的key，返回404
            http.Error(w, "Not found key", http.StatusNotFound)
            return
        }
    } else {
        // Contact the server and print out its response.
	    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
        r, err := client[target_id].GetKV(ctx, &pb.GetRequest{Key: key})
        if err != nil || !r.Success {
            http.Error(w, "Not found key", http.StatusNotFound)
            return
        }
        json = r.GetJson()
    }

    fmt.Fprintf(w, "Received GET return %+v\n", string(json))
}

func deleteHandler(w http.ResponseWriter, key string) {
    target_id := keyHashFunc(key)

    fmt.Printf("DELETE target_id is %d\n", target_id)

    if target_id == my_id {
        _, exists := data[key]
        if !exists {
            fmt.Fprintf(w, "%d\n", 0)
            return
        }
        delete(data, key)
        fmt.Fprintf(w, "%d\n", 1)
    } else {
        // Contact the server and print out its response.
	    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
        r, err := client[target_id].DeleteKV(ctx, &pb.DeleteRequest{Key: key})
        if err != nil || !r.Success {
            fmt.Fprintf(w, "%d\n", 0)
            return
        }
        fmt.Fprintf(w, "%d\n", 1)
    }
    fmt.Fprintf(w, "Received DELETE request for key: %v\n", key)
}

func handler(w http.ResponseWriter, r *http.Request) {
    // 根据请求路径拆分URL
    path := r.URL.Path
    
    // 根路径（处理 POST 请求）
    if r.Method == "POST" && path == "/" {
        postHandler(w, r)
        return
    }

    parts := strings.Split(strings.Trim(path, "/"), "/")
    
    // 处理 GET /{key} 和 DELETE /{key}
    if len(parts) == 1 {
        key := parts[0]
        switch r.Method {
        case "GET":
            getHandler(w, key)
        case "DELETE":
            deleteHandler(w, key)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
        return
    }

    // 如果没有匹配的路由，返回404
    http.Error(w, "Not found", http.StatusNotFound)
}

func main() {
    flag.Parse()
    my_id = *id
    http_port = *port
    rpc_port = "5005" + fmt.Sprint(my_id)
    // Initialize data map.
    data = make(map[string]json.RawMessage)
    // 创建 FNV-1a 哈希

    go func(){
        // Start KV RPC server.
        lis, err := net.Listen("tcp", fmt.Sprintf(":%s", rpc_port))
        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }
        s := grpc.NewServer()
        pb.RegisterServiceKVServer(s, &server{})
        log.Printf("server listening at %v", lis.Addr())
        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
        
    }()

    // Sleep for 100ms
    time.Sleep(1000 * time.Millisecond)

    // Start KV RPC client.
    for i := 0; i <= 2; i++ {
        if i == my_id {
            client[i] = nil
            continue
        }

        // Set up a connection to the server.
        conn, err := grpc.Dial("localhost:5005" + fmt.Sprintf("%d", i), grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
            log.Fatalf("did not connect: %v", err)
        }
        defer conn.Close()
        client[i] = pb.NewServiceKVClient(conn)
    }

    // Start HTTP service.
    http.HandleFunc("/", handler)
    fmt.Println("Starting server at :" + fmt.Sprint(http_port))
    err := http.ListenAndServe(":" + fmt.Sprint(http_port), nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}