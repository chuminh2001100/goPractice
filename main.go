package main

import (
	"context"
	"encoding/json"
	"mime/multipart"

	// "errors"
	pb "HelloService" // Import generated code
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	// "github.com/go-co-op/gocron"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-zoo/bone"
	"github.com/pkg/errors"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var task = func() {
	fmt.Println("get log hehee")
}

// Request struct là dữ liệu đầu vào cho endpoint
type Request23 struct {
	Name string `json:"name"`
}


type UploadReq struct {
	file multipart.File
	check bool
}


type RequestReply struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Id      string `json:"id"`
}

type ResponseReply struct {
	Id string `json:"id"`
}

type Request_IDstudent struct {
	Name  string `json:"name"`
	Age   int32  `json:"age"`
	Class string `json:"class"`
}

type Response_IDstudent struct {
	Success string `json:"success"`
}

// Response struct là dữ liệu trả về từ endpoint
type Response struct {
	Greeting string `json:"greeting"`
}

// HelloWorldService là dịch vụ cung cấp phương thức SayHello
type HelloWorldService interface {
	SayHello(ctx context.Context, name string) (string, error)
	SayReply(ctx context.Context, name string) (string, error)
	UpSertAtt(ctx context.Context, name string) (string, error)
}

type loggingMiddleware struct {
	logger log.Logger
	next   HelloWorldService
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           HelloWorldService
}

type HelloWorldServiceImpl struct {
	clientGrpc pb.HelloServiceClient
	name       string
}

// SayHello trả về chuỗi chào "Hello, World!"
func (s HelloWorldServiceImpl) SayHello(ctx context.Context, name string) (string, error) {
	return "Hello, " + name + "!" + "Thanh" + "Hello" + s.name, nil
}

func (s HelloWorldServiceImpl) SayReply(ctx context.Context, name string) (string, error) {
	return "Hello, " + name + "!" + "Hung" + "Hello" + s.name, nil
}

func (s HelloWorldServiceImpl) UpSertAtt(ctx context.Context, name string) (string, error) {
	name123 := "API kafka 1235466665555555"
	response, err := s.clientGrpc.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		fmt.Println("un wirelab")
		return "fail", err
	}
	fmt.Printf("Goi rgpc xong roi hehe: %s", response.GetMessage())
	fmt.Println("There Here")
	hanldeSerivce := "Hello, " + name123
	return hanldeSerivce, nil
}

func (mw loggingMiddleware) SayHello(ctx context.Context, name string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"input", name,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.SayHello(ctx, name)
	return
}

func (mw loggingMiddleware) SayReply(ctx context.Context, name string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "SayReply",
			"input", name,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.SayReply(ctx, name)
	return
}

func (mw loggingMiddleware) UpSertAtt(ctx context.Context, name string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Upsert Attributes",
			"input", name,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	fmt.Println("Second Here")
	output, err = mw.next.UpSertAtt(ctx, name)
	return
}

// SayHelloEndpoint định nghĩa endpoint của phương thức SayHello

func SayHelloEndpoint(svc HelloWorldService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(Request23) // Kiểm tra kiểu của request
		if !ok {
			return nil, ErrInvalidRequest
		}

		greeting, err := svc.SayHello(ctx, req.Name)
		if err != nil {
			return nil, err
		}

		return Response{Greeting: greeting}, nil
	}
}

func SayReplyEndpoint(svc HelloWorldService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(RequestReply) // Kiểm tra kiểu của request
		if !ok {
			return nil, ErrInvalidRequest
		}
        usernameRegex := regexp2.MustCompile(`^(?=.*[A-Za-z])[A-Za-z0-9_]{6,32}$`, 0x0200)
		matched, _ := usernameRegex.MatchString(req.Name)
		if !matched {
			fmt.Println("name not ok");
			return nil, ErrInvalidRequest
		}
		// greeting, err := svc.SayHello(ctx, req.Name)
		// if err != nil {
		// 	return nil, err
		// }

		return ResponseReply{Id: req.Id}, nil
	}
}

func UpsertEndpoint(svc HelloWorldService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UploadReq) // Kiểm tra kiểu của request
		if !ok {
			return nil, ErrInvalidRequest
		}
		fmt.Printf("check req: %+v\n", req)
		greeting, err := svc.UpSertAtt(ctx, "hello em")
		if err != nil {
			return nil, err
		}
		fmt.Println("Cam on em")
		return Response_IDstudent{Success: greeting}, nil
	}
}

