import { useEffect } from 'react'
import { BRAND } from '@/constants'

// Lightweight per-route SEO tag management for the SPA. The crawl-critical tags
// live statically in index.html; this keeps the title, description, and the
// self-referential canonical / og:url accurate as users navigate between routes
// (index.html only ships the "/" canonical, so subpages must override it).
export function useDocumentTitle(title: string, description?: string) {
  useEffect(() => {
    document.title = title ? `${title} | ${BRAND.name}` : BRAND.name

    if (description) {
      setMeta('name', 'description', description)
      setMeta('property', 'og:description', description)
    }

    // Self-referential canonical + og:url so /about, /contact, /privacy, /terms
    // point at themselves instead of all canonicalizing to the homepage.
    const url = `${BRAND.url}${window.location.pathname}`
    setLink('canonical', url)
    setMeta('property', 'og:url', url)
  }, [title, description])
}

function setMeta(attr: 'name' | 'property', key: string, content: string) {
  let tag = document.querySelector(`meta[${attr}="${key}"]`)
  if (!tag) {
    tag = document.createElement('meta')
    tag.setAttribute(attr, key)
    document.head.appendChild(tag)
  }
  tag.setAttribute('content', content)
}

function setLink(rel: string, href: string) {
  let tag = document.querySelector(`link[rel="${rel}"]`)
  if (!tag) {
    tag = document.createElement('link')
    tag.setAttribute('rel', rel)
    document.head.appendChild(tag)
  }
  tag.setAttribute('href', href)
}
