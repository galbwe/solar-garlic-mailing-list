AWS_REGION=$1
ECR_REGISTRY=$2

aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REGISTRY
