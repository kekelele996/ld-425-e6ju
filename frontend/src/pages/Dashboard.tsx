import { useEffect, useMemo, useState } from 'react'
import { Card, Col, Row, Spin, Typography, List, Tag } from 'antd'
import ReactECharts from 'echarts-for-react'
import StatusBadge from '@/components/common/StatusBadge'
import ProgressBar from '@/components/common/ProgressBar'
import AlertBanner from '@/components/common/AlertBanner'
import { useProjectStore } from '@/stores/projectStore'
import { useConstructionStore } from '@/stores/constructionStore'
import { useBudgetStore } from '@/stores/budgetStore'
import { formatCurrency } from '@/utils/formatBudget'
import { ConstructionStatus, AcceptanceStatus } from '@/types/enums'

export default function Dashboard() {
  const { projects, fetchProjects } = useProjectStore()
  const { nodes, fetchNodes } = useConstructionStore()
  const { budgets, fetchBudgets } = useBudgetStore()
  const [loaded, setLoaded] = useState(false)

  useEffect(() => {
    Promise.all([fetchProjects(), fetchNodes(), fetchBudgets()]).finally(() => setLoaded(true))
  }, [fetchProjects, fetchNodes, fetchBudgets])

  const activeProjects = useMemo(() => projects.filter((p) => p.status !== 'Archived'), [projects])

  const ganttOption = useMemo(() => {
    const rows = nodes.slice(0, 7).map((node) => ({
      name: node.name,
      start: node.planned_start_date || '',
      end: node.planned_end_date || '',
    }))
    return {
      tooltip: { trigger: 'item' },
      grid: { left: 90, right: 20, top: 20, bottom: 30 },
      xAxis: { type: 'time' },
      yAxis: { type: 'category', data: rows.map((r) => r.name) },
      series: [
        {
          type: 'custom',
          renderItem: (_params: unknown, api: any) => {
            const categoryIndex = api.value(0)
            const start = api.coord([api.value(1), categoryIndex])
            const end = api.coord([api.value(2), categoryIndex])
            const height = api.size([0, 1])[1] * 0.5
            return {
              type: 'rect',
              shape: { x: start[0], y: start[1] - height / 2, width: Math.max(end[0] - start[0], 2), height },
              style: { fill: '#1677ff' },
            }
          },
          encode: { x: [1, 2], y: 0 },
          data: rows.map((r) => [r.name, r.start, r.end]),
        },
      ],
    }
  }, [nodes])

  const budgetOption = useMemo(() => {
    const totalBudget = budgets.reduce((sum, item) => sum + item.budget_amount, 0)
    const totalActual = budgets.reduce((sum, item) => sum + item.actual_amount, 0)
    return {
      tooltip: { trigger: 'item' },
      series: [
        {
          type: 'pie',
          radius: ['58%', '80%'],
          label: { show: true, formatter: '{b}\n{d}%' },
          data: [
            { name: '已执行', value: totalActual },
            { name: '剩余预算', value: Math.max(totalBudget - totalActual, 0) },
          ],
        },
      ],
    }
  }, [budgets])

  const pendingAcceptance = nodes.filter(
    (node) => node.status === ConstructionStatus.Completed && node.acceptance_status === AcceptanceStatus.Pending,
  )
  const delayedNodes = nodes.filter((node) => node.status === ConstructionStatus.Delayed)

  if (!loaded) return <Spin style={{ display: 'block', margin: '80px auto' }} />

  return (
    <div>
      <Typography.Title level={3}>项目总览</Typography.Title>

      {delayedNodes.length > 0 ? (
        <AlertBanner type="error" message={`${delayedNodes.length} 个施工节点已延期`} description={delayedNodes.map((n) => n.name).join('、')} />
      ) : null}

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        {activeProjects.map((project) => (
          <Col xs={24} md={12} xl={8} key={project.id}>
            <Card
              title={project.name}
              extra={<StatusBadge status={project.status} />}
            >
              <Typography.Paragraph type="secondary">{project.house_type} · {project.area}m² · {project.decor_style}</Typography.Paragraph>
              <Typography.Paragraph>合同金额：{formatCurrency(project.contract_amount)}</Typography.Paragraph>
              <ProgressBar
                percent={project.status === 'Completed' ? 100 : project.status === 'InProgress' ? 55 : 15}
                status="active"
                label="施工进度（示意）"
              />
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={14}>
          <Card title="施工进度甘特图">
            <ReactECharts option={ganttOption} style={{ height: 320 }} />
          </Card>
        </Col>
        <Col xs={24} lg={10}>
          <Card title="预算执行率">
            <ReactECharts option={budgetOption} style={{ height: 320 }} />
          </Card>
        </Col>
      </Row>

      <Card title="待验收节点提醒" style={{ marginTop: 16 }}>
        {pendingAcceptance.length === 0 ? (
          <Typography.Text type="secondary">暂无待验收节点</Typography.Text>
        ) : (
          <List
            dataSource={pendingAcceptance}
            renderItem={(node) => (
              <List.Item>
                <List.Item.Meta title={node.name} description={`计划完成 ${node.planned_end_date || '-'}`} />
                <Tag color="orange">待验收</Tag>
              </List.Item>
            )}
          />
        )}
      </Card>
    </div>
  )
}
