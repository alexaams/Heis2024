rungo(){
    go run main.go
}
trap 'rungo' SIGINT SIGTERM

rungo