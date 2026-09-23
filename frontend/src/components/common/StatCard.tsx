import { Card, Statistic } from 'antd'

interface StatCardProps {
  title: string
  value: number
  prefix?: string
  suffix?: string
  precision?: number
  valueStyle?: React.CSSProperties
}

export default function StatCard({ title, value, prefix, suffix, precision = 2, valueStyle }: StatCardProps) {
  return (
    <Card>
      <Statistic title={title} value={value} prefix={prefix} suffix={suffix} precision={precision} valueStyle={valueStyle} />
    </Card>
  )
}
