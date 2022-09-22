docker-compose up -d

go run cmd/api/main.go

Rabbitmq producer will allways run. 


example curls:
curl --request POST \
  --url http://localhost:8080/message \
  --header 'Content-Type: application/json' \
  --data '{
    "sender": "String", 
    "receiver": "String", 
    "message": "String"
}'

curl --request GET \
  --url 'http://localhost:8080/message/list?sender=String&receiver=String'
