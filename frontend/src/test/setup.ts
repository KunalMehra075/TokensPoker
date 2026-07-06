// Registers jest-dom matchers (toBeInTheDocument, toHaveClass, ...) and clears
// the DOM + persisted stores between tests so cases stay isolated.
import '@testing-library/jest-dom/vitest'
import { afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'

afterEach(() => {
  cleanup()
  localStorage.clear()
})
