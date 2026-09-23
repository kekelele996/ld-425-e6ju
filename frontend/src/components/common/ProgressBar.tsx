import { Progress } from 'antd'

interface ProgressBarProps {
  percent: number
  status?: 'normal' | 'success' | 'exception' | 'active'
  label?: string
}

export default function ProgressBar({ percent, status = 'normal', label }: ProgressBarProps) {
  return (
    <div>
      <Progress percent={Math.max(0, Math.min(100, Math.round(percent)))} status={status} />
      {label ? <div style={{ color: '#8c8c8c', fontSize: 12 }}>{label}</div> : null}
    </div>
  )
}
