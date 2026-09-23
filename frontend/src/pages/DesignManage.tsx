import { useEffect, useMemo, useState } from 'react'
import { Button, Card, Form, Input, InputNumber, Modal, Select, Space, Table, Typography, message, Descriptions } from 'antd'
import { PlusOutlined, SendOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'
import StatusBadge from '@/components/common/StatusBadge'
import StepIndicator from '@/components/common/StepIndicator'
import VersionTag from '@/components/common/VersionTag'
import EmptyState from '@/components/common/EmptyState'
import { useDesignStore } from '@/stores/designStore'
import { useProjectStore } from '@/stores/projectStore'
import { useAuthStore } from '@/stores/authStore'
import { createDesign, reviewDesign, submitDesign } from '@/api/design'
import { extractErrorMessage } from '@/utils/request'
import { PhaseStatus, Role, DesignPhaseName } from '@/types/enums'
import type { DesignPhase } from '@/types'

const phaseSteps = ['方案设计', '效果图', '施工图', '软装方案']

export default function DesignManage() {
  const { designs, fetchDesigns } = useDesignStore()
  const { projects, fetchProjects } = useProjectStore()
  const user = useAuthStore((state) => state.user)
  const [projectId, setProjectId] = useState<number>()
  const [createOpen, setCreateOpen] = useState(false)
  const [reviewTarget, setReviewTarget] = useState<DesignPhase | null>(null)
  const [compareTarget, setCompareTarget] = useState<DesignPhase | null>(null)
  const [form] = Form.useForm()
  const [reviewForm] = Form.useForm()

  useEffect(() => {
    fetchProjects()
  }, [fetchProjects])

  useEffect(() => {
    fetchDesigns(projectId)
  }, [fetchDesigns, projectId])

  const filteredDesigns = useMemo(() => {
    if (!projectId) return designs
    return designs.filter((d) => d.project_id === projectId)
  }, [designs, projectId])

  const canEdit = user?.role === Role.Admin || user?.role === Role.Designer || user?.role === Role.ProjectManager
  const canReview = user?.role === Role.Admin || user?.role === Role.Owner

  const onCreate = async () => {
    const values = await form.validateFields()
    try {
      await createDesign({ ...values, project_id: projectId! })
      message.success('创建成功')
      setCreateOpen(false)
      form.resetFields()
      await fetchDesigns(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  const onReview = async () => {
    const values = await reviewForm.validateFields()
    if (!reviewTarget) return
    try {
      await reviewDesign(reviewTarget.id, values)
      message.success('审核完成')
      setReviewTarget(null)
      reviewForm.resetFields()
      await fetchDesigns(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  const onCompare = (record: DesignPhase) => {
    setCompareTarget(record)
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <Typography.Title level={3} style={{ margin: 0 }}>设计管理</Typography.Title>
        <Select
          style={{ width: 260 }}
          placeholder="选择项目"
          allowClear
          value={projectId}
          onChange={setProjectId}
          options={projects.map((p) => ({ label: p.name, value: p.id }))}
        />
        {canEdit ? (
          <Button type="primary" icon={<PlusOutlined />} disabled={!projectId} onClick={() => setCreateOpen(true)}>
            新增设计阶段
          </Button>
        ) : null}
      </Space>

      <StepIndicator current={filteredDesigns.filter((d) => d.status === PhaseStatus.Approved).length} items={phaseSteps} />

      <Card style={{ marginTop: 16 }}>
        {filteredDesigns.length === 0 ? (
          <EmptyState description="暂无设计阶段" />
        ) : (
          <Table
            rowKey="id"
            dataSource={filteredDesigns}
            pagination={false}
            columns={[
              { title: '阶段名称', dataIndex: 'name' },
              { title: '状态', dataIndex: 'status', render: (v) => <StatusBadge status={v} /> },
              { title: '版本', dataIndex: 'version', render: (v) => <VersionTag version={v} /> },
              { title: '设计说明', dataIndex: 'description', ellipsis: true },
              { title: '审核意见', dataIndex: 'review_comment', ellipsis: true },
              {
                title: '操作',
                render: (_, record) => (
                  <Space>
                    {canEdit ? (
                      <Button size="small" icon={<SendOutlined />} onClick={async () => {
                        try {
                          await submitDesign(record.id, {})
                          message.success('已提交审核')
                          await fetchDesigns(projectId)
                        } catch (error) {
                          message.error(extractErrorMessage(error))
                        }
                      }}>
                        提交审核
                      </Button>
                    ) : null}
                    {canReview ? (
                      <Button size="small" type="primary" icon={<CheckOutlined />} onClick={() => setReviewTarget(record)}>
                        审核
                      </Button>
                    ) : null}
                    <Button size="small" onClick={() => onCompare(record)}>版本对比</Button>
                  </Space>
                ),
              },
            ]}
          />
        )}
      </Card>

      <Modal title="新增设计阶段" open={createOpen} onOk={onCreate} onCancel={() => setCreateOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="阶段名称" rules={[{ required: true }]}>
            <Select options={DesignPhaseName.map((n) => ({ label: n, value: n }))} />
          </Form.Item>
          <Form.Item name="designer_id" label="设计师 ID">
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="description" label="设计说明">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal title="审核设计" open={!!reviewTarget} onOk={onReview} onCancel={() => setReviewTarget(null)}>
        <Form form={reviewForm} layout="vertical" initialValues={{ approved: true }}>
          <Form.Item name="approved" label="审核结论" rules={[{ required: true }]}>
            <Select options={[{ label: '通过', value: true }, { label: '驳回修改', value: false }]} />
          </Form.Item>
          <Form.Item name="comment" label="审核意见">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="版本历史对比"
        open={!!compareTarget}
        footer={<Button onClick={() => setCompareTarget(null)}>关闭</Button>}
        width={700}
        onCancel={() => setCompareTarget(null)}
      >
        {compareTarget ? (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="阶段">{compareTarget.name}</Descriptions.Item>
            <Descriptions.Item label="状态"><StatusBadge status={compareTarget.status} /></Descriptions.Item>
            <Descriptions.Item label="版本"><VersionTag version={compareTarget.version} /></Descriptions.Item>
            <Descriptions.Item label="设计说明">{compareTarget.description || '-'}</Descriptions.Item>
            <Descriptions.Item label="文件">{compareTarget.file_urls?.join(', ') || '-'}</Descriptions.Item>
            <Descriptions.Item label="审核意见">{compareTarget.review_comment || '-'}</Descriptions.Item>
          </Descriptions>
        ) : null}
      </Modal>
    </div>
  )
}