func (mw instrumentingMiddleware) SayHello(ctx context.Context, name string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SayHello", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.SayHello(ctx, name)
	return
}

func (mw instrumentingMiddleware) SayReply(ctx context.Context, name string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SayReply", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.SayReply(ctx, name)
	return
}

func (mw instrumentingMiddleware) UpSertAtt(ctx context.Context, name string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Upsert Att", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	fmt.Println("first Here")
	output, err = mw.next.UpSertAtt(ctx, name)
	fmt.Println("Four Here")
	return
}

func DecodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// Lấy giá trị của tham số "name" từ query string
	name := r.URL.Query().Get("name")
	test := bone.GetQuery(r, "name")
	fmt.Printf("Test bone Query: %s\n", test)
	fmt.Println(test)
	// Tạo một đối tượng Request từ giá trị "name" trích xuất
	request := Request23{Name: name}

	return request, nil
}

func DecodeRequestReply(_ context.Context, r *http.Request) (interface{}, error) {
	// Lấy giá trị của tham số "name" từ query string
	if !strings.Contains(r.Header.Get("Content-Type"),"application/json"){
		return nil, errors.New("unsupported content type")
	}
	// Tạo một đối tượng Request từ giá trị "name" trích xuất
	request := RequestReply{
		Id: bone.GetValue(r,"id"),
	}
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("Un confirmed")
		return nil, errors.Wrap(err, "Fail adu vjp")
	}
	return request, nil
}

func DecodeRequestUpsert(_ context.Context, r *http.Request) (interface{}, error) {
	// Lấy giá trị của tham số "name" từ query string
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		return nil, errors.New("unsupported content type")
	}
	// req := Request_IDstudent{}
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	fmt.Println("Un confirmed")
	// 	return nil, errors.Wrap(err, "Fail adu vjp")
	// }

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("Error file: %s\n", err)
		return nil, errors.Wrap(err, "Fail adu vjp")
	}
    fmt.Printf("In full handler: %+v\n", handler)
    fmt.Printf("In full file: %+v\n", file)
	req := UploadReq {
		file: file,
		check: true,
	}

	return req, nil
}

// EncodeResponse mã hóa kết quả từ endpoint thành response HTTP
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Println("Amee")
	return json.NewEncoder(w).Encode(response)
}

func main() {
	// Khởi tạo logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	// Khởi tạo router bằng bone
	r := bone.New()
	// Khởi tạo dịch vụ HelloWorldService
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("connect fail")
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)
	name := "API kafka"
	response, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		fmt.Println("un wirelab")
	}
	fmt.Println(response)
	var helloService HelloWorldService
	helloService = HelloWorldServiceImpl{
		clientGrpc: c,
		name:       "ChuVanMinh",
	}
	// Tạo endpoint
	helloService = loggingMiddleware{logger, helloService}
	helloService = instrumentingMiddleware{requestCount, requestLatency, countResult, helloService}
	sayHelloEndpoint23 := SayHelloEndpoint(helloService)
	r.Get("/hello", httptransport.NewServer(
		sayHelloEndpoint23,
		DecodeRequest,
		EncodeResponse,
	))
	r.Put("/hello/:id", httptransport.NewServer(
		SayReplyEndpoint(helloService),
		DecodeRequestReply,
		EncodeResponse,
	))

	r.Post("/api/upsert", httptransport.NewServer(
		UpsertEndpoint(helloService),
		DecodeRequestUpsert,
		EncodeResponse,
	))
	metricsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	// s := gocron.NewScheduler(time.UTC)
	// _, _ = s.Every(1).Second().Do(task)
	// _, _ = s.Every(1).Minute().Do(task)
	// _, _ = s.Every(1).Month(1).Do(task)
	// fmt.Println(len(s.Jobs())) // Print the number of jobs before clearing
	// s.Clear()                  // Clear all the jobs
	// fmt.Println(len(s.Jobs())) // Print the number of jobs after clearing
	// s.StartAsync()
	// Gán handler của metrics vào router
	r.Handle("/metrics", metricsHandler)
	done := make(chan bool)

	// Khởi động server trên cổng 8080
	go func() {
		// Khởi động server trên cổng 8080
		if err := http.ListenAndServe(":8086", r); err != nil {
			fmt.Println("Error starting server:", err)
		}
		fmt.Println("DM AO THAT DAY, SOC SON")
		done <- true
	}()

	fmt.Println("Server is running...")
	<-done
	fmt.Println("Server stopped")
}

var ErrInvalidRequest = errors.New("invalid request")
