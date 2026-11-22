package prompt

import (
  schema "github.com/verdverm/veg/schemas/prompt"
)

// Same prompt from creators as a workflow task
Prompt: schema.Prompt & {
  @task(prompt.Prompt)
}