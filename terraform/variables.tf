variable "gh_action_role" {
  type        = string
  description = "AWS IAM ARN for Terraform GitHub Actions"
}

variable "domain_name" {
  type        = string
  description = "The domain name to host site"
}

variable "comment" {
  description = "AWS Route53 Hosted Zone comment"
  type        = string
}

variable "aliases" {
  description = "Additional aliases for Cloudfront"
  type        = list(string)
  default     = []
}

variable "host_name" {
  description = "The host name for the site"
  type        = string
  default     = null
}
