import { useEffect } from 'react'
import { useCalculator } from './useCalculator'
import type { BinaryOperation } from './api'
import './Calculator.css'

export function Calculator() {
  const calculator = useCalculator()

  useEffect(() => {
    function onKeyDown(event: KeyboardEvent) {
      if (event.metaKey || event.ctrlKey || event.altKey) {
        return
      }
      if (/^[0-9]$/.test(event.key)) {
        event.preventDefault()
        calculator.inputDigit(event.key)
        return
      }
      const operations: Record<string, BinaryOperation> = {
        '+': 'add',
        '-': 'subtract',
        '*': 'multiply',
        '/': 'divide',
        '%': 'percentage',
        '^': 'power',
      }
      if (event.key in operations) {
        event.preventDefault()
        calculator.pressOperator(operations[event.key])
        return
      }
      if (event.key === 'Enter' || event.key === '=') {
        event.preventDefault()
        calculator.pressEquals()
      } else if (event.key === '.') {
        event.preventDefault()
        calculator.inputDecimal()
      } else if (event.key === 'Backspace') {
        event.preventDefault()
        calculator.deleteDigit()
      } else if (event.key === 'Escape') {
        event.preventDefault()
        calculator.clear()
      }
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [calculator])

  const statusText = calculator.error ? '' : calculator.loading ? 'Calculating…' : calculator.preview

  return (
    <main className="page">
      <section className="device" aria-busy={calculator.loading} aria-label="Calculator">
        <h1>Calculator</h1>
        <div className="screen">
          <p className="preview">{statusText || '\u00a0'}</p>
          <output
            className={calculator.error ? 'display error' : 'display'}
            aria-label={calculator.error ? undefined : 'Result'}
            role={calculator.error ? 'alert' : undefined}
          >
            {calculator.error ?? calculator.display}
          </output>
        </div>
        <div className="keys">
          <Key name="Clear" label="C" className="function" disabled={false} onClick={calculator.clear} />
          <Key name="Delete" label="DEL" className="function" disabled={calculator.loading} onClick={calculator.deleteDigit} />
          <Key name="Square root" label="√" className="function" disabled={calculator.loading} onClick={calculator.pressSqrt} />
          <Operator name="Divide" label="÷" operation="divide" pending={calculator.pending} disabled={calculator.loading} onPress={calculator.pressOperator} />
          <Digit digit="7" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Digit digit="8" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Digit digit="9" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Operator name="Multiply" label="×" operation="multiply" pending={calculator.pending} disabled={calculator.loading} onPress={calculator.pressOperator} />
          <Digit digit="4" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Digit digit="5" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Digit digit="6" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Operator name="Subtract" label="−" operation="subtract" pending={calculator.pending} disabled={calculator.loading} onPress={calculator.pressOperator} />
          <Digit digit="1" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Digit digit="2" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Digit digit="3" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Operator name="Add" label="+" operation="add" pending={calculator.pending} disabled={calculator.loading} onPress={calculator.pressOperator} />
          <Operator name="Percentage" label="%" operation="percentage" variant="function" pending={calculator.pending} disabled={calculator.loading} onPress={calculator.pressOperator} />
          <Digit digit="0" loading={calculator.loading} onDigit={calculator.inputDigit} />
          <Key name="Decimal" label="." disabled={calculator.loading} onClick={calculator.inputDecimal} />
          <Operator name="Power" label="x^y" operation="power" variant="function" pending={calculator.pending} disabled={calculator.loading} onPress={calculator.pressOperator} />
          <Key name="Equals" label="=" className="equals" disabled={calculator.loading} onClick={calculator.pressEquals} />
        </div>
      </section>
    </main>
  )
}

function Digit({
  digit,
  loading,
  onDigit,
}: {
  digit: string
  loading: boolean
  onDigit: (digit: string) => void
}) {
  return <Key name={digit} label={digit} disabled={loading} onClick={() => onDigit(digit)} />
}

function Operator({
  name,
  label,
  operation,
  pending,
  disabled,
  onPress,
  variant = 'operator',
}: {
  name: string
  label: string
  operation: BinaryOperation
  pending: BinaryOperation | null
  disabled: boolean
  onPress: (operation: BinaryOperation) => void
  variant?: 'operator' | 'function'
}) {
  const active = pending === operation
  const className = active ? `${variant} active` : variant
  return (
    <Key
      name={name}
      label={label}
      className={className}
      disabled={disabled}
      pressed={active}
      onClick={() => onPress(operation)}
    />
  )
}

function Key({
  name,
  label,
  className,
  disabled,
  pressed,
  onClick,
}: {
  name: string
  label: string
  className?: string
  disabled: boolean
  pressed?: boolean
  onClick: () => void
}) {
  return (
    <button
      type="button"
      className={className ? `key ${className}` : 'key'}
      aria-label={name}
      aria-pressed={pressed}
      disabled={disabled}
      onClick={onClick}
    >
      {label}
    </button>
  )
}
