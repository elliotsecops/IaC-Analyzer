resource "aws_instance" "oversized_instance" {
  ami           = "ami-12345678"
  instance_type = "t3.large"
}

resource "aws_ebs_volume" "large_volume" {
  availability_zone = "us-west-2a"
  size              = 2000
}

resource "aws_eip" "unattached_eip" {
  vpc = true
}
