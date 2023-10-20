// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

output "id" {
  description = "The ARN of the Domain."
  value       = aws_codeartifact_domain.this.id
}

output "arn" {
  description = "The ARN of the Domain."
  value       = aws_codeartifact_domain.this.arn
}

output "owner" {
  description = "The AWS account ID that owns the domain."
  value       = aws_codeartifact_domain.this.owner
}

output "repository_count" {
  description = "The number of repositories in the domain."
  value       = aws_codeartifact_domain.this.repository_count
}

output "created_time" {
  description = "A timestamp that represents the date and time the domain was created in RFC3339 format."
  value       = aws_codeartifact_domain.this.created_time
}

output "asset_size_bytes" {
  description = "The total size of all assets in the domain."
  value       = aws_codeartifact_domain.this.asset_size_bytes
}

output "tags_all" {
  description = " A map of tags assigned to the resource, including those inherited from the provider default_tags configuration block."
  value       = aws_codeartifact_domain.this.tags_all
}
