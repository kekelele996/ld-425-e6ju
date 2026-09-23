import { Tag } from 'antd'

interface StatusBadgeProps {
  status: string
}

const colorMap: Record<string, string> = {
  Designing: 'blue',
  Quoting: 'gold',
  InProgress: 'processing',
  Completed: 'green',
  Archived: 'default',
  NotStarted: 'default',
  Revision: 'orange',
  Approved: 'green',
  NotPurchased: 'default',
  Ordered: 'blue',
  Delivered: 'cyan',
  Installed: 'green',
  Pending: 'default',
  Delayed: 'red',
  Passed: 'green',
  Failed: 'red',
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  return <Tag color={colorMap[status] || 'default'}>{status}</Tag>
}
