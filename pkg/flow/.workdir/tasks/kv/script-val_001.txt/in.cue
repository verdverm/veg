@flow()

set: {
	@task(kv.Mem)
	key: "foo"
	val: {
		bar: "baz"
		num: int
	}
}

get: {

	wait: {
		@task(os.Sleep)
		duration: "1s"
	}

	load: {
		@task(kv.Mem)
		$dep: wait.done
		key:  "foo"
	}

	print: {text: load.val.bar + "\n"} @task(os.Stdout)
}
