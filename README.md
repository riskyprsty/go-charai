## go-charai

Golang Character AI implementation with simple rest API
### About

This project establishes a Gin Gonic HTTP framework that functions as an endpoint for interacting with the Character-AI service. So we can send a message and receive a reply from the AI character.

This project use [github.com/harmony-ai-solutions/CharacterAI-Golang](https://github.com/harmony-ai-solutions/CharacterAI-Golang) to communicate with Character AI

### Usage
**Please modify the code yourself, this code just intended for testing purposes**

Edit the `token` and `character` variable with your Character AI Token.

### Endpoints
#### Checking Character ID
`GET http://localhost:8080/ai?charid="<character_id>"`
```json
[
    {
        "status": "success",
        "body": "Ushio Noa"
    }
] 

```
#### Chat with AI
`GET http://localhost:8080/message?charid="<character_id>"&body="<message>"`
```json
[
    {
        "status": "success",
        "body": "Hi....."
    }
]                                                       
```

### Disclaimer

This project uses the following dependencies listed below:
- https://github.com/harmony-ai-solutions/CharacterAI-Golang
- https://github.com/gin-gonic/gin
