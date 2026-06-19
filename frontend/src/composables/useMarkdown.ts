import MarkdownIt from 'markdown-it'
import DOMPurify from 'dompurify'

const md = new MarkdownIt({
  html: false, // 生 HTML を無効化（XSS 対策の一次防御。ユーザ入力HTMLはエスケープされる）
  linkify: true,
  breaks: true,
})

/**
 * Discord 風の spoiler 記法 ||text|| を span に変換する。
 * markdown-it の text ルールが `|` を特殊文字として扱わないためインラインルール化が難しく、
 * レンダリング後の HTML に対する後処理で実装する（Terms は config 由来の信頼コンテンツ）。
 */
function applySpoilers(html: string): string {
  return html.replace(/\|\|([^|]+?)\|\|/g, '<span class="spoiler">$1</span>')
}

/** Markdown を HTML 文字列へ変換する（DOMPurify 未適用。構造テスト用） */
export function renderMarkdownRaw(src: string): string {
  return applySpoilers(md.render(src ?? ''))
}

/** Markdown を Discord 風にレンダリングし、DOMPurify で無害化した HTML を返す */
export function renderMarkdown(src: string): string {
  return DOMPurify.sanitize(renderMarkdownRaw(src), {
    ADD_ATTR: ['class', 'target', 'rel'],
  })
}

export function useMarkdown() {
  return { render: renderMarkdown }
}
