import type { ReactNode } from 'react'
import { MarketingHeader } from '@/components/layout/MarketingHeader'
import { Footer } from '@/components/layout/Footer'
import { BRAND } from '@/constants'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'

function LegalLayout({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="flex min-h-screen flex-col">
      <MarketingHeader />
      <main className="flex-1">
        <div className="container-app max-w-3xl py-16">
          <h1 className="text-3xl font-bold tracking-tight text-foreground">{title}</h1>
          <div className="mt-8 space-y-5 text-[15px] leading-relaxed text-muted-foreground [&_h2]:mt-8 [&_h2]:text-lg [&_h2]:font-semibold [&_h2]:text-foreground [&_a]:text-primary [&_a]:underline">
            {children}
          </div>
        </div>
      </main>
      <Footer />
    </div>
  )
}

export function PrivacyPage() {
  useDocumentTitle('Privacy Policy', `How ${BRAND.name} handles your data.`)
  return (
    <LegalLayout title="Privacy Policy">
      <p>
        This policy explains what {BRAND.name} collects and why. The short version: we collect
        as little as we can get away with, because storing data you do not need is a liability,
        not a feature.
      </p>
      <h2>What we collect</h2>
      <p>
        When you join, you give us a display name and an email address. The email is an identity
        label so your rooms, votes, and history belong to you across sessions. We do not use it
        as a password and we do not sell it. We also store the estimation data you create:
        rooms, tasks, the cards you select, and final decisions.
      </p>
      <h2>Cookies and local storage</h2>
      <p>
        We keep a small identity token in your browser local storage so you stay signed in. If
        we add privacy-friendly analytics (such as Plausible) or, much later, advertising on our
        public marketing pages (such as Google AdSense), the relevant cookies and their purpose
        will be disclosed here before they go live.
      </p>
      <h2>How we use your data</h2>
      <p>
        To run the product: show your rooms, sync votes in real time, and keep your history.
        That is it. We may review aggregate, anonymized estimation trends to improve the
        product, never tied back to an individual without consent.
      </p>
      <h2>Your choices</h2>
      <p>
        You can request deletion of your account and associated data at any time by emailing{' '}
        <a href={`mailto:${BRAND.contactEmail}`}>{BRAND.contactEmail}</a>.
      </p>
      <h2>Contact</h2>
      <p>
        Questions about privacy? Reach us at{' '}
        <a href={`mailto:${BRAND.contactEmail}`}>{BRAND.contactEmail}</a>.
      </p>
    </LegalLayout>
  )
}

export function TermsPage() {
  useDocumentTitle('Terms of Service', `The rules for using ${BRAND.name}.`)
  return (
    <LegalLayout title="Terms of Service">
      <p>
        By using {BRAND.name} you agree to these terms. They are intentionally short and written
        in plain language.
      </p>
      <h2>Using the service</h2>
      <p>
        {BRAND.name} is a collaborative estimation tool. Use it for lawful purposes. Do not abuse
        it, attempt to break it, or use it to harass other people. We may suspend access that
        threatens the service or other users.
      </p>
      <h2>Your content</h2>
      <p>
        You own the estimation content you create. You grant us the limited permission needed to
        store and display it back to you and your room members so the product can function.
      </p>
      <h2>Availability</h2>
      <p>
        The service is provided as is, on a best-effort basis. We aim for high availability but
        do not promise it will never go down. For a free tool, that is a fair trade.
      </p>
      <h2>Changes</h2>
      <p>
        We may update these terms as the product evolves. Continued use after a change means you
        accept the updated terms.
      </p>
      <h2>Contact</h2>
      <p>
        Questions? Email <a href={`mailto:${BRAND.contactEmail}`}>{BRAND.contactEmail}</a>.
      </p>
    </LegalLayout>
  )
}

export function AboutPage() {
  useDocumentTitle('About', `Why we built ${BRAND.name}.`)
  return (
    <LegalLayout title="About FreeTokensPoker">
      <p>
        {BRAND.name} started from a small, repeated annoyance. Every sprint planning call had the
        same unscripted moment: someone says we will probably use Claude, someone else says GPT
        is fine, a third person mentions cost, and then everyone moves on without writing a
        single thing down.
      </p>
      <p>
        Planning Poker solved this years ago for engineering effort. A private vote, a
        simultaneous reveal, an honest conversation about why people disagree. We thought the
        same ritual deserved to cover the new questions AI introduced: how many tokens, which
        model, what it will cost, how long it takes when an agent is doing half the work.
      </p>
      <p>
        So we kept the ritual and changed the cards. The goal is not to predict the future of AI
        spending with spreadsheet precision. The goal is to get a team aligned, on the record,
        in under a minute, before the work begins.
      </p>
      <p>
        We are keeping the core free because adoption matters more than early revenue. If you
        run a team that builds with AI, we would love for {BRAND.name} to become a small part of
        your sprint ritual.
      </p>
    </LegalLayout>
  )
}

export function ContactPage() {
  useDocumentTitle('Contact', `Get in touch with the ${BRAND.name} team.`)
  return (
    <LegalLayout title="Contact us">
      <p>
        We read everything. Feature requests, bug reports, sharp criticism, and the occasional
        kind word are all welcome.
      </p>
      <h2>Email</h2>
      <p>
        The fastest way to reach us is{' '}
        <a href={`mailto:${BRAND.contactEmail}`}>{BRAND.contactEmail}</a>. We try to reply within
        a couple of business days.
      </p>
      <h2>What helps</h2>
      <p>
        If you are reporting a problem, tell us what you expected, what happened instead, and the
        room or browser you were using. Specifics turn a vague report into a fixed bug.
      </p>
    </LegalLayout>
  )
}
