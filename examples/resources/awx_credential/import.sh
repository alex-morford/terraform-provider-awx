# The inputs field contains values that are not returned by the AWX API and thus not available in state.
# The first plan/apply after import will result in a modification to the inputs & kind so that the state can be updated.

terraform import awx_credential.example 5
