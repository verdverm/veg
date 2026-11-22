package agentic

import (
	"github.com/verdverm/veg/schemas"
	"github.com/verdverm/veg/schemas/common"
)

// Chat represents a call to an LLM via the `veg chat with' command.
// You can put these in your module to provide ChatGPT like interactions
// for the other components in your module, or make a module just for Chats.
Agent: {
	schemas.#veg // needed for reFerences

  #veg: agentic: { root: true, kind: "agent" }

	name: common.NameLabel

	// default model
	model?: string

	// description for man and machine
	description?: string

	// content or file path
	instruction: string

	tools: [...string]
	toolsets: [...string]
	mcp: [...string]
	subagents: [...string]

	// name of an agentic environment
	environ?: string
	
}
