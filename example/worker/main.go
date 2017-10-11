package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tkanos/konsumerou"
	"github.com/tkanos/konsumerou/example/worker/config"
	"github.com/tkanos/konsumerou/example/worker/middleware"
	"github.com/tkanos/konsumerou/example/worker/myservice"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func init() {
	// read config file and log an error if it's not present
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	// Tracing domain.
	var tracer opentracing.Tracer
	{
		zipkinAddr := config.Config.ZipkinURI
		if zipkinAddr != "" {
			//create a collector
			collector, err := zipkin.NewHTTPCollector(zipkinAddr)
			if err != nil {
				logger.Fatal("Unable to create a zipkin collector : ", err)
			}
			defer collector.Close()

			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, true, "kafka-consumer:"+strconv.Itoa(config.Config.Port), "kafka-consumer"),
				zipkin.TraceID128Bit(true),
			)
			if err != nil {
				logger.Fatal("Unbale to create a Zipkin Tracer : ", err)
			}

			// explicitly set our tracer to be the default tracer.
			opentracing.SetGlobalTracer(tracer)

		} else {
			tracer = opentracing.GlobalTracer() // no-op
		}
	}

	// create web server just for the /healthz
	httpAddr := ":" + strconv.Itoa(config.Config.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(httpAddr, mux)
	}()

	// check kafka config
	if config.Config.KafkaBrokers == "" {
		logger.Fatalln("kafka broker is empty")
	}

	//setup login failed event consumer
	done := make(chan bool)
	myService := myserviceEventListener(done)
	defer myService.Close()

	//quit application
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	done <- true

	fmt.Fprintf(os.Stdout, "the user choose to interrupt the program")
}

func myserviceEventListener(done chan bool) konsumerou.Listener {
	//subscribe to topic

	//create our service with tracing and logging
	service := myservice.NewService()
	service = myservice.NewServiceTracing(service)
	service = myservice.NewServiceLogging(service)

	//ProcessMessage endpoint
	handler := myservice.MakeMyServiceEndpoint(service)

	//add metrics middleware to service
	handler = middleware.NewMetricsService("myService_ProcessMessage", handler)
	//add log middleware to service
	handler = middleware.NewLogService(logger, handler)

	// create config to handle offset
	clusterConfig := cluster.NewConfig()
	clusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	listener, err := konsumerou.NewListener(
		strings.Split(config.Config.KafkaBrokers, ";"),
		"my-group",
		config.Config.MyServiceKafkaTopic,
		handler, clusterConfig,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start consumer: %s", err)
		os.Exit(-3)
	}

	//service subscription
	err = listener.Subscribe(done)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start user login failed event consumer: %s", err)
		os.Exit(-3)
	}

	return listener
}
