import { Alert } from 'antd'

interface AlertBannerProps {
  type?: 'success' | 'info' | 'warning' | 'error'
  message: string
  description?: string
  showIcon?: boolean
}

export default function AlertBanner({ type = 'info', message, description, showIcon = true }: AlertBannerProps) {
  return <Alert type={type} message={message} description={description} showIcon={showIcon} />
}
