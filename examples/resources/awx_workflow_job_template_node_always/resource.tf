resource "awx_workflow_job_template_node_always" "example_node_always" {
  id               = 201
  success_node_ids = [241, 914]
}
