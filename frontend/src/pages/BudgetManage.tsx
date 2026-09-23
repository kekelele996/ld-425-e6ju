import { useEffect, useMemo, useState } from 'react'
import { Button, Card, Form, Input, InputNumber, Modal, Select, Space, Table, Typography, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import ReactECharts from 'echarts-for-react'
import StatusBadge from '@/components/common/StatusBadge'
import AlertBanner from '@/components/common/AlertBanner'
import StatCard from '@/components/common/StatCard'
import { useBudgetStore } from '@/stores/budgetStore'
import { useProjectStore } from '@/stores/projectStore'
import { useMaterialStore } from '@/stores/materialStore'
import { useAuthStore } from '@/stores/authStore'
import { createBudget } from '@/api/budget'
import { extractErrorMessage } from '@/utils/request'
import { formatCurrency } from '@/utils/formatBudget'
import { BudgetCategory, Role } from '@/types/enums'
import type { BudgetItem } from '@/types'

export default function BudgetManage() {
  const { budgets, fetchBudgets } = useBudgetStore()
  const { projects, fetchProjects } = useProjectStore()
  const { materials, fetchMaterials } = useMaterialStore()
  const user = useAuthStore((state) => state.user)
  const [projectId, setProjectId] = useState<number>()
  const [createOpen, setCreateOpen] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchProjects()
    fetchMaterials()
  }, [fetchProjects, fetchMaterials])

  useEffect(() => {
    fetchBudgets(projectId)
  }, [fetchBudgets, projectId])

  const filtered = useMemo(() => {
    if (!projectId) return budgets
    return budgets.filter((b) => b.project_id === projectId)
  }, [budgets, projectId])

  const totalBudget = filtered.reduce((sum, item) => sum + item.budget_amount, 0)
  const totalActual = filtered.reduce((sum, item) => sum + item.actual_amount, 0)
  const totalMaterialCost = materials.filter((m) => !projectId || m.project_id === projectId).reduce((sum, m) => sum + m.total_price, 0)

  const overBudgetItems = filtered.filter((item) => item.variance > 0)

  const chartOption = useMemo(() => ({
    tooltip: { trigger: 'axis' },
    legend: { data: ['预算', '实际'] },
    grid: { left: 60, right: 20, top: 40, bottom: 30 },
    xAxis: { type: 'category', data: filtered.map((item) => item.category) },
    yAxis: { type: 'value' },
    series: [
      { name: '预算', type: 'bar', data: filtered.map((item) => item.budget_amount), itemStyle: { color: '#91caff' } },
      { name: '实际', type: 'bar', data: filtered.map((item) => item.actual_amount), itemStyle: { color: '#1677ff' } },
    ],
  }), [filtered])

  const canEdit = user?.role === Role.Admin || user?.role === Role.ProjectManager

  const onCreate = async () => {
    const values = await form.validateFields()
    try {
      await createBudget({ ...values, project_id: projectId! })
      message.success('创建成功')
      setCreateOpen(false)
      form.resetFields()
      await fetchBudgets(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <Typography.Title level={3} style={{ margin: 0 }}>预算管理</Typography.Title>
        <Select
          style={{ width: 240 }}
          placeholder="选择项目"
          allowClear
          value={projectId}
          onChange={setProjectId}
          options={projects.map((p) => ({ label: p.name, value: p.id }))}
        />
        {canEdit ? (
          <Button type="primary" icon={<PlusOutlined />} disabled={!projectId} onClick={() => setCreateOpen(true)}>
            新增预算项
          </Button>
        ) : null}
      </Space>

      {overBudgetItems.length > 0 ? (
        <AlertBanner type="warning" message={`${overBudgetItems.length} 个预算项超支`} description={overBudgetItems.map((i) => i.category).join('、')} />
      ) : null}

      <Space size="middle" style={{ margin: '16px 0' }} wrap>
        <StatCard title="预算总额" value={totalBudget} prefix="¥" />
        <StatCard title="实际花费" value={totalActual} prefix="¥" />
        <StatCard title="材料费用" value={totalMaterialCost} prefix="¥" />
      </Space>

      <Card title="预算 vs 实际" style={{ marginBottom: 16 }}>
        <ReactECharts option={chartOption} style={{ height: 320 }} />
      </Card>

      <Card>
        <Table<BudgetItem>
          rowKey="id"
          dataSource={filtered}
          pagination={false}
          columns={[
            { title: '预算类别', dataIndex: 'category', render: (v) => <StatusBadge status={v} /> },
            { title: '预算金额', dataIndex: 'budget_amount', render: (v) => formatCurrency(v) },
            { title: '实际花费', dataIndex: 'actual_amount', render: (v) => formatCurrency(v) },
            {
              title: '差异金额',
              dataIndex: 'variance',
              render: (v: number) => <span style={{ color: v > 0 ? '#cf1322' : '#389e0d' }}>{formatCurrency(v)}</span>,
            },
            { title: '备注', dataIndex: 'remark' },
          ]}
        />
      </Card>

      <Modal title="新增预算项" open={createOpen} onOk={onCreate} onCancel={() => setCreateOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="category" label="预算类别" rules={[{ required: true }]}>
            <Select options={BudgetCategory.map((c) => ({ label: c, value: c }))} />
          </Form.Item>
          <Form.Item name="budget_amount" label="预算金额"><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="actual_amount" label="实际花费"><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="remark" label="备注"><Input /></Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
