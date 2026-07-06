import { Link } from 'react-router-dom'
import {
  ArrowRight,
  Coins,
  DollarSign,
  CalendarDays,
  Cpu,
  EyeOff,
  Users,
  Zap,
  History as HistoryIcon,
} from 'lucide-react'
import { MarketingHeader } from '@/components/layout/MarketingHeader'
import { Footer } from '@/components/layout/Footer'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useDocumentTitle } from '@/hooks/useDocumentTitle'

const MODES = [
  { icon: Coins, name: 'AI Tokens', sample: '500K to 50M', blurb: 'How much context and generation will this actually burn through.' },
  { icon: DollarSign, name: 'AI Cost', sample: '$1 to $250', blurb: 'The number finance will eventually ask you about. Beat them to it.' },
  { icon: CalendarDays, name: 'Engineering Days', sample: '1 to 21', blurb: 'Familiar Fibonacci effort, now with AI doing half the typing.' },
  { icon: Cpu, name: 'Best AI Model', sample: 'GPT, Claude, Gemini', blurb: 'Settle the model debate with a vote instead of the loudest voice.' },
]

const STEPS = [
  { icon: Users, title: 'Create a room', body: 'Name it, share the six character code. Your team is in within seconds, no installs and no calendar invite archaeology.' },
  { icon: EyeOff, title: 'Everyone votes in private', body: 'Each person picks a card without seeing anyone else. The senior engineer cannot anchor the room before the intern even reads the ticket.' },
  { icon: Zap, title: 'Reveal together', body: 'One click flips every card at once. The interesting part is not the average, it is the two people who are wildly far apart.' },
  { icon: HistoryIcon, title: 'Decide and keep the record', body: 'Lock in a final number, archive the round, move on. Every estimate is saved so future you can see how wrong past you was.' },
]

const FAQ = [
  {
    q: 'What is FreeTokensPoker?',
    a: 'It is Planning Poker rebuilt for teams that build with AI. Instead of only estimating story points, your team estimates AI tokens, cost, engineering days, and which model to use, all through the same private vote then simultaneous reveal ritual you already know.',
  },
  {
    q: 'Do I need an account or a password?',
    a: 'No password, ever. You enter your name and email once so your votes and history have a home. That is the entire signup. Email is an identity label here, not a security gate.',
  },
  {
    q: 'Is it really free?',
    a: 'Yes. Unlimited rooms and unlimited estimations on the free tier. The name is about access, not a countdown to a paywall. Future team and enterprise features may be paid, but the core estimation workflow stays free.',
  },
  {
    q: 'How many people can join a room?',
    a: 'Bring the whole squad. Rooms are designed for a normal engineering team, and everything updates live over a realtime connection, so votes and reveals appear instantly for everyone.',
  },
  {
    q: 'Why estimate AI usage at all when AI is cheap right now?',
    a: 'Because cheap is a temporary launch promotion, not a law of physics. Teams that build the habit of discussing model choice and token budget now will not panic when the invoices start to matter.',
  },
  {
    q: 'Can I see past estimations?',
    a: 'Every completed round is stored with the task, the mode, individual votes, and the final decision. Your history is one click away from any room.',
  },
]

