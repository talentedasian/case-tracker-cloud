provider "aws" {
  region  = "ap-southeast-1"
  profile = "terraform"

  assume_role {
    role_arn     = "arn:aws:iam::407464631290:role/dynamodb-management-role"
    session_name = "terraform-dynamodb-role"
  }

}

resource "aws_dynamodb_table" "case_tracker" {
  name           = "case_tracker"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "inmate_id"
  range_key      = "sort_key"
  deletion_protection_enabled = true

  attribute {
    name = "inmate_id"
    type = "S"
  }

  attribute {
    name = "sort_key"
    type = "S"
  }

  tags = {
    Name = "case-tracker"
  }
}

import {
  to = aws_dynamodb_table.case_tracker
  id = "case_tracker"
}

