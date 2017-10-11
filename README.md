# konsumerou 

konsumerou is a kafka consumer router. 
When you have your http endpoint on your main you have to define 2 things:
- the route ("/mypoint" for example)
- And an handler (that will execute the code of this endpoint)

with konsumerou it's almost the same, you give to konsumerou:
- the topic you want to listen
- your handler that will execute the code of this endpoint

```golang
// Create the service endpoint
handler := myservice.MakeMyServiceEndpoint(service)

listener, err := konsumerou.NewListener(ctx, 
		"localhost:9092", //kafka brokers
		"my-group",       // group id
		"my-topic",       // the topic name
		handler, //the handler
		nil   // cluster config
	)


// Subscribe your service to the topic
listener.Subscribe(done)
defer listener.Close()
```

konsumerou share the same philosophy than go-kit.
Konsumerou wants that your service does not have to handle the transport.
So like this you can create middlewares for your service to handle tracing, logging, metrics ....

So you have to create an endpoint :

```golang
type MyServiceMessageProcessor interface {
	ProcessMessage(context.Context, *MyServiceMessage) error
}

func MakeMyServiceEndpoint(s MyServiceMessageProcessor) konsumerou.Handler {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) error {
		message := MyServiceMessage{}
    if err := json.Unmarshal(msg, &message); err != nil {
		  return err
	  }

		return s.ProcessMessage(ctx, message)
	}
}
```

You can find an example on the example folder.
Enjoy
