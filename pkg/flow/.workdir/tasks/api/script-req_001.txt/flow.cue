package veg

import "encoding/json"

pick: {
	args: cow: string
}

tasks: [string]: {
	Out: _
	...
}

tasks: {
  @flow()
	call: { 
    @task(api.Call)
    req: {
      host: "https://postman-echo.com"
      // method: "GET"
      path: "/get"
      query: {
        cow: "moo"
      }
    }
    resp: {
      body: _
      statusCode: 200
    }
  }
  mask: {
    @task(st.Mask)
    val: call.resp.body
    mask: { headers: "x-amzn-trace-id": string }
  }
	out: { text: json.Indent(json.Marshal(mask.out), "", "  ") +"\n" } @task(os.Stdout)
}

