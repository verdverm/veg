import "encoding/json"

apicall: {
  in: string
	r: { filename: in, contents: string } @task(os.ReadFile)
  j: json.Unmarshal(r.contents)
	r1: { req: j, resp: body: string } @task(api.Call)
  resp: r1
}


main: {
  @flow(main)

  input: { text: "apicalling" } @task(os.Stdout)

  call: apicall & { in: "req.json" }
  // final: { text: call.resp.body } @task(os.Stdout)
  final: { text: call.r1.resp.body } @task(os.Stdout)

  call2: apicall & { in: "req2.json" }
  final2: { text: call2.r1.resp.body } @task(os.Stdout)
}

