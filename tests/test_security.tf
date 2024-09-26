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

resource "aws_s3_bucket" "public_read" {
  bucket = "public-read-bucket"
  acl    = "public-read"
}

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
