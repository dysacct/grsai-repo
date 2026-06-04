import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'GRSAI NewAPI',
  description: 'OpenAI-compatible API adapter for GRS AI chat and image models.',
  lang: 'zh-CN',
  cleanUrls: true,
  lastUpdated: true,
  metaChunk: true,
  themeConfig: {
    logo: { src: '/logo.svg', alt: 'GRSAI NewAPI' },
    siteTitle: 'GRSAI NewAPI',
    nav: [
      { text: '指南', link: '/guide/getting-started' },
      { text: '接口', link: '/api/overview' },
      { text: '运维', link: '/ops/deploy' }
    ],
    sidebar: [
      {
        text: '指南',
        items: [
          { text: '快速开始', link: '/guide/getting-started' },
          { text: '配置说明', link: '/guide/configuration' },
          { text: '客户端接入', link: '/guide/client-examples' }
        ]
      },
      {
        text: '接口',
        items: [
          { text: '接口总览', link: '/api/overview' },
          { text: '模型列表', link: '/api/models' },
          { text: '聊天补全', link: '/api/chat-completions' },
          { text: '图片生成与编辑', link: '/api/images' }
        ]
      },
      {
        text: '运维',
        items: [
          { text: '部署', link: '/ops/deploy' },
          { text: '故障排查', link: '/ops/troubleshooting' }
        ]
      }
    ],
    search: {
      provider: 'local'
    },
    outline: {
      label: '本页目录',
      level: [2, 3]
    },
    docFooter: {
      prev: '上一页',
      next: '下一页'
    },
    lastUpdated: {
      text: '最后更新',
      formatOptions: {
        dateStyle: 'medium',
        timeStyle: 'short'
      }
    },
    returnToTopLabel: '回到顶部',
    sidebarMenuLabel: '菜单',
    darkModeSwitchLabel: '外观',
    lightModeSwitchTitle: '切换到浅色模式',
    darkModeSwitchTitle: '切换到深色模式'
  }
})
