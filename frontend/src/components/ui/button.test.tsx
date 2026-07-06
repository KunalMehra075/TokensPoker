import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Button } from './button'

describe('Button', () => {
  it('renders its children', () => {
    render(<Button>Start a room</Button>)
    expect(screen.getByRole('button', { name: 'Start a room' })).toBeInTheDocument()
  })

  it('applies the danger variant classes', () => {
    render(<Button variant="danger">Delete</Button>)
    expect(screen.getByRole('button', { name: 'Delete' })).toHaveClass('bg-danger')
  })

  it('fires onClick when pressed', async () => {
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Vote</Button>)
    await userEvent.click(screen.getByRole('button', { name: 'Vote' }))
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('is disabled and unclickable while loading', async () => {
    const onClick = vi.fn()
    render(
      <Button loading onClick={onClick}>
        Reveal
      </Button>,
    )
    const btn = screen.getByRole('button')
    expect(btn).toBeDisabled()
    await userEvent.click(btn)
    expect(onClick).not.toHaveBeenCalled()
  })

  it('respects an explicit disabled prop', () => {
    render(<Button disabled>Nope</Button>)
    expect(screen.getByRole('button', { name: 'Nope' })).toBeDisabled()
  })
})
