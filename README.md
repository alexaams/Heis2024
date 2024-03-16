# Heis2024
A real-time project focused on the simultaneous control of multiple elevators.

This project is part of the course **TTK4145 Real Time Programming**, aimed at teaching us about creating robust real-time systems with a focus on code quality.

We chose to write this project in Go, leveraging its robust support for concurrent programming and efficient internal communication capabilities, which are ideal for real-time systems.

Our solution employs UDP (User Datagram Protocol) for lightweight communication between elevators, structured in a peer-to-peer (P2P) network where each elevator node exchanges information directly.

## Configuration
In the script found in `config_folder/config`, we have set the parameters that describe the system's size, such as the number of elevators and the amount of floors. The other directory in `config_folder/types` contains most of the defined packages and descriptors for this project.

### Shared Variables
We would also like to mention that all variables with the **G_** prefix function as the global shared variables in this project.

## How to Run the Code
From the working directory, run the commands:

    chmod +x elevatorstart.sh 
    ./elevatorstart.sh

Running this will create a terminal which has to be killed to stop the program.

If you wish to run the code directly, from the working directory, run the command: 

    go run main.go

# Authors
* Alexander Riis Amsjø
* Håkon Waage
* Jørgen Hazeland Baugerud
