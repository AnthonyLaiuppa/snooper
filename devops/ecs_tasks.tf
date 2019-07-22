data "template_file" "snooper" {
  template = file("templates/ecs/snooper.json.tpl")
  vars = {
    REPOSITORY_URL = aws_ecr_repository.snooper.repository_url
    AWS_REGION     = var.AWS_REGION
    SWHURL = var.SWHURL
  }
}

resource "aws_ecs_task_definition" "snooper" {
  family                   = "snooper"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  container_definitions    = data.template_file.snooper.rendered
  execution_role_arn       = aws_iam_role.ecs_task_assume.arn
}

resource "aws_ecs_service" "snooper" {
  name            = "snooper-${terraform.workspace}"
  cluster         = aws_ecs_cluster.fargate.id
  launch_type     = "FARGATE"
  task_definition = aws_ecs_task_definition.snooper.arn
  desired_count   = 1

  network_configuration {
    # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
    # force an interpolation expression to be interpreted as a list by wrapping it
    # in an extra set of list brackets. That form was supported for compatibilty in
    # v0.11, but is no longer supported in Terraform v0.12.
    #
    # If the expression in the following list itself returns a list, remove the
    # brackets to avoid interpretation as a list of lists. If the expression
    # returns a single list item then leave it as-is and remove this TODO comment.
    subnets = [module.base_vpc.private_subnets[0]]
    # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
    # force an interpolation expression to be interpreted as a list by wrapping it
    # in an extra set of list brackets. That form was supported for compatibilty in
    # v0.11, but is no longer supported in Terraform v0.12.
    #
    # If the expression in the following list itself returns a list, remove the
    # brackets to avoid interpretation as a list of lists. If the expression
    # returns a single list item then leave it as-is and remove this TODO comment.
    security_groups = [aws_security_group.ecs.id]
  }
}

