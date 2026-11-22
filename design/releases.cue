package design

import (
	"github.com/verdverm/veg/catalogs/gen/cli/schema"
)

CliReleases: schema.GoReleaser & {
	Disable:  false
	Draft:    false
	Author:   "verdverm"
	Homepage: "https://github.com/verdverm/veg"

	GitHub: {
		Owner: "verdverm"
		Repo:  "veg"
	}

	Docker: {
		Maintainer: "verdverm <verdverm@users.noreply.github.com>"
		Repo:       "veg"
	}
}
