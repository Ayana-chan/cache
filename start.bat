start go run run_nodeserver.go -5001 -10
start go run run_nodeserver.go -5002 -10
start go run run_nodeserver.go -5003 -10

start go run run_router.go localhost:5001 localhost:5002 localhost:5003 -6000 -5