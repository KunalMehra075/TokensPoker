import { useQuery } from '@tanstack/react-query'
import { fetchModes } from '@/services/api'
import type { EstimationMode, ModeDefinition } from '@/types'

// Estimation modes come from the backend so cards are never hardcoded in the UI.
export function useModes() {
  return useQuery({
    queryKey: ['modes'],
    queryFn: fetchModes,
    staleTime: Infinity,
  })
}

export function findMode(modes: ModeDefinition[] | undefined, mode: EstimationMode) {
  return modes?.find((m) => m.mode === mode)
}
