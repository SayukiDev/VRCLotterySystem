import { describe, expect, it } from 'vitest'
import { renderMarkdown, renderMarkdownRaw } from './useMarkdown'

describe('renderMarkdownRaw（構造）', () => {
  it('基本的な Markdown を HTML に変換する', () => {
    const html = renderMarkdownRaw('**bold** and *italic*')
    expect(html).toContain('<strong>bold</strong>')
    expect(html).toContain('<em>italic</em>')
  })

  it('リストを変換する', () => {
    const html = renderMarkdownRaw('- a\n- b')
    expect(html).toContain('<ul>')
    expect(html).toContain('<li>a</li>')
  })

  it('||spoiler|| を spoiler span に変換する', () => {
    const html = renderMarkdownRaw('秘密は ||これ|| です')
    expect(html).toContain('<span class="spoiler">これ</span>')
  })

  it('html:false により生 HTML はエスケープされ実行可能な script タグは生成されない', () => {
    const html = renderMarkdownRaw('hello <script>alert(1)</script>')
    expect(html).not.toContain('<script>')
    expect(html).toContain('&lt;script&gt;')
  })
})

describe('renderMarkdown（無害化）', () => {
  it('実行可能な script タグを含まない', () => {
    const html = renderMarkdown('hello <script>alert(1)</script>')
    expect(html).not.toContain('<script>')
  })

  it('通常テキストは保持される', () => {
    const html = renderMarkdown('こんにちは')
    expect(html).toContain('こんにちは')
  })
})
