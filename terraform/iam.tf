provider "aws" {
    alias = "iam-aws"
    region = "ap-southeast-1"
    profile = "terraform"

    assume_role {
        role_arn     = "arn:aws:iam::407464631290:role/terraform-iam-role"
        session_name = "terraform-iam-role"
    }
}

resource "aws_iam_user" "dynamodb-assumer" {
    provider = aws.iam-aws
    name = "dynamodb-role-assumer"
}

resource "aws_iam_role" "dynamodb-role" {
    provider = aws.iam-aws
    name = "dynamodb-role"
    assume_role_policy = data.aws_iam_policy_document.dynamodb-assumer-assume-role.json
}

resource "aws_iam_role_policy_attachment" "dynamodb-role" {
    provider = aws.iam-aws
    role = aws_iam_role.dynamodb-role.name
    policy_arn = aws_iam_policy.dynamodb-policy.arn
}

resource "aws_iam_access_key" "dynamodb-assumer" {
    provider = aws.iam-aws
    user = aws_iam_user.dynamodb-assumer.name
}

data "aws_iam_policy_document" "dynamodb-assumer-assume-role" {
  statement {
    effect    = "Allow"
    actions   = ["sts:AssumeRole"]

    principals {
        type = "AWS"
        identifiers = [aws_iam_user.dynamodb-assumer.arn]
    }
  }
}

data "aws_iam_policy_document" "dynamodb-policy" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:BatchGetItem", "dynamodb:BatchWriteItem", "dynamodb:ConditonCheckItem", "dynamodb:GetItem", "dynamodb:ListTables", "dynamodb:PutItem", "dynamodb:Scan", "dynamodb:UpdateItem"]
    resources = [aws_dynamodb_table.case_tracker.arn]
  }
}

resource "aws_iam_policy" "dynamodb-policy" {
  provider = aws.iam-aws
  name   = "dynamo-db-app-policy"
  policy = data.aws_iam_policy_document.dynamodb-policy.json
}


