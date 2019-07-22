resource "aws_ecr_repository" "snooper" {
  name = "snooper"

  provisioner "local-exec" {
    command = "packer build ../packer.json -var ECR_URL=${aws_ecr_repository.snooper.registry_id} && sleep 120"
  }
}

output "snooper-repo" {
  value = aws_ecr_repository.snooper.repository_url
}

