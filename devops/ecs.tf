resource "aws_ecs_cluster" "fargate" {
  name = "${terraform.workspace}-fargate"
}

