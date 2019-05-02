# singleton-echoloop-app
singleton echoloop app in go

### Usage example
```shell
$ go run *.go -text "test1" &                           
2019/05/02 19:25:56 attempt to start server at :8080
test1
test1
test1
$ go run *.go -text "test2"
2019/05/02 19:26:13 attempt to start server at :8080
2019/05/02 19:26:13 listen tcp :8080: bind: address already in use
2019/05/02 19:26:13 can't listen on 8080 port, the first instance must be listening
2019/05/02 19:26:13 sending <test2> to the first instance (http://localhost:8080)
2019/05/02 19:26:13 response: ok!
2019/05/02 19:26:13 echo loop for <test2> finished!
test1
test2
test1
test2
test1
test2
```
### Description
This echoloop prints parameter `text` to standard output every second.    
As you can see, after starting program is trying to listen `:8080` for incoming http connections.   
If port is already in use then the program assumes that an instance is running and tries to send parameters to `http://localhost:8080`.
Then running instance alternates outputs.