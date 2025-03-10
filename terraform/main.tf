module "website" {
  source  = "cloudposse/s3-bucket/aws"
  version = "4.10.0"

  bucket_name             = data.aws_route53_zone.this.name
  s3_object_ownership     = "BucketOwnerEnforced"
  block_public_policy     = false
  restrict_public_buckets = false
  versioning_enabled      = false
  context                 = module.this.context
  website_configuration = [
    {
      index_document = "index.html"
      error_document = "/404.html"
      routing_rules  = null
    }
  ]
}

module "cdn" {
  source  = "cloudposse/cloudfront-s3-cdn/aws"
  version = "0.97.0"

  name                        = "cdn"
  comment                     = var.comment
  origin_bucket               = module.website.bucket_id
  aliases                     = [local.hostname, local.mta_sts]
  external_aliases            = var.aliases
  dns_alias_enabled           = true
  dns_allow_overwrite         = true
  website_enabled             = true
  s3_website_password_enabled = true
  allow_ssl_requests_only     = false
  price_class                 = "PriceClass_All"
  default_ttl                 = 86400
  min_ttl                     = 3600
  max_ttl                     = 2592000
  minimum_protocol_version    = "TLSv1.2_2021"
  parent_zone_id              = data.aws_route53_zone.this.id
  acm_certificate_arn         = data.aws_acm_certificate.this.arn
  context                     = module.this.context
  custom_error_response = [
    {
      error_caching_min_ttl = null
      error_code            = 404
      response_code         = 200
      response_page_path    = "/404.html"
    }
  ]
}

locals {
  hostname = join(".", compact([var.host_name, data.aws_route53_zone.this.name]))
  mta_sts  = join(".", ["mta-sts", data.aws_route53_zone.this.name])
}
