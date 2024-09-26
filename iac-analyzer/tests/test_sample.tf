# tests/test_sample.tf

# Security Group with Open SSH Access
resource "aws_security_group" "open_ssh" {
  name        = "open_ssh"
  description = "Open SSH access"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Public S3 Bucket
resource "aws_s3_bucket" "public_read" {
  bucket = "public-read-bucket"
  acl    = "public-read"
}

# Unencrypted RDS Instance
resource "aws_db_instance" "unencrypted_db" {
  allocated_storage    = 20
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "foo"
  password             = "bar"
  parameter_group_name = "default.mysql5.7"
}

# Oversized EC2 Instance
resource "aws_instance" "oversized_instance" {
  ami           = "ami-12345678"
  instance_type = "t3.large"
}

# Large EBS Volume
resource "aws_ebs_volume" "large_volume" {
  availability_zone = "us-west-2a"
  size              = 2000
}

# Unattached Elastic IP
resource "aws_eip" "unattached_eip" {
  vpc = true
}