// The Go API accepts POST /api/calculate.
// Binary operations send { operation, a, b }. Square root sends { operation: "sqrt", a }.
// Success is { result }. Failure is { code, message }.

export type BinaryOperation = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'percentage'

export type Operation = BinaryOperation | 'sqrt'

type CalculatePayload = {
  operation: Operation
  a: number
  b?: number
}

type SuccessBody = {
  result: number
}

type ErrorBody = {
  code?: string
  message?: string
}

export class CalculateError extends Error {
  readonly code: string

  constructor(code: string, message: string) {
    super(message)
    this.name = 'CalculateError'
    this.code = code
  }
}

export async function calculate(operation: Operation, a: number, b?: number): Promise<number> {
  const payload: CalculatePayload = { operation, a }
  if (b !== undefined) {
    payload.b = b
  }

  let response: Response
  try {
    response = await fetch('/api/calculate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
  } catch {
    throw new CalculateError('network_error', 'Could not reach the calculator service')
  }

  const body = await readBody(response)
  if (!response.ok) {
    const errorBody = body as ErrorBody | null
    const code = errorBody?.code || 'request_failed'
    const message = errorBody?.message || 'Calculation failed'
    throw new CalculateError(code, message)
  }

  const result = (body as SuccessBody | null)?.result
  if (typeof result !== 'number' || !Number.isFinite(result)) {
    throw new CalculateError('invalid_response', 'Calculator service returned an unreadable response')
  }
  return result
}

async function readBody(response: Response): Promise<unknown> {
  try {
    return await response.json()
  } catch {
    return null
  }
}
