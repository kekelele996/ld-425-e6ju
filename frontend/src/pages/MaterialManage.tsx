import { useEffect, useMemo, useState } from 'react'
import { Button, Card, Form, Input, InputNumber, Modal, Select, Space, Table, Typography, message, Statistic } from 'antd'
import { PlusOutlined, ShoppingCartOutlined } from '@ant-design/icons'
import StatusBadge from '@/components/common/StatusBadge'
import ProgressTag from '@/components/common/ProgressTag'
import EmptyState from '@/components/common/EmptyState'
import { useMaterialStore } from '@/stores/materialStore'
import { useProjectStore } from '@/stores/projectStore'
import { useBudgetStore } from '@/stores/budgetStore'
import { useAuthStore } from '@/stores/authStore'
import { createMaterial, updateMaterialStatus } from '@/api/material'
import { extractErrorMessage } from '@/utils/request'
import { formatCurrency } from '@/utils/formatBudget'
import { MaterialCategory, MaterialSpace, PurchaseStatus, Role } from '@/types/enums'
import type { MaterialItem } from '@/types'

const purchaseSteps = [PurchaseStatus.NotPurchased, PurchaseStatus.Ordered, PurchaseStatus.Delivered, PurchaseStatus.Installed]

export default function MaterialManage() {
  const { materials, fetchMaterials } = useMaterialStore()
  const { projects, fetchProjects } = useProjectStore()
  const { budgets, fetchBudgets } = useBudgetStore()
  const user = useAuthStore((state) => state.user)
  const [projectId, setProjectId] = useState<number>()
  const [category, setCategory] = useState<string>()
  const [space, setSpace] = useState<string>()
  const [createOpen, setCreateOpen] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchProjects()
    fetchBudgets()
  }, [fetchProjects, fetchBudgets])

  useEffect(() => {
    fetchMaterials(projectId)
  }, [fetchMaterials, projectId])

  const filtered = useMemo(() => {
    return materials.filter((item) => {
      if (category && item.category !== category) return false
      if (space && item.space !== space) return false
      return true
    })
  }, [materials, category, space])

  const totalCost = filtered.reduce((sum, item) => sum + item.total_price, 0)
  const materialBudget = budgets.filter((b) => b.category === 'Material').reduce((sum, b) => sum + b.budget_amount, 0)

  const canEdit = user?.role === Role.Admin || user?.role === Role.Designer || user?.role === Role.ProjectManager
  const canProcure = user?.role === Role.Admin || user?.role === Role.Designer || user?.role === Role.Contractor || user?.role === Role.ProjectManager

  const onCreate = async () => {
    const values = await form.validateFields()
    try {
      await createMaterial({ ...values, project_id: projectId! })
      message.success('创建成功')
      setCreateOpen(false)
      form.resetFields()
      await fetchMaterials(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  const advanceStatus = async (record: MaterialItem) => {
    const current = purchaseSteps.indexOf(record.purchase_status)
    const next = purchaseSteps[current + 1]
    if (!next) return
    try {
      await updateMaterialStatus(record.id, next)
      message.success('采购状态已更新')
      await fetchMaterials(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <Typography.Title level={3} style={{ margin: 0 }}>材料管理</Typography.Title>
        <Select
          style={{ width: 240 }}
          placeholder="选择项目"
          allowClear
          value={projectId}
          onChange={setProjectId}
          options={projects.map((p) => ({ label: p.name, value: p.id }))}
        />
        <Select
          style={{ width: 140 }}
          placeholder="品类"
          allowClear
          value={category}
          onChange={setCategory}
          options={MaterialCategory.map((c) => ({ label: c, value: c }))}
        />
        <Select
          style={{ width: 140 }}
          placeholder="使用空间"
          allowClear
          value={space}
          onChange={setSpace}
          options={MaterialSpace.map((s) => ({ label: s, value: s }))}
        />
        {canEdit ? (
          <Button type="primary" icon={<PlusOutlined />} disabled={!projectId} onClick={() => setCreateOpen(true)}>
            新增材料
          </Button>
        ) : null}
      </Space>

      <Card style={{ marginBottom: 16 }}>
        <Space size="large">
          <Statistic title="材料费用汇总" value={totalCost} precision={2} prefix="¥" />
          <Statistic title="材料预算" value={materialBudget} precision={2} prefix="¥" />
        </Space>
      </Card>

      <Card>
        {filtered.length === 0 ? (
          <EmptyState description="暂无材料" />
        ) : (
          <Table
            rowKey="id"
            dataSource={filtered}
            pagination={{ pageSize: 10 }}
            columns={[
              { title: '材料名称', dataIndex: 'name' },
              { title: '品类', dataIndex: 'category', render: (v) => <ProgressTag text={v} color="cyan" /> },
              { title: '规格型号', dataIndex: 'spec' },
              { title: '品牌', dataIndex: 'brand' },
              { title: '数量', dataIndex: 'quantity' },
              { title: '单价', dataIndex: 'unit_price', render: (v) => formatCurrency(v) },
              { title: '总价', dataIndex: 'total_price', render: (v) => formatCurrency(v) },
              { title: '采购状态', dataIndex: 'purchase_status', render: (v) => <StatusBadge status={v} /> },
              { title: '使用空间', dataIndex: 'space' },
              {
                title: '操作',
                render: (_, record) =>
                  canProcure && record.purchase_status !== PurchaseStatus.Installed ? (
                    <Button
                      size="small"
                      icon={<ShoppingCartOutlined />}
                      onClick={() => advanceStatus(record)}
                    >
                      推进采购
                    </Button>
                  ) : null,
              },
            ]}
          />
        )}
      </Card>

      <Modal title="新增材料" open={createOpen} onOk={onCreate} onCancel={() => setCreateOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="材料名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="category" label="品类" rules={[{ required: true }]}>
            <Select options={MaterialCategory.map((c) => ({ label: c, value: c }))} />
          </Form.Item>
          <Form.Item name="spec" label="规格型号"><Input /></Form.Item>
          <Form.Item name="brand" label="品牌"><Input /></Form.Item>
          <Form.Item name="quantity" label="数量" rules={[{ required: true }]}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="unit" label="单位" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="unit_price" label="单价"><InputNumber min={0} style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="supplier" label="供应商"><Input /></Form.Item>
          <Form.Item name="space" label="使用空间">
            <Select options={MaterialSpace.map((s) => ({ label: s, value: s }))} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
