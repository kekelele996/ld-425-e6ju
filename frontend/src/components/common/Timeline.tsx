import { Timeline as AntTimeline, Tag } from 'antd'

export interface TimelineItem {
  title: string
  time?: string
  status?: string
  color?: string
  description?: string
}

interface TimelineProps {
  items: TimelineItem[]
}

const colorMap: Record<string, string> = {
  Pending: 'gray',
  InProgress: 'blue',
  Completed: 'green',
  Delayed: 'red',
}

export default function Timeline({ items }: TimelineProps) {
  return (
    <AntTimeline
      items={items.map((item) => ({
        color: item.color || colorMap[item.status || ''] || 'blue',
        children: (
          <div>
            <div style={{ fontWeight: 600 }}>
              {item.title}
              {item.status ? <Tag style={{ marginLeft: 8 }}>{item.status}</Tag> : null}
            </div>
            {item.time ? <div style={{ color: '#8c8c8c', fontSize: 12 }}>{item.time}</div> : null}
            {item.description ? <div style={{ color: '#595959' }}>{item.description}</div> : null}
          </div>
        ),
      }))}
    />
  )
}
