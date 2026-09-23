import { Tag } from 'antd'

interface ProgressTagProps {
  text: string
  color?: string
}

export default function ProgressTag({ text, color = 'blue' }: ProgressTagProps) {
  return <Tag color={color}>{text}</Tag>
}
