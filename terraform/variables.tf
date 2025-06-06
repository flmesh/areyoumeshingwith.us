variable "gh_action_role" {
  description = "AWS IAM ARN for Terraform GitHub Actions"
  type        = string
}

variable "repository" {
  description = "Source respository"
  type        = string
}

variable "domain_name" {
  description = "The domain name to host site"
  type        = string
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
