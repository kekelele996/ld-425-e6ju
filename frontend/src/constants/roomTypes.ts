// 房间/空间类型常量。
export const ROOM_TYPES = ['客厅', '卧室', '厨房', '卫生间', '阳台'] as const
export type RoomType = (typeof ROOM_TYPES)[number]

export const HOUSE_TYPES = ['一室', '两室', '三室', '四室', '别墅', '复式'] as const
export type HouseType = (typeof HOUSE_TYPES)[number]
