import { useRef, useState } from 'react'
import { calculate, type BinaryOperation } from './api'

const maxInputLength = 16

type Phase = 'typing' | 'operator' | 'result'

type Repeat = {
  operation: BinaryOperation
  operand: number
}

const symbols: Record<BinaryOperation, string> = {
  add: '+',
  subtract: '−',
  multiply: '×',
  divide: '÷',
  power: 'x^y',
  percentage: '%',
}

function formatResult(value: number): string {
  if (Object.is(value, -0)) {
    return '0'
  }
  return String(value)
}

function parseDisplay(value: string): number | null {
  if (value === '.' || value === '-' || value === '-.') {
    return null
  }
  const parsed = Number(value)
  if (!Number.isFinite(parsed)) {
    return null
  }
  return parsed
}

function errorMessage(error: unknown): string {
  if (error instanceof Error && error.message) {
    return error.message
  }
  return 'Calculation failed'
}

export function useCalculator() {
  const [display, setDisplay] = useState('0')
  const [stored, setStored] = useState<number | null>(null)
  const [pending, setPending] = useState<BinaryOperation | null>(null)
  const [phase, setPhase] = useState<Phase>('typing')
  const [repeat, setRepeat] = useState<Repeat | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const busy = useRef(false)
  const requestId = useRef(0)

  const preview = pending !== null && stored !== null ? `${formatResult(stored)} ${symbols[pending]}` : ''

  function resetEntryState() {
    setStored(null)
    setPending(null)
    setRepeat(null)
    setPhase('result')
  }

  function guard(task: (isCurrent: () => boolean) => Promise<void>) {
    if (busy.current) {
      return
    }
    busy.current = true
    const id = requestId.current + 1
    requestId.current = id
    setLoading(true)
    setError(null)
    const isCurrent = () => id === requestId.current

    void task(isCurrent)
      .catch((caught: unknown) => {
        if (!isCurrent()) {
          return
        }
        setError(errorMessage(caught))
        resetEntryState()
      })
      .finally(() => {
        if (!isCurrent()) {
          return
        }
        busy.current = false
        setLoading(false)
      })
  }

  function inputDigit(digit: string) {
    if (busy.current) {
      return
    }
    setError(null)
    setDisplay((current) => {
      if (phase === 'operator' || phase === 'result') {
        return digit
      }
      if (current === '0') {
        return digit
      }
      if (current.replace('-', '').length >= maxInputLength) {
        return current
      }
      return current + digit
    })
    if (phase === 'result' && pending === null) {
      setStored(null)
      setRepeat(null)
    }
    if (phase !== 'typing') {
      setPhase('typing')
    }
  }

  function inputDecimal() {
    if (busy.current) {
      return
    }
    setError(null)
    setDisplay((current) => {
      if (phase === 'operator' || phase === 'result') {
        return '0.'
      }
      if (current.includes('.')) {
        return current
      }
      return `${current}.`
    })
    if (phase === 'result' && pending === null) {
      setStored(null)
      setRepeat(null)
    }
    if (phase !== 'typing') {
      setPhase('typing')
    }
  }

  function deleteDigit() {
    if (busy.current) {
      return
    }
    setError(null)
    if (phase === 'operator') {
      return
    }
    if (phase === 'result' && pending === null) {
      setStored(null)
      setRepeat(null)
    }
    setDisplay((current) => {
      const next = current.length <= 1 ? '0' : current.slice(0, -1)
      return next === '-' ? '0' : next
    })
    setPhase('typing')
  }

  function clear() {
    requestId.current += 1
    busy.current = false
    setLoading(false)
    setDisplay('0')
    setStored(null)
    setPending(null)
    setRepeat(null)
    setError(null)
    setPhase('typing')
  }

  function pressOperator(operation: BinaryOperation) {
    if (busy.current) {
      return
    }
    const current = parseDisplay(display)
    if (current === null) {
      setError('Enter a number first')
      return
    }
    setError(null)

    // A binary operator only records the left operand and operation.
    // The API request waits until equals, including when the operator is changed
    // after the right operand has been typed.
    if (pending !== null && phase === 'operator') {
      setPending(operation)
      setRepeat(null)
      return
    }
    if (pending !== null && stored !== null && phase === 'typing') {
      setPending(operation)
      setRepeat(null)
      setPhase('operator')
      return
    }

    setStored(current)
    setPending(operation)
    setRepeat(null)
    setPhase('operator')
  }

  function pressEquals() {
    if (busy.current) {
      return
    }
    const current = parseDisplay(display)
    if (current === null) {
      setError('Enter a number first')
      return
    }

    if (pending !== null && stored !== null) {
      const left = stored
      const operation = pending
      guard(async (isCurrent) => {
        const result = await calculate(operation, left, current)
        if (!isCurrent()) {
          return
        }
        setDisplay(formatResult(result))
        setStored(result)
        setRepeat({ operation, operand: current })
        setPending(null)
        setPhase('result')
      })
      return
    }

    if (repeat !== null && stored !== null && phase === 'result') {
      const { operation, operand } = repeat
      const left = stored
      guard(async (isCurrent) => {
        const result = await calculate(operation, left, operand)
        if (!isCurrent()) {
          return
        }
        setDisplay(formatResult(result))
        setStored(result)
        setPhase('result')
      })
    }
  }

  function pressSqrt() {
    if (busy.current) {
      return
    }
    const current = parseDisplay(display)
    if (current === null) {
      setError('Enter a number first')
      return
    }
    guard(async (isCurrent) => {
      const result = await calculate('sqrt', current)
      if (!isCurrent()) {
        return
      }
      setDisplay(formatResult(result))
      setPhase('result')
    })
  }

  return {
    display,
    preview,
    error,
    loading,
    pending,
    inputDigit,
    inputDecimal,
    deleteDigit,
    clear,
    pressOperator,
    pressEquals,
    pressSqrt,
  }
}
