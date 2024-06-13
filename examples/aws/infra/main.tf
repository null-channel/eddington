
provider "aws" {
  region  = var.region
  profile = "<profile>"
}

terraform {
  backend "s3" {
    bucket  = "<bucket_name>"
    key     = "<key>"
    region  = "<profile>"
    profile = "<profile>"
  }
}
module "vpc" {
  source               = "terraform-aws-modules/vpc/aws"
  version              = "3.18.1"
  name                 = var.vpc_name
  cidr                 = var.cidr
  azs                  = var.azs
  private_subnets      = var.private_subnets
  public_subnets       = var.public_subnets
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"           = "1"
  }
  tags = {
    Name = var.cluster_name
  }
}
resource "aws_db_subnet_group" "db_subnet_group" {
  name       = "subnets-lsaf-${var.application_name}-${var.env}-private"
  subnet_ids = module.vpc.private_subnets
  tags = {
    Name = var.cluster_name
  }
}


module "eks" {
  source                          = "terraform-aws-modules/eks/aws"
  version                         = "19.14.0"
  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = false
  cluster_name                    = var.cluster_name
  cluster_version                 = var.eks_version
  subnet_ids                      = module.vpc.private_subnets

  # EKS Addons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    aws-ebs-csi-driver = {
      addon_version            = "v1.18.0-eksbuild.1"
      service_account_role_arn = module.ebs_role.iam_role_arn
    }

  }

  vpc_id = module.vpc.vpc_id
  node_security_group_additional_rules = {

    vpc_inbound = {
      type        = "ingress"
      description = "communication between cluster nodes"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = module.vpc.private_subnets_cidr_blocks
    }
  }

  # Extend cluster security group rules
  cluster_security_group_additional_rules = {

    ingress_nodes_ephemeral_ports_tcp = {
      description                = "Nodes on ephemeral ports"
      protocol                   = "tcp"
      from_port                  = 1025
      to_port                    = 65535
      type                       = "ingress"
      source_node_security_group = true
    }

    vpc_all_outbound = {
      type        = "egress"
      description = "outbound communication from nodes including to DB (RDS)"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }

  }

  eks_managed_node_groups = {
    null_cloud_controle_plan_node_group = {
      desired_capacity = 2
      max_capacity     = 3
      min_capacity     = 1
      instance_types   = var.eks_instance_types
    }
  }
  tags = {
    Name = var.cluster_name
  }

}

module "ebs_role" {

  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name             = "${var.cluster_name}_eks_ebs"
  attach_ebs_csi_policy = true

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["kube-system:ebs-csi-controller-sa"]
    }
  }
  tags = {
    Name = var.cluster_name
  }

}
