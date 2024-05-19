kubectl apply -f psql.deployment.yaml -n persistence

docker image rm secap-auth:latest
docker build -t secap-auth:latest .

minikube image rm secap-auth:latest
minikube image load secap-auth:latest

kubectl apply -f .k8s/deployment.yaml -n secap-compass