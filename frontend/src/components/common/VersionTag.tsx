import { Tag } from 'antd'

interface VersionTagProps {
  version: number
}

export default function VersionTag({ version }: VersionTagProps) {
  return <Tag color="purple">v{version}</Tag>
}
