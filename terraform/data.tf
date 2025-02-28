data "aws_route53_zone" "this" {
  name = var.domain_name
}

data "aws_acm_certificate" "this" {
  domain      = var.domain_name
  types       = ["AMAZON_ISSUED"]
  statuses    = ["ISSUED"]
  most_recent = true
}

data "aws_region" "current" {}
