rungo(){
    attemps=0
    max_attemts=5

    while [ $attempts -lt $max_attempts ]; do
        echo "Trying to connect elevator..."
        sleep 2
        go run main.go
        ((attemps++))
    done

    echo "Failed to reconnect elevator"
    exit 1
}

killelevator(){
    echo "terminating elevator"
    pkill -f 'startelevator.sh'
}

trap 'rungo' SIGINT
trap 'killelevator' SIGTERM SIGUSR1

rungo