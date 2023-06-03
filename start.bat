start go run run_nodeserver.go -5001 -5
start go run run_nodeserver.go -5002 -5
start go run run_nodeserver.go -5003 -5

start go run run_router.go localhost:5001 localhost:5002 localhost:5003 -5999 -5