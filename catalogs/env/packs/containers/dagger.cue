package containers

import (
	"github.com/verdverm/veg/catalogs/env/utils"
)

dagger: {
	#ver: string | *"1.0.0-beta.28"

	// hack for now to maintain consistency in pack.<tool>.cli.install
	cli: install: [
		utils.githubBin & {#repo: "dagger/dagger", #ver: dagger.#ver},
	]
}
