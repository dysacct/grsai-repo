# 客户端接入

只要客户端支持自定义 `baseURL`，通常就可以把 OpenAI 风格请求指向本服务。

## JavaScript

```ts
import OpenAI from 'openai'

const client = new OpenAI({
  apiKey: 'local-dev-token',
  baseURL: 'http://localhost:8080/v1'
})

const completion = await client.chat.completions.create({
  model: 'gpt-5.5',
  messages: [
    { role: 'user', content: '用一句话介绍这个项目' }
  ]
})

console.log(completion.choices[0]?.message?.content)
```

当前服务不会校验客户端传入的 OpenAI API Key，真正的上游鉴权由服务端环境变量 `GRSAI_API_KEY` 完成。生产环境如果暴露给外部用户，建议在网关或服务内补充访问控制。

## Python

```py
from openai import OpenAI

client = OpenAI(
    api_key="local-dev-token",
    base_url="http://localhost:8080/v1",
)

response = client.chat.completions.create(
    model="gpt-5.5",
    messages=[
        {"role": "user", "content": "给我一个接口接入检查清单"}
    ],
)

print(response.choices[0].message.content)
```

## 流式输出

```ts
const stream = await client.chat.completions.create({
  model: 'gpt-5.5',
  stream: true,
  messages: [
    { role: 'user', content: '流式输出一段部署建议' }
  ]
})

for await (const part of stream) {
  process.stdout.write(part.choices[0]?.delta?.content ?? '')
}
```
