curl -X GET "http://localhost:8080/gettask?id=37371b8b-e91e-4d20-bb4f-0b4794c6274c" -H "Accept: application/json"
curl -X POST http://localhost:8080/newtask -H "Content-Type: application/json"
curl -X DELETE "http://localhost:8080/deletetask?id=37371b8b-e91e-4d20-bb4f-0b4794c6274c" -H "Accept: application/json"