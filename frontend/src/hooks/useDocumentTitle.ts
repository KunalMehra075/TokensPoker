import { useEffect } from 'react'
import { BRAND } from '@/constants'

// Lightweight per-route title + meta description management for the SPA. The
// crawl-critical tags live statically in index.html; this keeps tab titles and
// descriptions accurate as users navigate.
export function useDocumentTitle(title: string, description?: string) {
  useEffect(() => {
    document.title = title ? `${title} | ${BRAND.name}` : BRAND.name
    if (description) {
      let tag = document.querySelector('meta[name="description"]')
      if (!tag) {
        tag = document.createElement('meta')
        tag.setAttribute('name', 'description')
        document.head.appendChild(tag)
      }
      tag.setAttribute('content', description)
    }
  }, [title, description])
}
