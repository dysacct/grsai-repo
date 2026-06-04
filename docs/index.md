---
layout: home

hero:
  name: GRSAI NewAPI
  text: OpenAI-compatible adapter for GRS AI.
  tagline: 一个轻量、清晰、可直接接入的 Go API 转换层，把 GRS AI 的聊天与绘图能力整理成常见的 OpenAI 风格接口。
  image:
    src: /logo.svg
    alt: GRSAI NewAPI
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: 查看接口
      link: /api/overview

features:
  - title: OpenAI 风格接口
    details: 提供 /v1/models、/v1/chat/completions、/v1/images/generations、/v1/images/edits，方便现有客户端迁移。
  - title: 流式聊天转发
    details: 对 stream=true 的聊天请求透传 SSE 响应，适合聊天 UI、长文本生成和实时输出。
  - title: 图片能力桥接
    details: 将 OpenAI 图片参数映射到 GRS AI 绘图接口，支持文生图和图生图。
---

<div class="signal-row">
  <div><strong>4 个公开接口</strong><span>模型、聊天、文生图、图生图</span></div>
  <div><strong>1 个核心密钥</strong><span>通过 GRSAI_API_KEY 接入上游</span></div>
  <div><strong>Go 标准路由</strong><span>少依赖、易部署、可审计</span></div>
</div>

## 这个项目适合做什么

GRSAI NewAPI 适合放在你的服务侧，作为客户端和 GRS AI 之间的兼容层。客户端继续按 OpenAI 风格请求 `/v1/chat/completions` 或图片接口，服务端负责补齐鉴权、参数映射、SSE 透传和错误整理。

## 当前能力

<div class="api-grid">
  <div class="api-card">
    <strong>模型列表</strong>
    <code>GET /v1/models</code>
    <p>返回当前服务声明支持的聊天模型和图片模型。</p>
  </div>
  <div class="api-card">
    <strong>聊天补全</strong>
    <code>POST /v1/chat/completions</code>
    <p>支持普通 JSON 响应，也支持 <code>stream: true</code> 的 SSE 流式响应。</p>
  </div>
  <div class="api-card">
    <strong>图片生成</strong>
    <code>POST /v1/images/generations</code>
    <p>接收 prompt、model、size、quality 等参数，返回图片 URL 或 base64。</p>
  </div>
  <div class="api-card">
    <strong>图片编辑</strong>
    <code>POST /v1/images/edits</code>
    <p>支持 JSON base64/URL 输入，也支持 multipart/form-data 上传图片。</p>
  </div>
</div>

## 下一步

从 [快速开始](/guide/getting-started) 启动服务，再按 [客户端接入](/guide/client-examples) 里的示例把现有 OpenAI SDK 指向你的 NewAPI 地址。
