## D2 API

[D2](https://github.com/terrastruct/d2) is modern diagram scripting language. I was fascinated by it and noticed that one of the [issue](https://github.com/terrastruct/d2/issues/207) was about building a web-playground where people can try things out without much fuss of setting it up locally.

### API Deploy details

URL to target: https://d2api.fly.dev/

Endpoint: `/getSvg`
Request type: `POST`
Payload type: `json`
Payload details:

```json
{
    diagramCode: "# D2 script here.
    x -> y"
}
```

Response type: `text/html`
Response: The generated svg for the provided D2 script, if the D2 script is valid.

### D2 Playground

Using this API, I have built a basic D2 playground where you can play around and see if D2 fits your taste and use case.

Playground URL: https://d2play.yashbhalodi.me/
