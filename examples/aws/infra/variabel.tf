variable "vpc_name" {
  type = string
}
variable "cluster_name" {
  type = string
}

variable "region" {
  type = string
}
variable "cidr" {
  type = string
}
variable "private_subnets" {}
variable "public_subnets" {}
variable "eks_version" {
  type = string
}
variable "eks_instance_types" {
  type = list(string)
}
variable "cidr_blocks" {
  type = string
}


variable "azs" {
  type = list(any)
}

