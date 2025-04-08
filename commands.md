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
pull  
```docker pull ghcr.io/yamada3-1990/twitter-like-memo:~~~```  
run  
```docker run -p 9000:9000 ghcr.io/yamada3-1990/twitter-like-memo:~~~```  

~~build~~  
~~```docker build -t twitter-like-memo/app:latest .```~~

~~run(backend)~~  
~~```docker run -d -p 9000:9000 twitter-like-memo/app:latest```~~  

~~run(front)~~  
~~```docker run -d -p 3000:3000 twitter-like-memo/app:latest```~~  

docker build(frontend)  
```docker build -t twitter-like-memo/app:latest .```  

docker run(frontend)  
```docker run -d -p 5173:5173 twitter-like-memo/app:latest```  
-> http://localhost:5173/ にアクセス


```docker-compose up```  
```docker-compose down```


フロントエンド開発環境
```cd frontend/twitter-like-memo```  
```npm run dev```
