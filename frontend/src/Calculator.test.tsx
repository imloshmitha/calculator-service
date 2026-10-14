import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, afterEach, expect, test, vi } from 'vitest'
import App from './App'

function jsonResponse(body: unknown, status = 200) {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

async function press(name: string) {
  await userEvent.click(screen.getByRole('button', { name }))
}

function result() {
  return screen.getByRole('status', { name: 'Result' })
}

function requestBody(callIndex = 0) {
  const fetchMock = vi.mocked(fetch)
  const init = fetchMock.mock.calls[callIndex][1]
  return JSON.parse(String(init?.body)) as { operation: string; a: number; b?: number }
}

beforeEach(() => {
  vi.stubGlobal('fetch', vi.fn())
})

afterEach(() => {
  vi.unstubAllGlobals()
})

test('renders a zeroed calculator and its controls', () => {
  render(<App />)

  expect(result()).toHaveTextContent('0')
  expect(screen.getByRole('button', { name: 'Clear' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Decimal' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Add' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Subtract' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Multiply' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Divide' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Power' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Square root' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Percentage' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: 'Equals' })).toBeInTheDocument()
})

test('enters digits and ignores a repeated decimal point', async () => {
  render(<App />)

  await press('1')
  await press('9')
  await press('Decimal')
  await press('Decimal')
  await press('5')

  expect(result()).toHaveTextContent('19.5')
  expect(fetch).not.toHaveBeenCalled()
})

test('shows the API result instead of calculating locally', async () => {
  vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 123 }))
  render(<App />)

  await press('1')
  await press('9')
  await press('Add')
  await press('9')
  await press('Equals')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('123')
  expect(fetch).toHaveBeenCalledTimes(1)
  expect(requestBody()).toEqual({ operation: 'add', a: 19, b: 9 })
  expect(vi.mocked(fetch).mock.calls[0][0]).toBe('/api/calculate')
})

test('does not call the API for a binary operation until equals', async () => {
  vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 123 }))
  render(<App />)

  await press('1')
  await press('9')
  await press('Add')
  await press('9')

  expect(fetch).not.toHaveBeenCalled()
  expect(result()).toHaveTextContent('9')

  await press('Equals')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('123')
  expect(fetch).toHaveBeenCalledTimes(1)
  expect(requestBody()).toEqual({ operation: 'add', a: 19, b: 9 })
})

test.each([
  ['Subtract', 'subtract'],
  ['Multiply', 'multiply'],
  ['Divide', 'divide'],
  ['Power', 'power'],
])('%s waits for equals before calling the API', async (button, operation) => {
  vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 7 }))
  render(<App />)

  await press('8')
  await press(button)
  await press('2')

  expect(fetch).not.toHaveBeenCalled()
  expect(result()).toHaveTextContent('2')

  await press('Equals')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('7')
  expect(requestBody()).toEqual({ operation, a: 8, b: 2 })
})

test('continues from the previous API result after equals', async () => {
  vi.mocked(fetch)
    .mockResolvedValueOnce(jsonResponse({ result: 100 }))
    .mockResolvedValueOnce(jsonResponse({ result: 94 }))
  render(<App />)

  await press('1')
  await press('9')
  await press('Add')
  await press('9')
  await press('Equals')
  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('100')
  expect(fetch).toHaveBeenCalledTimes(1)

  await press('Subtract')
  await press('6')
  expect(fetch).toHaveBeenCalledTimes(1)
  expect(result()).toHaveTextContent('6')

  await press('Equals')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('94')
  expect(requestBody(0)).toEqual({ operation: 'add', a: 19, b: 9 })
  expect(requestBody(1)).toEqual({ operation: 'subtract', a: 100, b: 6 })
})

test('asks the API for b percent of a', async () => {
  vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 30 }))
  render(<App />)

  await press('2')
  await press('0')
  await press('0')
  await press('Percentage')
  await press('1')
  await press('5')
  await press('Equals')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('30')
  expect(requestBody()).toEqual({ operation: 'percentage', a: 200, b: 15 })
})

test('sends square root with only operand a', async () => {
  vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 3 }))
  render(<App />)

  await press('9')
  await press('Square root')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('3')
  expect(requestBody()).toEqual({ operation: 'sqrt', a: 9 })
})

test('replacing the operator before a second number does not call the API', async () => {
  vi.mocked(fetch).mockResolvedValueOnce(jsonResponse({ result: 6 }))
  render(<App />)

  await press('8')
  await press('Add')
  await press('Subtract')
  await press('2')
  await press('Equals')

  expect(await screen.findByRole('status', { name: 'Result' })).toHaveTextContent('6')
  expect(fetch).toHaveBeenCalledTimes(1)
  expect(requestBody()).toEqual({ operation: 'subtract', a: 8, b: 2 })
})

test('shows an API error for division by zero and can be cleared', async () => {
  vi.mocked(fetch).mockResolvedValueOnce(
    jsonResponse({ code: 'division_by_zero', message: 'division by zero' }, 400),
  )
  render(<App />)

  await press('8')
  await press('Divide')
  await press('0')
  await press('Equals')

  const alert = await screen.findByRole('alert')
  expect(alert).toHaveTextContent(/^division by zero$/)
  expect(alert.tagName).toBe('OUTPUT')
  expect(screen.queryByRole('status', { name: 'Result' })).not.toBeInTheDocument()
  expect(requestBody()).toEqual({ operation: 'divide', a: 8, b: 0 })

  await press('Clear')

  expect(screen.queryByRole('alert')).not.toBeInTheDocument()
  expect(result()).toHaveTextContent(/^0$/)
})

test('shows a message when the calculator service cannot be reached', async () => {
  vi.mocked(fetch).mockRejectedValueOnce(new TypeError('Failed to fetch'))
  render(<App />)

  await press('2')
  await press('Power')
  await press('3')
  await press('Equals')

  expect(await screen.findByRole('alert')).toHaveTextContent('Could not reach the calculator service')
})

test('does not call the API when equals is pressed with no operation', async () => {
  render(<App />)

  await press('4')
  await press('Equals')

  expect(result()).toHaveTextContent('4')
  expect(fetch).not.toHaveBeenCalled()
})
