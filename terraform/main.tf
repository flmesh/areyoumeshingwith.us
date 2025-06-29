resource "aws_amplify_app" "website" {
  name       = var.domain_name
  repository = var.repository

  platform          = "WEB"

  environment_variables = {
    BASEURL = var.domain_name
  }
  
  custom_rule {
    source = "/.well-known/<*>"
    status = "200"
    target = "/well-known/<*>.txt"
  }

  custom_rule {
    source = "/<*>"
    status = "404"
    target = "/404.html"
  }
}

resource "aws_amplify_branch" "main" {
  app_id                      = aws_amplify_app.website.id
  branch_name                 = "main"
  stage                       = "PRODUCTION"
  enable_pull_request_preview = true
  framework                   = "Web"
}

resource "aws_amplify_domain_association" "website" {
  app_id                 = aws_amplify_app.website.id
  domain_name            = var.domain_name
  enable_auto_sub_domain = true

  certificate_settings {
    custom_certificate_arn = data.aws_acm_certificate.this.arn
    type                   = "CUSTOM"
  }

  sub_domain {
    branch_name = aws_amplify_branch.main.branch_name
    prefix      = ""
  }

  sub_domain {
    branch_name = aws_amplify_branch.main.branch_name
    prefix      = "mta-sts"
  }
}
