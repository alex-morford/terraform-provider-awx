resource "awx_job_template_instance_group" "default" {
  instance_group_ids = [1]
  job_template_id    = 100
}