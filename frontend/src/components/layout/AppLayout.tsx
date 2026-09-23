import { Layout, Menu, Dropdown, Space, Avatar, Typography } from 'antd'
import {
  DashboardOutlined,
  BulbOutlined,
  GoldOutlined,
  ToolOutlined,
  AccountBookOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import { Role } from '@/types/enums'

const { Header, Sider, Content } = Layout

const menuItems = [
  { key: '/dashboard', icon: <DashboardOutlined />, label: '项目总览' },
  { key: '/designs', icon: <BulbOutlined />, label: '设计管理' },
  { key: '/materials', icon: <GoldOutlined />, label: '材料管理' },
  { key: '/construction', icon: <ToolOutlined />, label: '施工进度' },
  { key: '/budget', icon: <AccountBookOutlined />, label: '预算管理' },
]

const roleLabels: Record<string, string> = {
  [Role.Admin]: '管理员',
  [Role.Designer]: '设计师',
  [Role.Contractor]: '施工方',
  [Role.Owner]: '业主',
  [Role.ProjectManager]: '项目经理',
}

export default function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const user = useAuthStore((state) => state.user)
  const logout = useAuthStore((state) => state.logout)

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider breakpoint="lg" collapsedWidth="0">
        <div style={{ color: '#fff', padding: 20, fontSize: 18, fontWeight: 700 }}>装修管理平台</div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end' }}>
          <Dropdown
            menu={{
              items: [
                {
                  key: 'logout',
                  icon: <LogoutOutlined />,
                  label: '退出登录',
                  onClick: () => {
                    logout()
                    navigate('/login')
                  },
                },
              ],
            }}
          >
            <Space style={{ cursor: 'pointer' }}>
              <Avatar>{user?.name?.slice(0, 1) || 'U'}</Avatar>
              <Typography.Text>{user?.name || user?.username}</Typography.Text>
              <Typography.Text type="secondary">{roleLabels[user?.role || ''] || user?.role}</Typography.Text>
            </Space>
          </Dropdown>
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
