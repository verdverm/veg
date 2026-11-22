tasks: {
  @flow()
	r: { filename: "in.txt", contents: string } @task(os.ReadFile)
  o: { text: r.contents } @task(os.Stdout)
}

