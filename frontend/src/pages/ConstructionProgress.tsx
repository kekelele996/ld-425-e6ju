import { useEffect, useMemo, useState } from 'react'
import { Button, Card, Form, Input, Modal, Select, Space, Table, Typography, message, Upload } from 'antd'
import { CheckOutlined, PlayCircleOutlined, UploadOutlined } from '@ant-design/icons'
import StatusBadge from '@/components/common/StatusBadge'
import Timeline from '@/components/common/Timeline'
import StepIndicator from '@/components/common/StepIndicator'
import { useConstructionStore } from '@/stores/constructionStore'
import { useProjectStore } from '@/stores/projectStore'
import { useAuthStore } from '@/stores/authStore'
import { acceptConstruction, createConstruction, updateConstructionStatus } from '@/api/construction'
import { uploadFile } from '@/utils/upload'
import { extractErrorMessage } from '@/utils/request'
import { ConstructionStatus, Role, ConstructionName } from '@/types/enums'
import type { ConstructionNode } from '@/types'

export default function ConstructionProgress() {
  const { nodes, fetchNodes } = useConstructionStore()
  const { projects, fetchProjects } = useProjectStore()
  const user = useAuthStore((state) => state.user)
  const [projectId, setProjectId] = useState<number>()
  const [acceptTarget, setAcceptTarget] = useState<ConstructionNode | null>(null)
  const [createOpen, setCreateOpen] = useState(false)
  const [photoUrls, setPhotoUrls] = useState<string[]>([])
  const [form] = Form.useForm()
  const [acceptForm] = Form.useForm()

  useEffect(() => {
    fetchProjects()
  }, [fetchProjects])

  useEffect(() => {
    fetchNodes(projectId)
  }, [fetchNodes, projectId])

  const filtered = useMemo(() => {
    if (!projectId) return nodes
    return nodes.filter((n) => n.project_id === projectId)
  }, [nodes, projectId])

  const canOperate = user?.role === Role.Admin || user?.role === Role.Contractor || user?.role === Role.ProjectManager

  const startNode = async (node: ConstructionNode) => {
    try {
      await updateConstructionStatus(node.id, ConstructionStatus.InProgress)
      message.success('节点已开工')
      await fetchNodes(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  const completeNode = async (node: ConstructionNode) => {
    try {
      await updateConstructionStatus(node.id, ConstructionStatus.Completed)
      message.success('节点已完工')
      await fetchNodes(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  const onAccept = async () => {
    const values = await acceptForm.validateFields()
    if (!acceptTarget) return
    try {
      await acceptConstruction(acceptTarget.id, { ...values, photos: photoUrls })
      message.success('验收完成')
      setAcceptTarget(null)
      setPhotoUrls([])
      acceptForm.resetFields()
      await fetchNodes(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  const onCreate = async () => {
    const values = await form.validateFields()
    try {
      await createConstruction({ ...values, project_id: projectId! })
      message.success('创建成功')
      setCreateOpen(false)
      form.resetFields()
      await fetchNodes(projectId)
    } catch (error) {
      message.error(extractErrorMessage(error))
    }
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }} wrap>
        <Typography.Title level={3} style={{ margin: 0 }}>施工进度</Typography.Title>
        <Select
          style={{ width: 240 }}
          placeholder="选择项目"
          allowClear
          value={projectId}
          onChange={setProjectId}
          options={projects.map((p) => ({ label: p.name, value: p.id }))}
        />
        {user?.role === Role.Admin || user?.role === Role.ProjectManager ? (
          <Button type="primary" disabled={!projectId} onClick={() => setCreateOpen(true)}>新增节点</Button>
        ) : null}
      </Space>

      <StepIndicator
        current={filtered.filter((n) => n.status === ConstructionStatus.Completed).length}
        items={filtered.map((n) => n.name)}
      />

      <Card style={{ marginTop: 16 }}>
        <Timeline
          items={filtered.map((node) => ({
            title: node.name,
            status: node.status,
            time: `${node.planned_start_date || '-'} ~ ${node.planned_end_date || '-'}`,
            description: `验收状态：${node.acceptance_status}${node.acceptance_note ? ' · ' + node.acceptance_note : ''}`,
          }))}
        />
      </Card>

      <Card title="节点操作" style={{ marginTop: 16 }}>
        <Table
          rowKey="id"
          dataSource={filtered}
          pagination={false}
          columns={[
            { title: '节点名称', dataIndex: 'name' },
            { title: '状态', dataIndex: 'status', render: (v) => <StatusBadge status={v} /> },
            { title: '验收状态', dataIndex: 'acceptance_status', render: (v) => <StatusBadge status={v} /> },
            { title: '计划开始', dataIndex: 'planned_start_date' },
            { title: '计划结束', dataIndex: 'planned_end_date' },
            {
              title: '操作',
              render: (_, record) =>
                canOperate ? (
                  <Space>
                    {record.status === ConstructionStatus.Pending ? (
                      <Button size="small" icon={<PlayCircleOutlined />} onClick={() => startNode(record)}>开工</Button>
                    ) : null}
                    {record.status === ConstructionStatus.InProgress ? (
                      <Button size="small" type="primary" onClick={() => completeNode(record)}>完工</Button>
                    ) : null}
                    {record.status === ConstructionStatus.Completed ? (
                      <Button size="small" icon={<CheckOutlined />} onClick={() => setAcceptTarget(record)}>验收</Button>
                    ) : null}
                  </Space>
                ) : null,
            },
          ]}
        />
      </Card>

      <Modal title="新增施工节点" open={createOpen} onOk={onCreate} onCancel={() => setCreateOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="节点名称" rules={[{ required: true }]}>
            <Select options={ConstructionName.map((n) => ({ label: n, value: n }))} />
          </Form.Item>
          <Form.Item name="planned_start_date" label="计划开始日期"><Input placeholder="YYYY-MM-DD" /></Form.Item>
          <Form.Item name="planned_end_date" label="计划结束日期"><Input placeholder="YYYY-MM-DD" /></Form.Item>
        </Form>
      </Modal>

      <Modal title="施工验收" open={!!acceptTarget} onOk={onAccept} onCancel={() => setAcceptTarget(null)}>
        <Form form={acceptForm} layout="vertical" initialValues={{ accepted: true }}>
          <Form.Item name="accepted" label="验收结论" rules={[{ required: true }]}>
            <Select options={[{ label: '验收通过', value: true }, { label: '验收不通过', value: false }]} />
          </Form.Item>
          <Form.Item name="note" label="验收说明"><Input.TextArea rows={3} /></Form.Item>
          <Form.Item label="验收照片">
            <Upload
              beforeUpload={async (file) => {
                try {
                  const url = await uploadFile(file)
                  setPhotoUrls((prev) => [...prev, url])
                } catch (error) {
                  message.error(extractErrorMessage(error))
                }
                return false
              }}
            >
              <Button icon={<UploadOutlined />}>上传照片</Button>
            </Upload>
            {photoUrls.map((url) => (
              <div key={url} style={{ marginTop: 8 }}>{url}</div>
            ))}
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
