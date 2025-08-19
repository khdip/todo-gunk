This is a simple web app where you can create a task, update it, delete it and mark it as completed.
In this project I implemented GRPC using a package called gunk. Using gunk you don't need .proto file to generate the pb.go files. You can just write go code inside the .gunk file.
Please change the environment variables of the server and client accordingly in the config file.
Find the config file of the server navigate to **todo/env**.
Find the config file of the client navigate to **cms/env**.
To run the server: *go run todo/main.go*
To run the client: *go run cms/main.go*
To view the web app in the browser go to: localhost:8080 (or the port number you specified in the config file).
