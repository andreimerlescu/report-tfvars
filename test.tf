variable "color" {
    description = "Deployment group color or label or tag"
    type = string
    default = "blue"
}

variable "nicknames" {
   description = "map of nicknames to full names"
   type = map(string)
   default = {
      "alfred" = "Alfred Wilson",
      "bob" = "Bob Gates",
      "paul" = "Paul LePage"
   }
}

variable "teams" {
    description = "a list of teams"
    type = list(string)
}

variable "departments" {
    type = list(string)
}

variable "company_name" {
    type = string
}

variable "active" {
    type = bool
    default = true
}

variable "max_team_size" {
   type = map(number)
	default = {
      "small" = 1,
      "medium" = 2,
      "large" = 3
   }
}

variable "managers" {
   description = "List of Engineering Managers"
   type = list(string)
   default = [
      "Alfred Wilson"
   ]
}

variable "public_key" {
    description = "Value of the SSH public key"
    type = string
    default = "value"
}

variable "key_name" {
    description = "Name of the SSH key used in resource creation"
    type = string
    default = "docker_cluster_sshkey"
}

variable "root_ebs_size" {
    description = "GB number of EBS volume size"
    type = number
    default = 60
}

variable "tmp_ebs_size" {
    description = "GB number of EBS volume size"
    type = number
    default = 30
}

variable "override" {
    description = "Flag to enable quantity"
    type = bool
    default = false
}

variable "quantity" {
    description = "Override quantity number for instance count regardless of strategy used."
    type = number
    default = 1
}

variable "scalar" {
    description = "Multiple instance count by this number to scale deployment with a scalar."
    type = number
    default = 1
}

variable "strategy" {
    description = "Strategy to use for deployment to determine instance size and quantity."
    type = string
    default = "dev-large"
}

variable "ami" {
    description = "AMI ID (ami-) for base machine for resource creation"
    type = string
    default ="ami-0862be96e41dcbf74" # US-EAST-1 Ubuntu 22.04 LTS HVM_SSD AMI & IGNORED COMMENT
}

variable "size" {
    description = "Requires override and defines resource instance type for all quantities"
    type = string
    default = "c6a.medium"
}

variable "subnet" {
    description = "Subnet ID (subnet-) for base machine for resource creation"
    type = string
    default = "subnet-5446533c"       # Ignored comment
}

variable "dcmsgid" {
    description = "Docker Cluster Member AWS Security Group ID"
    type = string
    default = "sg-0a4934bfbfe6d95ca" # Ignored comment
}

variable "region" {
    description = "AWS Region (us-east-2 without a|b|c|d|e|f zone)"
    type = string
    default = "us-east-2"
}

variable "zone" {
    description = "AWS Availability Zone (us-east-2a)"
    type = string
    default = "us-east-2a"
}

variable "device_names" {
  description = "List of device names to be used for EBS volumes"
  type        = list(string)
  default     = ["/dev/xvdh", "/dev/xvdi", "/dev/xvdj", "/dev/xvdk", "/dev/xvdl", "/dev/xvdm", "/dev/xvdn", "/dev/xvdo", "/dev/xvdp"]
}
