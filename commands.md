```go run backend/cmd/main.go```

## ■GET
```curl.exe -X GET 'http://127.0.0.1:9000/memos'```  

### /search/keyword
```curl.exe -X GET 'http://127.0.0.1:9000/search/keyword?keyword=おはよう'```  

### /search/tags
```curl.exe -X GET 'http://localhost:9000/search/tags?tags=tag1&tags=tag2'```  
```curl.exe -X GET 'http://localhost:9000/search/tags?tags='```

## ■POST
```curl.exe -X POST --url 'http://localhost:9000/memos' -d 'title=test&body=test&tags=tag1,tag2'```   
```curl.exe -X POST --url 'http://localhost:9000/memos' -d 'title=title&body=title&tags=tag2,tag3'```   
```curl.exe -X POST --url 'http://localhost:9000/memos' -d 'title=test title&body=test body&tags=tag1,tag2,tag3'```   
```curl.exe -X POST --url 'http://localhost:9000/memos' -d 'title=おはよう&body=おはようございます&tags=greeting'```   


## ■DELETE
```curl.exe -X DELETE --url 'http://localhost:9000/memos?title=test&body=test'```  

↓スペースある場合は %20 にエンコード  
```curl.exe -X DELETE --url 'http://localhost:9000/memos?title=明日は&body=Hello%20こんにちは'  ```      

## ■Docker  
build  
```docker build -t twitter-like-memo/app:latest .```

run  
```docker run -d -p 9000:9000 twitter-like-memo/app:latest```
