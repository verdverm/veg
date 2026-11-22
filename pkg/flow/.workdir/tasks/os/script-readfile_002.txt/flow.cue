tasks: {
  @flow()
	r: { filename: "in.json", contents: string } @task(os.ReadFile)
  o: { text: r.contents } @task(os.Stdout)
}

