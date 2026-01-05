# Backend for htg

## Starter
1. Copy content from `.env.sample`  then create `.env` then paste the content to this
2. startup docker

```bash
docker compose up -d
```

3. start server
```bash
go run .

# with port flag
# go run . --port 8080
```

or

```bash
go run main.go

# with port flag
# go run main.go --port 8080
```

4. Health check

```bash
curl http://localhost:<port>/health
```
