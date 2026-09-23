import { Steps } from 'antd'

interface StepIndicatorProps {
  current: number
  items: string[]
}

export default function StepIndicator({ current, items }: StepIndicatorProps) {
  return <Steps current={current} size="small" items={items.map((title) => ({ title }))} />
}
