
build:
	docker build -t apollonion-conversation-service .

run: 
	docker run --name apollonion-conversation-service -p 8080:8080 apollonion-conversation-service

connect:
	echo -n “Ground Control For Major Tom” | nc localhost 8080