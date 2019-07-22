module "base_vpc" {
  source = "github.com/terraform-aws-modules/terraform-aws-vpc"

  name = "${terraform.workspace} - base_vpc"
  cidr = "10.0.0.0/16"

  azs             = flatten([split(",", local.azs)])
  private_subnets = flatten([split(",", var.CIDR_PRIVATE)])
  public_subnets  = flatten([split(",", var.CIDR_PUBLIC)])

  enable_nat_gateway = true
  single_nat_gateway = true

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