export function LandingPage() {
  useDocumentTitle(
    '',
    'FreeTokensPoker is Planning Poker for AI teams. Estimate AI tokens, cost, model choice, and engineering effort together with a private vote and a simultaneous reveal.',
  )

  return (
    <div className="flex min-h-screen flex-col">
      <MarketingHeader />

      <main className="flex-1">
        {/* Hero */}
        <section className="container-app pt-16 pb-20 sm:pt-24">
          <div className="mx-auto max-w-3xl text-center">
            <Badge variant="primary" className="mb-5">
              <Zap className="h-3.5 w-3.5" />
              AI Planning Poker
            </Badge>
            <h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-5xl">
              Estimate the AI work before you spend the AI budget
            </h1>
            <p className="mx-auto mt-5 max-w-2xl text-lg text-muted-foreground">
              Your team already argues about which model to use and how many tokens a feature
              will eat. FreeTokensPoker turns that hallway debate into a quick, structured vote.
              Everyone estimates in private, then you reveal together and actually talk about
              the differences.
            </p>
            <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <Link to="/login">
                <Button size="lg">
                  Start a room
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
              <a href="#how-it-works">
                <Button size="lg" variant="secondary">
                  See how it works
                </Button>
              </a>
            </div>
            <p className="mt-4 text-sm text-muted-foreground">
              No password, no setup, no credit card. You will be voting in under a minute.
            </p>
          </div>
        </section>

        {/* Modes */}
        <section id="modes" className="border-t border-border bg-muted/40 py-20">
          <div className="container-app">
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="text-3xl font-bold tracking-tight text-foreground">
                Four ways to estimate, one familiar ritual
              </h2>
              <p className="mt-3 text-muted-foreground">
                Pick the dimension that matters for the decision in front of you. The cards
                change, the muscle memory does not.
              </p>
            </div>
            <div className="mt-12 grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
              {MODES.map((m) => (
                <div
                  key={m.name}
                  className="rounded-lg border border-border bg-white p-6 shadow-card transition-shadow hover:shadow-card-hover"
                >
                  <span className="inline-flex h-10 w-10 items-center justify-center rounded-md bg-blue-50 text-primary">
                    <m.icon className="h-5 w-5" />
                  </span>
                  <h3 className="mt-4 font-semibold text-foreground">{m.name}</h3>
                  <p className="mt-1 text-sm font-medium text-primary">{m.sample}</p>
                  <p className="mt-2 text-sm text-muted-foreground">{m.blurb}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* How it works */}
        <section id="how-it-works" className="py-20">
          <div className="container-app">
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="text-3xl font-bold tracking-tight text-foreground">
                Planning Poker, with the AI parts added
              </h2>
              <p className="mt-3 text-muted-foreground">
                If your team has ever held up fingers or flipped cards in a sprint planning
                call, you already know how to use this.
              </p>
            </div>
            <div className="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-4">
              {STEPS.map((s, i) => (
                <div key={s.title} className="relative rounded-lg border border-border bg-white p-6 shadow-card">
                  <span className="text-sm font-bold text-primary">0{i + 1}</span>
                  <span className="mt-3 inline-flex h-10 w-10 items-center justify-center rounded-md bg-blue-50 text-primary">
                    <s.icon className="h-5 w-5" />
                  </span>
                  <h3 className="mt-4 font-semibold text-foreground">{s.title}</h3>
                  <p className="mt-2 text-sm text-muted-foreground">{s.body}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Why */}
        <section className="border-t border-border bg-muted/40 py-20">
          <div className="container-app grid gap-10 lg:grid-cols-2 lg:items-center">
            <div>
              <h2 className="text-3xl font-bold tracking-tight text-foreground">
                Because nobody is writing down the AI decisions
              </h2>
              <div className="mt-5 space-y-4 text-muted-foreground">
                <p>
                  Right now the planning meeting sounds like this. We will probably use Claude.
                  I think GPT is enough. Let us just ask Cursor. Then everyone moves on, nothing
                  gets recorded, and three weeks later someone asks why the AI bill tripled.
                </p>
                <p>
                  FreeTokensPoker gives that conversation a shape. A private vote removes the
                  anchoring problem where the most senior person in the room accidentally
                  decides for everyone. The simultaneous reveal surfaces disagreement instead of
                  burying it. And the saved history means your estimates slowly turn into real
                  data about how your team actually thinks about AI.
                </p>
                <p>
                  We are not trying to replace your engineering judgment. We are trying to get
                  it out of three peoples heads and onto the table where the whole team can see
                  it.
                </p>
              </div>
            </div>
            <div className="rounded-lg border border-border bg-white p-8 shadow-card">
              <div className="grid grid-cols-3 gap-3">
                {['500K', '1M', '2M', '5M', '10M', '20M'].map((c, i) => (
                  <div
                    key={c}
                    className={`flex h-20 items-center justify-center rounded-md border text-lg font-bold ${
                      i === 2
                        ? 'border-primary bg-primary text-primary-foreground'
                        : 'border-border bg-muted/40 text-foreground'
                    }`}
                  >
                    {c}
                  </div>
                ))}
              </div>
              <p className="mt-5 text-center text-sm text-muted-foreground">
                Tokens mode. Everyone has the same deck. Only the conversation is different.
              </p>
            </div>
          </div>
        </section>

        {/* FAQ */}
        <section className="py-20">
          <div className="container-app mx-auto max-w-3xl">
            <h2 className="text-center text-3xl font-bold tracking-tight text-foreground">
              Questions teams actually ask
            </h2>
            <div className="mt-10 divide-y divide-border">
              {FAQ.map((item) => (
                <details key={item.q} className="group py-5">
                  <summary className="flex cursor-pointer list-none items-center justify-between text-base font-semibold text-foreground">
                    {item.q}
                    <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-open:rotate-90" />
                  </summary>
                  <p className="mt-3 text-sm leading-relaxed text-muted-foreground">{item.a}</p>
                </details>
              ))}
            </div>
          </div>
        </section>

        {/* CTA */}
        <section className="border-t border-border py-20">
          <div className="container-app">
            <div className="mx-auto max-w-2xl rounded-lg border border-border bg-foreground p-10 text-center">
              <h2 className="text-3xl font-bold tracking-tight text-white">
                Your next sprint has AI in it. Plan for it.
              </h2>
              <p className="mt-3 text-gray-300">
                Spin up a room, drop the code in your team channel, and run your first AI
                estimation round in the time it takes to lose this tab.
              </p>
              <div className="mt-7 flex justify-center">
                <Link to="/login">
                  <Button size="lg" className="bg-white text-foreground hover:bg-gray-100">
                    Start a room
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  )
}
