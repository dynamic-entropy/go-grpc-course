# # For Greet Service

# gen:
# 	protoc -I . --go_out=. ./greet/greetpb/greet.proto \
# 	&& protoc -I . --go-grpc_out=. ./greet/greetpb/greet.proto

# clean:
# 	rm -rf greet/greetpb/*.pb.go

# runserver:
# 	go run greet/server/server.go
	
# runclient:
# 	go run greet/client/client.go


# # For Calulate Service
# gen:
# 	protoc -I . --go_out=. ./calculate/calculatepb/calculate.proto \
# 	&& protoc -I . --go-grpc_out=. ./calculate/calculatepb/calculate.proto

# clean:
# 	rm -rf calculate/calculatepb/*.pb.go

# runserver:
# 	go run calculate/server/server.go
	
# runclient:
# 	go run calculate/client/client.go


#For Blog Service
gen:
	protoc -I . --go_out=. ./blog/blogpb/blog.proto \
	&& protoc -I . --go-grpc_out=. ./blog/blogpb/blog.proto

clean:
	rm -rf blog/blogpb/*.pb.go

runserver:
	go run blog/server/server.go
	
runclient:
	go run blog/client/client.go
