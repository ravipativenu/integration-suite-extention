docker build -t ravipativenu/integration-suite-extention:latest -f Dockerfile .
docker push ravipativenu/integration-suite-extention:latest
kubectl replace --force -f deployment.yaml -n default
